package http

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"fmt"
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

	cache   *diskv.Diskv
	limiter *Limiter
	sd      *statsd.Client
}

var ErrFailingRemote = fmt.Errorf("remote is known to be failing")
var Err400GreaterError = fmt.Errorf("remote is currently failing")

// var ErrBadHead = fmt.Errorf("url failed head check")
var ErrBannedUrl = fmt.Errorf("url banned")
var ErrBadFormat = fmt.Errorf("url has bad extension")
var ErrBadContentType = fmt.Errorf("url has bad content type")
var ErrTooLarge = fmt.Errorf("url is too large")
var ErrTooManyRedirects = fmt.Errorf("too many redirects")
var ErrGlobalLimitHit = fmt.Errorf("global limit hit")
var ErrHostLimitHit = fmt.Errorf("per host limit hit")
var ErrTooManyRetries = fmt.Errorf("too many retries here")
var ErrNotFound = fmt.Errorf("url not found")
var ErrNotInCache = fmt.Errorf("url not in cache")

func (t *retryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	// Do we have a context?

	req.Header.Set("User-Agent", "curl/8.4.0")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15")
	req.Header.Set("Accepts", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("DNT", "1")

	t.sd.Incr("http.roundtrip.start", 1)

	var err error

	k, err := getHTMLKey(req)
	if err != nil {
		return nil, err
	}

	t.sd.Incr("http.roundtrip.checkingDisk", 1)
	// Check our cache

	diskStream, _ := t.cache.ReadStream(k, false)

	var cached *http.Response

	if diskStream != nil {
		t.sd.Incr("http.roundtrip.diskCacheHit", 1)
		// defer diskStream.Close()
		// data, err := io.ReadAll(diskStream)

		// if err != nil {
		// 	return nil, err
		// }
		// buf := bytes.NewBuffer(data)
		bufReader := bufio.NewReader(diskStream)

		// todo check how old the response is based on the code

		cached, err = http.ReadResponse(bufReader, req)
		if err != nil {
			return nil, err
		}

		cachedHeader := cached.Header.Get("x-100kb-cached-at")

		lastAttemptAt, _ := strconv.Atoi(cachedHeader)

		if cached.StatusCode != http.StatusTooManyRequests {
			cached.Header.Set("x-100kb-from-cache", "1")
			// Todo, retry some errors?

			if int64(lastAttemptAt) > (time.Now().Unix()-60*60*24*3) && cached.StatusCode == http.StatusTeapot {
				return cached, nil
			}

			if cached.StatusCode >= 400 {
				return cached, nil
			}

			if int64(lastAttemptAt) > (time.Now().Unix()-60*60*24*7) && strings.Contains(cached.Header.Get("Content-Type"), "text/html") {
				return cached, nil
			}

			if cachedHeader == "" || int64(lastAttemptAt) > (time.Now().Unix()-60*60*24) {
				return cached, nil
			}
		}

	}

	t.sd.Incr("http.roundtrip.diskCacheMiss", 1)

	if cached != nil {
		if cached.Header.Get("etag") != "" {
			req.Header.Set("if-none-match", cached.Header.Get("etag"))
		} else if cached.Header.Get("Last-Modified") != "" {
			req.Header.Set("if-modified-since", cached.Header.Get("Last-Modified"))
		}
	}

	servingHosts := strings.Split(req.URL.Hostname(), ".")

	servingHost := strings.Join(servingHosts[len(servingHosts)-2:], ".")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancelFunc()

	Response429 := &http.Response{
		Status:     "FAILED",
		StatusCode: http.StatusTooManyRequests,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       io.NopCloser(bytes.NewBufferString("")),
		Header:     make(http.Header),
	}

	err = t.limiter.Wait(ctx)
	if err != nil {
		return Response429, nil
	}

	err = t.limiter.WaitHost(ctx, servingHost)
	if err != nil {
		return Response429, nil
	}

	// debugDump, err := httputil.DumpRequestOut(req, false)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("%s\n\n", string(debugDump))

	cacheState := "NO CACHE"
	cacheStatus := 0
	if cached != nil {
		cacheState = "CACHE"
		cacheStatus = cached.StatusCode
	}

	fmt.Println("🏃", req.Method, cacheState, cacheStatus, req.URL.String())

	// Send the request
	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		fmt.Println(err)
		failedResponse := &http.Response{
			Status:     "FAILED",
			StatusCode: http.StatusTeapot,
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       io.NopCloser(bytes.NewBufferString("")),
			Header:     make(http.Header),
		}

		failedResponse.Header.Set("x-100kb-cached-at", fmt.Sprintf("%d", time.Now().Unix()))

		// Cache it
		buf, err := httputil.DumpResponse(failedResponse, true)
		if err != nil {
			return nil, err
		}
		err = t.cache.Write(k, buf)
		if err != nil {
			return nil, err
		}

		return failedResponse, err

	}

	if resp.StatusCode == http.StatusNotModified {
		cached.StatusCode = http.StatusOK
		return cached, nil
	}

	resp.Header.Set("x-100kb-cached-at", fmt.Sprintf("%d", time.Now().Unix()))

	// Cache it
	buf, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	err = t.cache.Write(k, buf)
	if err != nil {
		return nil, err
	}

	// Return the response
	return resp, err

}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, 3000*time.Millisecond)
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
					Timeout: time.Duration(3000) * time.Millisecond,
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

		cache:   cache,
		limiter: limiter,
		sd:      sd,
	}

	httpC := http.Client{
		Transport: transport,
		Timeout:   3 * time.Second,
	}

	client := Client{
		httpClient: &httpC,
		db:         db,
		sd:         sd,
	}

	return &client, nil
}

func (c *Client) doGet(req *http.Request, attempt int) (*http.Response, error) {

	if attempt > 3 {
		return nil, ErrTooManyRetries
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusTeapot:
		{
			err = resp.Body.Close()
			if err != nil {
				return nil, err
			}
			return nil, ErrNotFound
		}
	case http.StatusTooManyRequests:
		{
			err = resp.Body.Close()
			if err != nil {
				return nil, err
			}
			// if resp.Header.Get("x-served-by") != "" {
			// 	servingHost = resp.Header.Get("x-served-by")
			// }

			if s, ok := parseRetryAfterHeader(resp.Header["Retry-After"]); ok {
				if s > 0 {
					time.Sleep(s)
				} else {
					time.Sleep(time.Duration(1100) * time.Millisecond * time.Duration(attempt+1))
				}
				return c.doGet(req, attempt+1)
			} else {
				time.Sleep(time.Duration(1100) * time.Millisecond * time.Duration(attempt+1))
				return c.doGet(req, attempt+1)
			}

		}
	case http.StatusNotFound:
		{
			err = resp.Body.Close()
			if err != nil {
				return nil, err
			}
			return nil, ErrNotFound
		}
	}

	isGoodContentType := false

	goodTypes := []string{"text/html", "application/rss+xml", "application/atom+xml", "text/xml", "application/xml", "application/rss"}

	for _, typ := range goodTypes {
		if strings.Contains(resp.Header.Get("Content-Type"), typ) {
			isGoodContentType = true
		}
	}

	if !isGoodContentType {
		resp.Body.Close()
		return nil, fmt.Errorf("%w %s %d %s", ErrBadContentType, req.URL.String(), resp.StatusCode, resp.Header.Get("Content-Type"))
	}

	size := resp.Header.Get("Content-Length")

	if size != "" {
		sizeInt, err := strconv.ParseInt(size, 10, 64)
		if err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("%w %s parseInt failed", ErrTooLarge, req.URL.String())
		}
		if sizeInt > 500000 && strings.Contains(resp.Header.Get("Content-Type"), "text/html") { // 100 kb
			resp.Body.Close()

			return nil, fmt.Errorf("%w %s", ErrTooLarge, req.URL.String())
		}
	}

	return resp, nil
}

func (c *Client) startRequest(method string, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Check it's a valid domain
	for _, bad := range BANNED_URLS {
		if req.URL.Hostname() == bad {
			return nil, ErrBannedUrl
		}
	}

	if strings.HasSuffix(req.URL.String(), ".mp4") {
		return nil, ErrBadFormat
	}
	if strings.HasSuffix(req.URL.String(), ".mp3") {
		return nil, ErrBadFormat
	}
	if strings.HasSuffix(req.URL.String(), ".pdf") {
		return nil, ErrBadFormat
	}
	if !strings.HasPrefix(req.URL.String(), "http") {
		return nil, ErrBadFormat
	}

	resp, err := c.doGet(req, 0)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) Get(u string) (*http.Response, error) {
	return c.startRequest(http.MethodGet, u)

}

func (c *Client) Head(u string) (*http.Response, error) {
	return c.startRequest(http.MethodHead, u)
}

func parseRetryAfterHeader(headers []string) (time.Duration, bool) {
	if len(headers) == 0 || headers[0] == "" {
		return 0, false
	}
	header := headers[0]
	// Retry-After: 120
	if sleep, err := strconv.ParseInt(header, 10, 64); err == nil {
		if sleep < 0 { // a negative sleep doesn't make sense
			return 0, false
		}
		return time.Second * time.Duration(sleep), true
	}

	// Retry-After: Fri, 31 Dec 1999 23:59:59 GMT
	retryTime, err := time.Parse(time.RFC1123, header)
	if err != nil {
		return 0, false
	}
	if until := retryTime.Sub(time.Now()); until > 0 {
		return until, true
	}
	// date is in the past
	return 0, true
}
