package http

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/peterbourgon/diskv/v3"
)

const RetryCount = 2

type Client struct {
	httpClient *http.Client
}

type retryableTransport struct {
	transport http.RoundTripper
	limiter   *Limiter
	cache     *diskv.Diskv
}

func (t *retryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	// Do we have a context?

	var err error

	// Check our cache
	k, err := getHTMLKey(req)
	if err != nil {
		return nil, err
	}
	diskStream, _ := t.cache.ReadStream(k, false)
	if err != nil {
		return nil, err
	}
	if diskStream != nil {
		defer diskStream.Close()
		data, err := io.ReadAll(diskStream)

		if err != nil {
			return nil, err
		}
		buf := bytes.NewBuffer(data)
		bufReader := bufio.NewReader(buf)
		return http.ReadResponse(bufReader, req)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	t.limiter.Wait(ctx)

	t.limiter.WaitHost(ctx, req.Host)

	// Send the request
	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

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
	return net.DialTimeout(network, addr, 1500*time.Millisecond)
}

func NewClient(cacheDir string) (*Client, error) {

	limiter, err := NewRateLimiter(100, 2, 30_000)
	if err != nil {
		return nil, err
	}

	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Duration(1500) * time.Millisecond,
				}
				return d.DialContext(ctx, "udp", "1.1.1.1:53")
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
	}

	httpC := http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	client := Client{
		httpClient: &httpC,
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
	urlRequest.Status = "9xx" // Unknown

	// Check it's a valid domain
	for _, bad := range BANNED_URLS {
		if req.URL.Hostname() == bad {
			urlRequest.Status = "8xx" // Unknown
			return urlRequest, nil
		}
	}

	if strings.HasSuffix(urlRequest.Url, ".mp4") {
		urlRequest.Status = "8xx" // Bad
		return urlRequest, nil
	}
	if strings.HasSuffix(urlRequest.Url, ".mp3") {
		urlRequest.Status = "8xx" // Bad
		return urlRequest, nil
	}
	if strings.HasSuffix(urlRequest.Url, ".pdf") {
		urlRequest.Status = "8xx" // Bad
		return urlRequest, nil
	}
	if !strings.HasPrefix(urlRequest.Url, "http") {
		urlRequest.Status = "8xx" // Bad
		return urlRequest, nil
	}

	// Now do a head request to check type

	if !c.headCheck(urlRequest.Url, 0) {
		urlRequest.Status = "8xx" // Bad
		return urlRequest, nil
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return urlRequest, err
	}

	urlRequest.Response = resp

	if resp.StatusCode < 400 {
		urlRequest.Status = "200"
	} else if resp.StatusCode < 500 {
		urlRequest.Status = "4xx"
	} else if resp.StatusCode < 600 {
		urlRequest.Status = "5xx"
	}

	contentType := resp.Header.Get("Content-Type")

	urlRequest.ContentType = contentType

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

	if !strings.Contains(head.Header.Get("Content-Type"), "text/html") && !strings.Contains(head.Header.Get("Content-Type"), "application/rss") {
		return false
	}

	size := head.Header.Get("Content-Length")

	if size != "" {
		sizeInt, err := strconv.ParseInt(size, 10, 64)
		if err != nil {
			return false
		}
		if sizeInt > 100000 { // 100 kb
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
