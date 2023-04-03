package main

import (
	"fmt"
	"math/rand"
	"time"
)

const secondsCount, jobsCount, workerCount = 3, 5, 3

func logic(input int) int {
	return input * 2
}

func worker(id int, jobs <-chan int, resultsChan chan<- [3]int) {
	for input := range jobs {
		fmt.Println("... worker ID:", id, "will start ...")
		secondsToSleep := rand.Intn(secondsCount)
		time.Sleep(time.Duration(secondsToSleep) * time.Second)
		output := logic(input)
		fmt.Println("... worker ID:", id, "worked", secondsToSleep, "sec. and end ...")
		result := [3]int{id, input, output}
		resultsChan <- result
	}
}

func main() {

	fmt.Println("... app start ...")
	jobsChan := make(chan int, jobsCount)
	resultsChan := make(chan [3]int, jobsCount)

	for w := 1; w <= workerCount; w++ {
		go worker(w, jobsChan, resultsChan)
	}

	for j := 1; j <= jobsCount; j++ {
		jobsChan <- j
	}

	for r := 1; r <= jobsCount; r++ {
		result := <-resultsChan
		fmt.Println("... worker ", result[0], result[1], result[2], " job end ...")
	}
	fmt.Println("... app end ...")
}
