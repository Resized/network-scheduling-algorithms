package token_bucket

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	maxCapacity   float64   // max capacity of the bucket
	currentTokens float64   // current number of tokens in the bucket
	fillRate      float64   // tokens per second
	lastTime      time.Time // last time the bucket was updated
	mu            sync.Mutex
}

func (tb *TokenBucket) Take(payloadSize int, t time.Time) error {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.calcTokens(t)
	if float64(payloadSize) > tb.maxCapacity {
		return fmt.Errorf("payload too large for bucket, size: %v, max capacity: %.3v", payloadSize, tb.maxCapacity)
	}
	if float64(payloadSize) > tb.currentTokens {
		return fmt.Errorf("not enough tokens in bucket, need %v, have %.3v", payloadSize, tb.currentTokens)
	}
	tb.currentTokens -= float64(payloadSize)
	return nil
}

func (tb *TokenBucket) calcTokens(t time.Time) {
	elapsed := t.Sub(tb.lastTime)
	tb.currentTokens += tb.fillRate * float64(elapsed.Milliseconds()) / 1000.0
	tb.lastTime = t
	if tb.currentTokens > tb.maxCapacity {
		tb.currentTokens = tb.maxCapacity
	}
}

func (tb *TokenBucket) Init(maxCapacity float64, fillRate float64) {
	tb.maxCapacity = maxCapacity
	tb.currentTokens = maxCapacity
	tb.fillRate = fillRate
	tb.lastTime = time.Now()
}

func (tb *TokenBucket) GetCurrentTokens() float64 {
	return tb.currentTokens
}
