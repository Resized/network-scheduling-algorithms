package main

import (
	"testing"
	"time"
)

func TestTokenBucketTableDriven(t *testing.T) {
	type testStruct = struct {
		name          string
		tb            TokenBucket
		p             Packet
		waitTime      time.Duration
		expectedError bool
	}
	var tests = []testStruct{
		{
			name: "packet too large",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 100,
				fillRate:      20,
				lastTime:      time.Now(),
			},
			p:             Packet{size: 101},
			waitTime:      0,
			expectedError: true,
		},
		{
			name: "not enough tokens",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 0,
				fillRate:      1,
				lastTime:      time.Now(),
			},
			p:             Packet{size: 20},
			waitTime:      0,
			expectedError: true,
		},
		{
			name: "enough tokens",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 100,
				fillRate:      20,
				lastTime:      time.Now(),
			},
			p:             Packet{size: 10},
			waitTime:      0,
			expectedError: false,
		},
		{
			name: "enough tokens with max packet size",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 100,
				fillRate:      20,
				lastTime:      time.Now(),
			},
			p:             Packet{size: 100},
			waitTime:      0,
			expectedError: false,
		},
		{
			name: "not enough tokens with max packet size",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 99,
				fillRate:      20,
				lastTime:      time.Now(),
			},
			p:             Packet{size: 100},
			waitTime:      0,
			expectedError: true,
		},
		{
			name: "tokens are just enough",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 0,
				fillRate:      20,
				lastTime:      time.Now(),
			},
			p:             Packet{size: 20},
			waitTime:      time.Millisecond * 1000,
			expectedError: false,
		},
		{
			name: "tokens are just not enough",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 0,
				fillRate:      20,
				lastTime:      time.Now(),
			},
			p:             Packet{size: 20},
			waitTime:      time.Millisecond * 999,
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.tb.lastTime = time.Now().Add(-test.waitTime)
			err := test.tb.sendPacket(test.p, time.Now())
			if err != nil && !test.expectedError {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && test.expectedError {
				t.Error("expected error but got none")
			}
		})
	}
}

func BenchmarkTokenBucket(b *testing.B) {
	tb := TokenBucket{maxCapacity: 100, currentTokens: 0, fillRate: 20, lastTime: time.Now()}
	p := Packet{size: 10}
	for i := 0; i < b.N; i++ {
		_ = tb.sendPacket(p, time.Now())
	}
}
