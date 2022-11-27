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

func (tb *TokenBucket) calcTokens(t time.Time) {
	if tb.lastTime.IsZero() {
		tb.lastTime = t
	}
	elapsed := t.Sub(tb.lastTime)
	tb.currentTokens += (tb.fillRate / 1000.0) * float64(elapsed.Milliseconds())
	tb.lastTime = t
	if tb.currentTokens > tb.maxCapacity {
		tb.currentTokens = tb.maxCapacity
	}
}

func (tb *TokenBucket) sendPacket(p Packet, t time.Time) error {
	tb.calcTokens(t)
	if float64(p.size) > tb.maxCapacity {
		return fmt.Errorf("packet too large for bucket, size: %v, max capacity: %.3v", p.size, tb.maxCapacity)
	}
	if float64(p.size) > tb.currentTokens {
		return fmt.Errorf("not enough tokens in bucket, need %v, have %.3v", p.size, tb.currentTokens)
	}
	tb.currentTokens -= float64(p.size)
	//fmt.Printf("sending packet of size %v, current tokens: %.3v\n", p.size, tb.currentTokens)
	return nil
}

func main() {
	var p = [10]Packet{}
	for i := 0; i < 10; i++ {
		p[i].size = i + 30
	}
	var tb = TokenBucket{maxCapacity: 100, currentTokens: 0, fillRate: 20, lastTime: time.Now()}
	for {
		for i := 0; i < 10; i++ {
			err := tb.sendPacket(p[i], time.Now())
			if err != nil {
				fmt.Printf("error sending packet: %v\n", err)
			}
			time.Sleep(time.Second)
		}
	}
}
