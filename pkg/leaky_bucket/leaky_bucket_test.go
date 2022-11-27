package leaky_bucket

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLeakyBucketPayloadTooLarge(t *testing.T) {
	var lb = LeakyBucket{}
	lb.Init(100, 20)
	err := lb.Add(101, time.Now())
	assert.Error(t, err)
}

func TestLeakyBucketNotEnoughRoom(t *testing.T) {
	var lb = LeakyBucket{}
	lb.Init(100, 20)
	lb.current = 90
	err := lb.Add(11, time.Now())
	assert.Error(t, err)
}

func TestLeakyBucketEnoughRoom(t *testing.T) {
	var lb = LeakyBucket{}
	lb.Init(100, 20)
	lb.current = 90
	err := lb.Add(10, time.Now())
	assert.NoError(t, err)
}

func TestLeakyBucketJustEnoughRoom(t *testing.T) {
	var lb = LeakyBucket{}
	lb.Init(100, 20)
	lb.current = lb.maxCapacity
	now := time.Now()
	lb.lastTime = now.Add(-time.Millisecond * 1000)
	err := lb.Add(20, now)
	assert.NoError(t, err)
}

func TestLeakyBucketJustNotEnoughRoom(t *testing.T) {
	var lb = LeakyBucket{}
	lb.Init(100, 20)
	lb.current = lb.maxCapacity
	now := time.Now()
	lb.lastTime = now.Add(-time.Millisecond * 1000)
	err := lb.Add(21, now)
	assert.Error(t, err)
}
