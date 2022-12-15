package main

import (
	"fmt"
	"github.com/Resized/traffic-policing-algorithms/pkg/leaky_bucket"
	"github.com/Resized/traffic-policing-algorithms/pkg/token_bucket"
	"log"
	"time"
)

func main() {
	fmt.Println("Please choose which algorithm to run:")
	fmt.Println("1) Leaky Bucket")
	fmt.Println("2) Token Bucket")
	var ans int
	_, err := fmt.Scanln(&ans)
	if err != nil {
		log.Fatal(err)
	}

	switch ans {
	case 1:
		var p = [10]int{10, 15, 20, 25, 30, 35, 40, 45, 50, 105}
		var lb = leaky_bucket.LeakyBucket{}
		lb.Init(100, 20)
		for {
			for i := 0; i < 10; i++ {
				err := lb.Add(p[i], time.Now())
				if err != nil {
					fmt.Printf("error sending payload: %v\n", err)
				} else {
					fmt.Printf("accepting payload of size %v, current tokens: %.3v\n", p[i], lb.GetCurrent())
				}
				time.Sleep(time.Second / 2)
			}
		}
	case 2:
		var p = [10]int{10, 15, 20, 25, 30, 35, 40, 45, 50, 105}
		var tb = token_bucket.TokenBucket{}
		tb.Init(100, 20)
		for {
			for i := 0; i < 10; i++ {
				err := tb.Take(p[i], time.Now())
				if err != nil {
					fmt.Printf("error sending payload: %v\n", err)
				} else {
					fmt.Printf("accepting payload of size %v, current tokens: %.3v\n", p[i], tb.GetCurrentTokens())
				}
				time.Sleep(time.Second / 2)
			}
		}
	default:
	}
}
