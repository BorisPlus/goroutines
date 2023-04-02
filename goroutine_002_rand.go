package main

import (
	"fmt"
	"math/rand"
	"time"
)

const secondsCount, iterationsCount = 3, 5

func worker() {
	// Добавим стохастичности в длительность работы
	fmt.Println("... job start ...")
	secondsToSleep := rand.Intn(secondsCount)
	time.Sleep(time.Duration(secondsToSleep) * time.Second)
	fmt.Println("... job worked", secondsToSleep, "sec. and end...")
}

func main() {
	fmt.Println("... app start ...")
	for i := 0; i < iterationsCount; i++ {
		go worker()
	}
	time.Sleep(time.Duration(iterationsCount*secondsCount)*time.Second + 1)
	fmt.Println("... app end ...")
}
