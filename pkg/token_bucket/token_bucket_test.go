package token_bucket

import (
	"testing"
	"time"
)

func TestTokenBucketTableDriven(t *testing.T) {
	type testStruct = struct {
		name                    string
		tb                      TokenBucket
		payloadSize             int
		elapsedTimeBetweenCalls time.Duration
		expectedError           bool
	}
	var tests = []testStruct{
		{
			name: "payload too large",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 100,
				fillRate:      20,
			},
			payloadSize:   101,
			expectedError: true,
		},
		{
			name: "not enough tokens",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 0,
				fillRate:      1,
			},
			payloadSize:   20,
			expectedError: true,
		},
		{
			name: "enough tokens",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 100,
				fillRate:      20,
			},
			payloadSize:   10,
			expectedError: false,
		},
		{
			name: "enough tokens with max payload size",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 100,
				fillRate:      20,
			},
			payloadSize:   100,
			expectedError: false,
		},
		{
			name: "not enough tokens with max payload size",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 99,
				fillRate:      20,
			},
			payloadSize:   100,
			expectedError: true,
		},
		{
			name: "tokens are just enough",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 0,
				fillRate:      20,
			},
			payloadSize:             20,
			elapsedTimeBetweenCalls: time.Millisecond * 1000,
			expectedError:           false,
		},
		{
			name: "tokens are just not enough",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 0,
				fillRate:      20,
			},
			payloadSize:             20,
			elapsedTimeBetweenCalls: time.Millisecond * 999,
			expectedError:           true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.tb.lastTime = time.Now().Add(-test.elapsedTimeBetweenCalls)
			err := test.tb.Take(test.payloadSize, time.Now())
			if err != nil && !test.expectedError {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && test.expectedError {
				t.Error("expected error but got none")
			}
		})
	}
}

func TestTokenBucketPayloadTooLarge(t *testing.T) {
	var tb = TokenBucket{}
	tb.InitTokenBucket(100, 20)
	err := tb.Take(101, time.Now())
	if err == nil {
		t.Error("expected error but got none")
	}
}

func TestTokenBucketNotEnoughTokens(t *testing.T) {
	var tb = TokenBucket{}
	tb.InitTokenBucket(100, 20)
	tb.currentTokens = 0
	err := tb.Take(20, time.Now())
	if err == nil {
		t.Error("expected error but got none")
	}
}

func TestTokenBucketEnoughTokens(t *testing.T) {
	var tb = TokenBucket{}
	tb.InitTokenBucket(100, 20)
	tb.currentTokens = 60
	err := tb.Take(20, time.Now())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTokenBucketEnoughTokensMaxPayloadSize(t *testing.T) {
	var tb = TokenBucket{}
	tb.InitTokenBucket(100, 20)
	tb.currentTokens = 100
	err := tb.Take(100, time.Now())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTokenBucketNotEnoughTokensMaxPayloadSize(t *testing.T) {
	var tb = TokenBucket{}
	tb.InitTokenBucket(100, 20)
	tb.currentTokens = 99
	err := tb.Take(100, time.Now())
	if err == nil {
		t.Error("expected error but got none")
	}
}

func TestTokenBucketJustEnoughTokens(t *testing.T) {
	var tb = TokenBucket{}
	tb.InitTokenBucket(100, 20)
	tb.currentTokens = 0
	tb.lastTime = time.Now().Add(-time.Millisecond * 1000)
	err := tb.Take(20, time.Now())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTokenBucketJustNotEnoughTokens(t *testing.T) {
	var tb = TokenBucket{}
	tb.InitTokenBucket(100, 20)
	tb.currentTokens = 0
	tb.lastTime = time.Now().Add(-time.Millisecond * 999)
	err := tb.Take(20, time.Now())
	if err == nil {
		t.Error("expected error but got none")
	}
}

func BenchmarkTokenBucket(b *testing.B) {
	tb := TokenBucket{}
	tb.InitTokenBucket(100, 20)
	payloadSize := 10
	for i := 0; i < b.N; i++ {
		_ = tb.Take(payloadSize, time.Now())
	}
}
