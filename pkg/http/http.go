package http

import (
	"bytes"
	"context"
	"io"
	"math"
	"net"
	"net/http"
	"time"
)

const RetryCount = 2

type retryableTransport struct {
	transport http.RoundTripper
	limiter   *Limiter
}

func backoff(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

func shouldRetry(err error, resp *http.Response) bool {
	if err != nil {
		// fmt.Println(err)
		return false
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return true
	}

	return false
}

func drainBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func (t *retryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request body
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	t.limiter.Wait(ctx)

	t.limiter.WaitHost(ctx, req.Host)

	// Send the request
	resp, err := t.transport.RoundTrip(req)

	// Retry logic
	retries := 0
	for shouldRetry(err, resp) && retries < RetryCount {
		// Wait for the specified backoff period
		// fmt.Println("retrying", req.URL)
		time.Sleep(backoff(retries))

		// We're going to retry, consume any response to reuse the connection.
		drainBody(resp)

		// Clone the request body again
		if req.Body != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Retry the request
		resp, err = t.transport.RoundTrip(req)

		retries++
	}

	// Return the response
	return resp, err
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, 1500*time.Millisecond)
}

func NewRetryableClient() *http.Client {

	limiter, err := NewRateLimiter(100, 2, 30_000)
	if err != nil {
		panic(err)
	}

	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Duration(1000) * time.Millisecond,
				}
				return d.DialContext(ctx, "udp", "1.1.1.1:53")
			},
		},
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}

	transport := &retryableTransport{
		transport: &http.Transport{
			Dial:                  dialTimeout,
			ResponseHeaderTimeout: 3 * time.Second,
			MaxIdleConns:          100,
			MaxConnsPerHost:       100,
			MaxIdleConnsPerHost:   100,
			DialContext:           dialContext,
		},
		limiter: limiter,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}
}
