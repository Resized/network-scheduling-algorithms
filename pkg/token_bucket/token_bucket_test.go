package token_bucket

import (
	"testing"
	"time"
)

func TestTokenBucketPayloadTooLarge(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	err := tb.Take(101, time.Now())
	assertError(true, err, t)
}

func TestTokenBucketNotEnoughTokens(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	tb.currentTokens = 0
	err := tb.Take(20, time.Now())
	assertError(true, err, t)
}

func TestTokenBucketEnoughTokens(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	tb.currentTokens = 60
	err := tb.Take(20, time.Now())
	assertError(false, err, t)
}

func TestTokenBucketEnoughTokensMaxPayloadSize(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	tb.currentTokens = 100
	err := tb.Take(100, time.Now())
	assertError(false, err, t)
}

func TestTokenBucketNotEnoughTokensMaxPayloadSize(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	tb.currentTokens = 99
	err := tb.Take(100, time.Now())
	assertError(true, err, t)
}

func TestTokenBucketJustEnoughTokens(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	tb.currentTokens = 0
	tb.lastTime = time.Now().Add(-time.Millisecond * 1000)
	err := tb.Take(20, time.Now())
	assertError(false, err, t)
}

func TestTokenBucketJustNotEnoughTokens(t *testing.T) {
	var tb = TokenBucket{}
	tb.Init(100, 20)
	tb.currentTokens = 0
	tb.lastTime = time.Now().Add(-time.Millisecond * 999)
	err := tb.Take(20, time.Now())
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

func BenchmarkTokenBucket(b *testing.B) {
	tb := TokenBucket{}
	tb.Init(100, 20)
	payloadSize := 10
	for i := 0; i < b.N; i++ {
		_ = tb.Take(payloadSize, time.Now())
	}
}
