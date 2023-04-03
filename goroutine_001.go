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
	time.Sleep(time.Duration(iterationsCount*secondsCount)*time.Second + 1)
	fmt.Println("... app end ...")
}
