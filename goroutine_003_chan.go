package main

import (
	"fmt"
	"math/rand"
	"time"
)

const secondsCount, iterationsCount = 3, 5

func worker(id int, c chan int) {
	fmt.Println("... worker ID:", id, "will start ...")
	secondsToSleep := rand.Intn(secondsCount)
	time.Sleep(time.Duration(secondsToSleep) * time.Second)
	fmt.Println("... worker ID:", id, "worked", secondsToSleep, "sec. and end ...")
	c <- id
}

func main() {
	fmt.Println("... app start ...")
	c := make(chan int) 
	for i := 0; i < iterationsCount; i++ {
		go worker(i, c)
	}
	for i := 0; i < iterationsCount; i++ {
		workerID := <-c
		fmt.Println("worker ID:", workerID, "finished")
	}
	fmt.Println("... app end ...")
}
