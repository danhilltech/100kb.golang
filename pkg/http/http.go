package http

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"regexp"
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
	existing, err := getURLRequestFromDB(req.URL.String(), req.Method, t.db)
	if err != nil {
		return nil, err
	}

	var urlRequest *URLRequest

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

	checkDisk := true

	if existing != nil {
		urlRequest = existing
	} else {
		// checkDisk = false

		urlRequest = &URLRequest{
			Url:           req.URL.String(),
			LastAttemptAt: time.Now().Unix(),
			Method:        req.Method,
			Domain:        req.URL.Hostname(),
		}
	}
	defer urlRequest.Save(t.db)

	// It failed last time, and we tried in last 24 hours
	if urlRequest.Status >= 400 && urlRequest.LastAttemptAt > (time.Now().Unix()-60*60*24) {
		return nil, ErrFailingRemote
	}

	// Within 2 hours, always use disk

	// TODO here check if its a feed or an article
	if existing != nil && urlRequest.LastAttemptAt < (time.Now().Unix()-60*60*2) {
		checkDisk = false
	}

	// If cloudflare limited us
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

	if urlRequest.Etag != "" {
		req.Header.Set("if-none-match", urlRequest.Etag)
	} else if urlRequest.LastModified != "" {
		req.Header.Set("if-modified-since", urlRequest.LastModified)
	}

	sleep := 0
	servingHost := req.URL.Hostname()

	for i := 0; ; i++ {

		if i > 3 {
			return nil, ErrTooManyRetries
		}

		if sleep > 0 {
			fmt.Println("Sleeping for", sleep, i, req.URL.String(), servingHost)
			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
		ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancelFunc()

		err = t.limiter.Wait(ctx)
		if err != nil {
			fmt.Println("Global limit")
			sleep += 200
			continue
		}

		err = t.limiter.WaitHost(ctx, servingHost)
		if err != nil {
			fmt.Println("Host limit", servingHost)
			sleep += 500
			continue
		}

		// debugDump, err := httputil.DumpRequestOut(req, false)
		// if err != nil {
		// 	return nil, err
		// }
		// fmt.Printf("%s\n\n", string(debugDump))

		fmt.Println("🏃", req.Method, req.URL.String())

		// Send the request
		resp, err := t.transport.RoundTrip(req)
		if err != nil {

			if isFatalError(err) {
				fmt.Println("🛜  fatal", err)
				urlRequest.Status = 509
				return nil, ErrFailingRemote
			}
			if err == context.DeadlineExceeded {
				fmt.Println("🛜  deadline", err)
				urlRequest.Status = 509
				sleep += 2000
			}
			if err == context.Canceled {
				fmt.Println("🛜  canceled", err)
				urlRequest.Status = 509
				sleep += 2000
			}

			fmt.Println("🛜 other", err, req.URL.String())

			continue
		}

		urlRequest.Status = int64(resp.StatusCode)
		urlRequest.ContentType = resp.Header.Get("Content-Type")

		switch resp.StatusCode {

		case http.StatusTooManyRequests:
			{
				if resp.Header.Get("x-served-by") != "" {
					servingHost = resp.Header.Get("x-served-by")
				}

				if s, ok := parseRetryAfterHeader(resp.Header["Retry-After"]); ok {

					sleep += int(s.Milliseconds())
					continue
				}

				sleep *= 2
				continue
			}
		case http.StatusNotFound:
			{

				return nil, ErrNotFound
			}
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

		urlRequest.Etag = resp.Header.Get("Etag")
		urlRequest.LastModified = resp.Header.Get("Last-modified")
		urlRequest.LastAttemptAt = time.Now().Unix()
		urlRequest.DiskPath = k

		// Return the response
		return resp, err
	}

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
		limiter: limiter,
		cache:   cache,
		db:      db,
		sd:      sd,
	}

	httpC := http.Client{
		Transport: transport,
		Timeout:   20 * time.Second,
	}

	client := Client{
		httpClient: &httpC,
		db:         db,
		sd:         sd,
	}

	return &client, nil
}

func (c *Client) GetWithSafety(u string) (*http.Response, error) {

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	// Now do a head request to check type

	err = c.headCheck(req.URL.String(), 0)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)

}

func (c *Client) headCheck(u string, n int) error {
	if n > 2 {
		return ErrTooManyRedirects
	}

	head, err := c.httpClient.Head(u)
	if err != nil {
		return err
	}
	defer head.Body.Close()

	if head.StatusCode >= 400 {
		return Err400GreaterError
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
		return ErrBadContentType
	}

	size := head.Header.Get("Content-Length")

	if size != "" {
		sizeInt, err := strconv.ParseInt(size, 10, 64)
		if err != nil {
			return ErrTooLarge
		}
		if sizeInt > 200000 { // 100 kb
			return ErrTooLarge
		}
	}

	return nil
}

func (c *Client) Get(u string) (*http.Response, error) {
	return c.httpClient.Get(u)
}

func (c *Client) Head(u string) (*http.Response, error) {
	return c.httpClient.Head(u)
}

var ( // A regular expression to match the error returned by net/http when the
	// configured number of redirects is exhausted. This error isn't typed
	// specifically so we resort to matching on the error string.
	redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)

	// A regular expression to match the error returned by net/http when the
	// scheme specified in the URL is invalid. This error isn't typed
	// specifically so we resort to matching on the error string.
	schemeErrorRe = regexp.MustCompile(`unsupported protocol scheme`)

	// A regular expression to match the error returned by net/http when a
	// request header or value is invalid. This error isn't typed
	// specifically so we resort to matching on the error string.
	invalidHeaderErrorRe = regexp.MustCompile(`invalid header`)

	// A regular expression to match the error returned by net/http when the
	// TLS certificate is not trusted. This error isn't typed
	// specifically so we resort to matching on the error string.
	notTrustedErrorRe = regexp.MustCompile(`certificate is not trusted`)

	tlsUnrecognizedNameRe = regexp.MustCompile(`unrecognized name`)
	refusedRe             = regexp.MustCompile(`connection refused`)
	noSuchHostRe          = regexp.MustCompile(`no such host`)
)

func isFatalError(err error) bool {

	_, ok := err.(*tls.CertificateVerificationError)
	if ok {
		return true
	}

	// Don't retry if the error was due to too many redirects.
	if redirectsErrorRe.MatchString(err.Error()) {
		return true
	}

	// Don't retry if the error was due to an invalid protocol scheme.
	if schemeErrorRe.MatchString(err.Error()) {
		return true
	}

	// Don't retry if the error was due to an invalid header.
	if invalidHeaderErrorRe.MatchString(err.Error()) {
		return true
	}

	// Don't retry if the error was due to TLS cert verification failure.
	if notTrustedErrorRe.MatchString(err.Error()) {
		return true
	}

	if tlsUnrecognizedNameRe.MatchString(err.Error()) {
		return true
	}
	if refusedRe.MatchString(err.Error()) {
		return true
	}
	if noSuchHostRe.MatchString(err.Error()) {
		return true
	}

	return false
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
