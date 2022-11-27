package token_bucket

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestTokenBucketPayloadTooLarge(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	err := tb.Take(101, time.Now())
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "payload too large"))
}

func TestTokenBucketNotEnoughTokens(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	tb.currentTokens = 0
	err := tb.Take(20, time.Now())
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "not enough tokens"))
}

func TestTokenBucketEnoughTokens(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	tb.currentTokens = 60
	err := tb.Take(20, time.Now())
	assert.NoError(t, err)
}

func TestTokenBucketEnoughTokensMaxPayloadSize(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	tb.currentTokens = 100
	err := tb.Take(100, time.Now())
	assert.NoError(t, err)
}

func TestTokenBucketNotEnoughTokensMaxPayloadSize(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	tb.currentTokens = 99
	err := tb.Take(100, time.Now())
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "not enough tokens"))
}

func TestTokenBucketJustEnoughTokens(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	tb.currentTokens = 0
	now := time.Now()
	tb.lastTime = now.Add(-time.Millisecond * 1000)
	err := tb.Take(20, now)
	assert.NoError(t, err)
}

func TestTokenBucketJustNotEnoughTokens(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	tb.currentTokens = 0
	now := time.Now()
	tb.lastTime = now.Add(-time.Millisecond * 999)
	err := tb.Take(20, now)
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "not enough tokens"))
}

func BenchmarkTokenBucket(b *testing.B) {
	tb := TokenBucket{}
	tb.Init(100, 20)
	payloadSize := 10
	for i := 0; i < b.N; i++ {
		_ = tb.Take(payloadSize, time.Now())
	}
}
