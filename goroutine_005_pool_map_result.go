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

func worker(id int, jobs <-chan int, resultsChan chan<- map[string]int) {
	for input := range jobs {
		fmt.Println("... worker ID:", id, "will start ...")
		secondsToSleep := rand.Intn(secondsCount)
		time.Sleep(time.Duration(secondsToSleep) * time.Second)
		output := logic(input)
		fmt.Println("... worker ID:", id, "worked", secondsToSleep, "sec. and end...")
		result := map[string]int{
			"id":     id,
			"input":  input,
			"output": output,
		}
		resultsChan <- result
	}
}

func main() {

	fmt.Println("... app start ...")
	// result := [3]int{0, 0, 0}
	// int k := 0
	jobsChan := make(chan int, jobsCount)
	resultsChan := make(chan map[string]int, jobsCount)

	for w := 1; w <= workerCount; w++ {
		go worker(w, jobsChan, resultsChan)
	}

	for j := 1; j <= jobsCount; j++ {
		jobsChan <- j
	}
	close(jobsChan)

	for r := 1; r <= jobsCount; r++ {
		// <- resultsChan
		result := <-resultsChan
		// fmt.Println("... worker ", id, inputed, outputed, " job end ...")
		fmt.Println("... worker ID:", result["id"],
			"start with INPUT:", result["input"],
			"and end with OUTPUT:", result["output"], " ...")
		// fmt.Println("... worker ", result, " job end ...")
	}
	fmt.Println("... app end ...")
}
