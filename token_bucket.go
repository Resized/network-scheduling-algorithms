package main

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

type Packet struct {
	size int
}

func (tb *TokenBucket) calcTokens() {
	now := time.Now()
	tb.currentTokens += (tb.fillRate / 1000.0) * float64(now.Sub(tb.lastTime).Milliseconds())
	tb.lastTime = now
	if tb.currentTokens > tb.maxCapacity {
		tb.currentTokens = tb.maxCapacity
	}
}

func (tb *TokenBucket) sendPacket(p Packet) error {
	tb.calcTokens()
	if float64(p.size) > tb.maxCapacity {
		return fmt.Errorf("packet too large for bucket, size: %v, max capacity: %.3v", p.size, tb.maxCapacity)
	}
	if float64(p.size) > tb.currentTokens {
		return fmt.Errorf("not enough tokens in bucket, need %v, have %.3v", p.size, tb.currentTokens)
	}
	fmt.Printf("sending packet of size %v, current tokens: %.3v\n", p.size, tb.currentTokens)
	tb.currentTokens -= float64(p.size)
	return nil
}
