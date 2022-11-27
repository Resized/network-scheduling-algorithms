package main

import (
	"fmt"
	"time"
	"token-bucket/pkg/token_bucket"
)

func main() {
	var p = [10]int{10, 15, 20, 25, 30, 35, 40, 45, 50, 105}
	var tb = token_bucket.TokenBucket{}
	tb.InitTokenBucket(100, 20)
	for {
		for i := 0; i < 10; i++ {
			err := tb.Take(p[i], time.Now(), true)
			if err != nil {
				fmt.Printf("error sending packet: %v\n", err)
			}
			time.Sleep(time.Second / 2)
		}
	}
}
