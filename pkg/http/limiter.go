package http

import (
	"context"
	"fmt"

	lru "github.com/opencoff/golang-lru"
	"golang.org/x/time/rate"
)

// Limiter controls how frequently events are allowed to happen globally or
// per-host. It uses a token-bucket limiter for the global limit and instantiates
// a token-bucket limiter for every unique host. The number of per-host limiters
// is limited to an upper bound ("cache size").
//
// A negative rate limit means "no limit" and a zero rate limit means "Infinite".
type Limiter struct {
	// Global rate limiter; thread-safe
	gl *rate.Limiter

	// Per-host limiter organized as an LRU cache; thread-safe
	h lru.Cache

	// per host rate limit (qps)
	p rate.Limit
	g rate.Limit

	// burst rate for per-host
	b int

	cache int
}

// Create a new token bucket rate limiter that limits globally at 'g'  requests/sec
// and per-host at 'p' requests/sec; It remembers the rate of the 'cachesize' most
// recent hosts (and their limits). The burst rates are pre-configured to be:
// Global burst limit: 3 * b; Per host burst limit:  2 * p
func NewRateLimiter(g, p, cachesize int) (*Limiter, error) {
	l, err := lru.New2Q(cachesize)
	if err != nil {
		return nil, fmt.Errorf("ratelimit: can't create LRU cache: %s", err)
	}

	b := 2 * p
	if b < 0 {
		b = 0
	}

	gl := limit(g)
	pl := limit(p)

	r := &Limiter{
		gl:    rate.NewLimiter(gl, 3*g),
		h:     l,
		p:     pl,
		g:     gl,
		b:     b,
		cache: cachesize,
	}

	return r, nil
}

// Wait blocks until the ratelimiter permits the configured global rate limit.
// It returns an error if the burst exceeds the configured limit or the
// context is cancelled.
func (r *Limiter) Wait(ctx context.Context) error {
	return r.gl.Wait(ctx)
}

// WaitHost blocks until the ratelimiter permits the configured per-host
// rate limit from host 'a'.
// It returns an error if the burst exceeds the configured limit or the
// context is cancelled.
func (r *Limiter) WaitHost(ctx context.Context, a string) error {
	rl := r.getRL(a)
	return rl.Wait(ctx)
}

// Allow returns true if the global rate limit can consume 1 token and
// false otherwise. Use this if you intend to drop/skip events that exceed
// a configured global rate limit, otherwise, use Wait().
func (r *Limiter) Allow() bool {
	return r.gl.Allow()
}

// AllowHost returns true if the per-host rate limit for host 'a' can consume
// 1 token and false otherwise. Use this if you intend to drop/skip events
// that exceed a configured global rate limit, otherwise, use WaitHost().
func (r *Limiter) AllowHost(a string) bool {
	rl := r.getRL(a)
	return rl.Allow()
}

// String returns a printable representation of the limiter
func (r Limiter) String() string {
	return fmt.Sprintf("ratelimiter: Global %0.2f rps, Per host %4.2 rps, LRU cache %d entries",
		r.g, r.p, r.cache)
}

// get or create a new per-host rate limiter.
// this function evicts the least used limiter from the LRU cache
func (r *Limiter) getRL(k string) *rate.Limiter {
	v, _ := r.h.Probe(k, func(k interface{}) interface{} {
		return rate.NewLimiter(r.p, r.b)
	})

	rl, ok := v.(*rate.Limiter)
	if !ok {
		panic(fmt.Sprintf("ratelimiter: bad type %t for host %s in per-host limiter", v, k))
	}
	return rl
}

func limit(r int) rate.Limit {
	var g rate.Limit

	switch {
	case r < 0:
		g = rate.Inf
	case r == 0:
		g = 0.0
	default:
		g = rate.Limit(r)
	}

	return g
}
