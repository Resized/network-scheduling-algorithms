package main

import (
	"testing"
	"time"
)

func TestTokenBucketTableDriven(t *testing.T) {
	var tests = []struct {
		name          string
		tb            TokenBucket
		p             Packet
		expectedError bool
	}{
		{
			name: "packet too large",
			tb: TokenBucket{
				maxCapacity:   100,
				currentTokens: 100,
				fillRate:      20,
				lastTime:      time.Now(),
			},
			p:             Packet{size: 101},
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
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.tb.sendPacket(test.p)
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
		_ = tb.sendPacket(p)
	}
}
