package clover

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter controls request rate and concurrency for Clover API calls.
type RateLimiter struct {
	limiter    *rate.Limiter   // token bucket: 16 req/s, burst 5
	semaphore  chan struct{}   // max 5 concurrent requests
	maxRetries int
}

// NewRateLimiter creates a rate limiter configured for Clover's per-token limits.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiter:    rate.NewLimiter(rate.Limit(16), 5), // 16 req/s, burst 5
		semaphore:  make(chan struct{}, 5),             // max 5 concurrent
		maxRetries: 10,
	}
}

// Do executes the given request respecting rate limits and concurrency.
// It retries on 429 with exponential backoff.
func (r *RateLimiter) Do(client *http.Client, req *http.Request) (*http.Response, error) {
	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		// Wait for rate limiter token
		ctx := req.Context()
		if err := r.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter wait: %w", err)
		}

		// Acquire concurrent slot
		select {
		case r.semaphore <- struct{}{}:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		resp, err := client.Do(req)
		<-r.semaphore // release slot

		if err != nil {
			return nil, err
		}

		// Success or non-429 error
		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		// 429 received — handle retry
		resp.Body.Close()

		if attempt == r.maxRetries {
			return nil, fmt.Errorf("rate limited after %d retries", r.maxRetries)
		}

		waitDuration := r.calculateWait(resp, attempt)
		select {
		case <-time.After(waitDuration):
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("rate limited after %d retries", r.maxRetries)
}

// calculateWait determines how long to wait before retrying.
// It respects retry-after header if present, otherwise uses exponential backoff with jitter.
func (r *RateLimiter) calculateWait(resp *http.Response, attempt int) time.Duration {
	// Check for retry-after header (only present on concurrent limit 429s)
	if retryAfter := resp.Header.Get("retry-after"); retryAfter != "" {
		if seconds, err := strconv.Atoi(retryAfter); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}

	// Exponential backoff: 2^attempt + random jitter
	backoff := math.Pow(2, float64(attempt))
	jitter := rand.Float64() // 0.0 - 1.0
	duration := time.Duration(backoff+jitter) * time.Second

	// Cap at reasonable max to avoid excessive waits
	maxWait := 60 * time.Second
	if duration > maxWait {
		duration = maxWait
	}

	return duration
}
