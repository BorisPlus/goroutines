package main

import (
	"fmt"
	"math/rand"
	"time"
)

const secondsCount, iterationsCount = 3, 5

func worker(id int, c chan int) {
	// Объявляя канал как аргумент, функция превращается в "грязную"
	// и таким образом реализует возможность возврата значения во вне
	fmt.Println("... worker ID:", id, "will start ...")
	secondsToSleep := rand.Intn(secondsCount)
	time.Sleep(time.Duration(secondsToSleep) * time.Second)
	fmt.Println("... worker ID:", id, "worked", secondsToSleep, "sec. and end...")
	c <- id // Отправляет значение обратно к main
}

func main() {
	fmt.Println("... app start ...")
	c := make(chan int) // Делает канал для связи
	for i := 0; i < iterationsCount; i++ {
		go worker(i, c)
	}
	// если сделать i < iterationsCount+1, то будет
	// fatal error: all goroutines are asleep - deadlock!
	for i := 0; i < iterationsCount; i++ {
		workerID := <-c // Получает значение от канала // blocked waiting for a notification
		fmt.Println("worker ID:", workerID, "finished")
	}
	fmt.Println("... app end ...")
}
