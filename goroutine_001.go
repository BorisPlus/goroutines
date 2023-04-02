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
	// Try to change above to
	// time.Sleep(time.Second)
	// or cut/comment this row
	// and you will see not all gorutines were run
	fmt.Println("... app end ...")
}
