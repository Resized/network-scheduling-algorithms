package leaky_bucket

import (
	"testing"
	"time"
)

func TestLeakyBucketPayloadTooLarge(t *testing.T) {
	var lb = LeakyBucket{}
	lb.Init(100, 20)
	err := lb.Add(101, time.Now())
	assertError(true, err, t)
}

func TestLeakyBucketNotEnoughRoom(t *testing.T) {
	var lb = LeakyBucket{}
	lb.Init(100, 20)
	lb.current = 90
	err := lb.Add(11, time.Now())
	assertError(true, err, t)
}

func TestLeakyBucketEnoughRoom(t *testing.T) {
	var lb = LeakyBucket{}
	lb.Init(100, 20)
	lb.current = 90
	err := lb.Add(10, time.Now())
	assertError(false, err, t)
}

func TestLeakyBucketJustEnoughRoom(t *testing.T) {
	var lb = LeakyBucket{}
	lb.Init(100, 20)
	lb.current = lb.maxCapacity
	now := time.Now()
	lb.lastTime = now.Add(-time.Millisecond * 1000)
	err := lb.Add(20, now)
	assertError(false, err, t)
}

func TestLeakyBucketJustNotEnoughRoom(t *testing.T) {
	var lb = LeakyBucket{}
	lb.Init(100, 20)
	lb.current = lb.maxCapacity
	now := time.Now()
	lb.lastTime = now.Add(-time.Millisecond * 1000)
	err := lb.Add(21, now)
	assertError(true, err, t)
}

func assertError(expectedError bool, err error, t *testing.T) {
	if err != nil && !expectedError {
		t.Errorf("unexpected error: %v", err)
	}
	if err == nil && expectedError {
		t.Error("expected error but got none")
	}
}
