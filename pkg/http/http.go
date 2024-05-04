package http

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/peterbourgon/diskv/v3"
	statsd "github.com/smira/go-statsd"
)

const RetryCount = 2

type Client struct {
	httpClient *http.Client
	db         *sql.DB
	sd         *statsd.Client
}

type retryableTransport struct {
	transport http.RoundTripper
	limiter   *Limiter
	cache     *diskv.Diskv
	db        *sql.DB
	sd        *statsd.Client
}

func (t *retryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	// Do we have a context?

	req.Header.Set("User-Agent", "curl/8.4.0")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15")
	req.Header.Set("Accepts", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("DNT", "1")

	t.sd.Incr("http.roundtrip.start", 1)

	var err error

	// Check for existing
	existing, err := getURLRequestFromDB(req.URL.String(), t.db)
	if err != nil {
		return nil, err
	}

	var urlRequest *URLRequest

	if existing != nil {
		urlRequest = existing
	} else {
		urlRequest = &URLRequest{
			Url:           req.URL.String(),
			LastAttemptAt: time.Now().Unix(),
		}
	}

	checkDisk := true

	if urlRequest.LastAttemptAt < (time.Now().Unix() - 60*60*24*5) {
		checkDisk = false
	}
	if urlRequest.Status == 403 {
		checkDisk = false
	}

	k, err := getHTMLKey(req)
	if err != nil {
		return nil, err
	}

	if checkDisk {
		t.sd.Incr("http.roundtrip.checkingDisk", 1)
		// Check our cache

		diskStream, _ := t.cache.ReadStream(k, false)

		if diskStream != nil {
			t.sd.Incr("http.roundtrip.diskCacheHit", 1)
			defer diskStream.Close()
			data, err := io.ReadAll(diskStream)

			if err != nil {
				return nil, err
			}
			buf := bytes.NewBuffer(data)
			bufReader := bufio.NewReader(buf)

			// todo check how old the response is based on the code

			return http.ReadResponse(bufReader, req)
		}
	}
	t.sd.Incr("http.roundtrip.diskCacheMiss", 1)

	urlRequest.LastAttemptAt = time.Now().Unix()
	urlRequest.Save(t.db)

	// ctx, cancelFunc := context.WithTimeout(context.Background(), 6*time.Second)
	// defer cancelFunc()

	// t.limiter.Wait(ctx)

	// t.limiter.WaitHost(ctx, req.Host)

	// Send the request
	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		if resp != nil {
			urlRequest.Status = int64(resp.StatusCode)
			urlRequest.ContentType = resp.Header.Get("Content-Type")
			urlRequest.LastAttemptAt = time.Now().Unix()
			urlRequest.Save(t.db)
		}
		return resp, err
	}

	// Cache it
	buf, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	urlRequest.Status = int64(resp.StatusCode)
	urlRequest.ContentType = resp.Header.Get("Content-Type")
	urlRequest.LastAttemptAt = time.Now().Unix()
	urlRequest.Save(t.db)

	err = t.cache.Write(k, buf)
	if err != nil {
		return nil, err
	}

	// Return the response
	return resp, err
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, 1500*time.Millisecond)
}

func NewClient(cacheDir string, db *sql.DB, sd *statsd.Client) (*Client, error) {

	limiter, err := NewRateLimiter(50, 2, 50_000)
	if err != nil {
		return nil, err
	}

	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: false,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Duration(1500) * time.Millisecond,
				}
				return d.DialContext(ctx, "tcp", "1.1.1.1:53")
			},
		},
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}

	cache := diskv.New(diskv.Options{
		BasePath: cacheDir,
		// CacheSizeMax:      1024 * 1024,
		// CacheSizeMax:      10_737_418_240, // 1 GB
		AdvancedTransform: AdvancedTransformExample,
		InverseTransform:  InverseTransformExample,
	})

	transport := &retryableTransport{
		transport: &http.Transport{
			Dial:                  dialTimeout,
			ResponseHeaderTimeout: 6 * time.Second,
			MaxIdleConns:          100,
			MaxConnsPerHost:       100,
			MaxIdleConnsPerHost:   100,
			DialContext:           dialContext,
		},
		limiter: limiter,
		cache:   cache,
		db:      db,
		sd:      sd,
	}

	httpC := http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	client := Client{
		httpClient: &httpC,
		db:         db,
		sd:         sd,
	}

	return &client, nil
}

func (c *Client) GetWithSafety(u string) (*URLRequest, error) {

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	urlRequest := &URLRequest{}
	urlRequest.Url = req.URL.String()
	urlRequest.Domain = req.URL.Hostname()
	urlRequest.LastAttemptAt = time.Now().Unix()
	urlRequest.Status = 900 // Unknown

	err = urlRequest.Save(c.db)
	if err != nil {
		return nil, err
	}

	// Check it's a valid domain
	for _, bad := range BANNED_URLS {
		if req.URL.Hostname() == bad {
			urlRequest.Status = 800 // Unknown
			err = urlRequest.Save(c.db)
			if err != nil {
				return nil, err
			}
			return urlRequest, nil
		}
	}

	if strings.HasSuffix(urlRequest.Url, ".mp4") {
		urlRequest.Status = 801 // Bad
		err = urlRequest.Save(c.db)
		if err != nil {
			return nil, err
		}
		return urlRequest, nil
	}
	if strings.HasSuffix(urlRequest.Url, ".mp3") {
		urlRequest.Status = 801 // Bad
		err = urlRequest.Save(c.db)
		if err != nil {
			return nil, err
		}
		return urlRequest, nil
	}
	if strings.HasSuffix(urlRequest.Url, ".pdf") {
		urlRequest.Status = 801 // Bad
		err = urlRequest.Save(c.db)
		if err != nil {
			return nil, err
		}
		return urlRequest, nil
	}
	if !strings.HasPrefix(urlRequest.Url, "http") {
		urlRequest.Status = 801 // Bad
		err = urlRequest.Save(c.db)
		if err != nil {
			return nil, err
		}
		return urlRequest, nil
	}

	// Now do a head request to check type

	if !c.headCheck(urlRequest.Url, 0) {
		urlRequest.Status = 802 // Bad
		err = urlRequest.Save(c.db)
		if err != nil {
			return nil, err
		}
		return urlRequest, nil
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return urlRequest, err
	}

	urlRequest.Response = resp

	urlRequest.Status = int64(resp.StatusCode)

	contentType := resp.Header.Get("Content-Type")

	urlRequest.ContentType = contentType

	err = urlRequest.Save(c.db)
	if err != nil {
		return nil, err
	}

	return urlRequest, nil

}

func (c *Client) headCheck(u string, n int) bool {
	if n > 2 {
		return false
	}

	head, err := c.httpClient.Head(u)
	if err != nil {
		return false
	}
	defer head.Body.Close()

	if head.StatusCode >= 400 {
		return false
	}

	loc := head.Header.Get("Location")
	if loc != "" {
		return c.headCheck(loc, n+1)
	}

	isGoodContentType := false

	goodTypes := []string{"text/html", "application/rss+xml", "application/atom+xml", "text/xml", "application/xml", "application/rss"}

	for _, typ := range goodTypes {
		if strings.Contains(head.Header.Get("Content-Type"), typ) {
			isGoodContentType = true

		}
	}

	if !isGoodContentType {
		return false
	}

	size := head.Header.Get("Content-Length")

	if size != "" {
		sizeInt, err := strconv.ParseInt(size, 10, 64)
		if err != nil {
			return false
		}
		if sizeInt > 200000 { // 100 kb
			return false
		}
	}

	return true
}

func (c *Client) Get(u string) (*http.Response, error) {
	return c.httpClient.Get(u)
}

func (c *Client) Head(u string) (*http.Response, error) {
	return c.httpClient.Head(u)
}

// func (c *Client) GetFromDisk(u string) (io.ReadCloser, error) {

// 	req, err := http.NewRequest("GET", u, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Check our cache
// 	k, err := getHTMLKey(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	c.httpClient.

// 	diskStream, _ := c.httpClient..ReadStream(k, false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if diskStream != nil {
// 		defer diskStream.Close()
// 		data, err := io.ReadAll(diskStream)

// 		if err != nil {
// 			return nil, err
// 		}
// 		buf := bytes.NewBuffer(data)
// 		bufReader := bufio.NewReader(buf)
// 		return http.ReadResponse(bufReader, req)
// 	}
// }
