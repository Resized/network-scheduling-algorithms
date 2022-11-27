package leaky_bucket

import (
	"fmt"
	"time"
)

type LeakyBucket struct {
	maxCapacity float64
	current     float64
	leakRate    float64
	lastTime    time.Time
}

func (lb *LeakyBucket) Add(payloadSize int, t time.Time, logging ...bool) error {
	lb.calcLeak(t)
	if float64(payloadSize) > lb.maxCapacity {
		return fmt.Errorf("payload too large for bucket, size: %v, max capacity: %.3v", payloadSize, lb.maxCapacity)
	}
	if lb.current+float64(payloadSize) > lb.maxCapacity {
		return fmt.Errorf("no room left for payload in bucket, size: %v, room left: %.3v", payloadSize, lb.maxCapacity-lb.current)
	}
	lb.current += float64(payloadSize)
	if len(logging) > 0 && logging[0] {
		fmt.Printf("accepting payload of size %v, current bucket: %.3v out of %v\n", payloadSize, lb.current, lb.maxCapacity)
	}
	return nil
}

func (lb *LeakyBucket) calcLeak(t time.Time) {
	if lb.lastTime.IsZero() {
		lb.lastTime = t
	}
	elapsed := t.Sub(lb.lastTime)
	lb.current -= lb.leakRate * float64(elapsed.Milliseconds()) / 1000.0
	lb.lastTime = t
	if lb.current < 0 {
		lb.current = 0
	}
}

func (lb *LeakyBucket) Init(maxCapacity float64, leakRate float64) {
	lb.maxCapacity = maxCapacity
	lb.leakRate = leakRate
	lb.current = 0
}
