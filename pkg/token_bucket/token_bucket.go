package token_bucket

import (
	"fmt"
	"time"
)

type TokenBucket struct {
	maxCapacity   float64   // max capacity of the bucket
	currentTokens float64   // current number of tokens in the bucket
	fillRate      float64   // tokens per second
	lastTime      time.Time // last time the bucket was updated
}

func (tb *TokenBucket) calcTokens(t time.Time) {
	if tb.lastTime.IsZero() {
		tb.lastTime = t
	}
	elapsed := t.Sub(tb.lastTime)
	tb.currentTokens += tb.fillRate * float64(elapsed.Milliseconds()) / 1000.0
	tb.lastTime = t
	if tb.currentTokens > tb.maxCapacity {
		tb.currentTokens = tb.maxCapacity
	}
}

func (tb *TokenBucket) Take(payloadSize int, t time.Time, logging ...bool) error {
	tb.calcTokens(t)
	if float64(payloadSize) > tb.maxCapacity {
		return fmt.Errorf("payload too large for bucket, size: %v, max capacity: %.3v", payloadSize, tb.maxCapacity)
	}
	if float64(payloadSize) > tb.currentTokens {
		return fmt.Errorf("not enough tokens in bucket, need %v, have %.3v", payloadSize, tb.currentTokens)
	}
	tb.currentTokens -= float64(payloadSize)
	if len(logging) > 0 && logging[0] {
		fmt.Printf("sending packet of size %v, current tokens: %.3v\n", payloadSize, tb.currentTokens)
	}
	return nil
}

func (tb *TokenBucket) Init(maxCapacity float64, fillRate float64) {
	tb.maxCapacity = maxCapacity
	tb.currentTokens = maxCapacity
	tb.fillRate = fillRate
}
