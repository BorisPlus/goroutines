package main

import (
	"fmt"
	"time"
)

const secondsCount, iterationsCount = 3, 5

func worker() {
	fmt.Println("... job start...")
	time.Sleep(time.Duration(secondsCount) * time.Second)
	fmt.Println("... job end...")
}

func main() {
	fmt.Println("... app start ...")
	for i := 0; i < iterationsCount; i++ {
		go worker()
	}
	fmt.Println("... app end ...")
}
