# Горутины

Развитие примера запуска горутин с учетом персональной обработкой каждой горутиной определенных входных параметров и полученный ею результатов в условии ограниченного числа работающих одновременно горутин

## Пример не отрабатывающих в фоне горутин

__Комментарий__: родительский процесс завершится, не дождавшись окончания фоновых

<details><summary>goroutine_000.go</summary>

```go
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
```

</details>

## Пример отрабатывающих в фоне горутин

__Комментарий__: родительский процесс ждет максимально минимально необходимое для отработки всех горутин время

<details><summary>goroutine_001.go</summary>


```go
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
```

</details>

## Выше + с эффектом разной длительности исполнения

__Комментарий__: введена случайного характера задержка исполнения задач

<details><summary>goroutine_002_rand.go</summary>

```go
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
```

</details>

## Пример отработки горутин ровно столько, сколько им надо по времени

__Комментарий__: канал следит за фактом исполнения горутин, в случе отработки всех - завершается, а не ждет максимально минимально необходимое для отработки всех горутин время

<details><summary>goroutine_003_chan.go</summary>

```go
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
```

</details>

## Выше + с ограниченным числом исполнителей

<details><summary>goroutine_004_pool.go</summary>

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

const secondsCount, jobsCount, workerCount = 3, 15, 3

func logic(input int) int {
    return input * 2
}

func worker(id int, jobs <-chan int, resultsChan chan<- [3]int) {
    for input := range jobs {
        fmt.Println("... worker ID:", id, "will start ...")
        secondsToSleep := rand.Intn(secondsCount)
        time.Sleep(time.Duration(secondsToSleep) * time.Second)
        output := logic(input)
        fmt.Println("... worker ID:", id, "worked", secondsToSleep, "sec. and end...")
        result := [3]int{id, input, output}
        resultsChan <- result
    }
}

func main() {

    fmt.Println("... app start ...")
    // result := [3]int{0, 0, 0}
    // const jobsCount = 5
    // int k := 0
    jobsChan := make(chan int, jobsCount)
    resultsChan := make(chan [3]int, jobsCount)

    for w := 1; w <= workerCount; w++ {
        go worker(w, jobsChan, resultsChan)
    }

    for j := 1; j <= jobsCount; j++ {
        jobsChan <- j
    }
    // close(jobsChan)

    for r := 1; r <= jobsCount; r++ {
        // <- resultsChan
        result := <-resultsChan
        // fmt.Println("... worker ", id, inputed, outputed, " job end ...")
        fmt.Println("... worker ", result[0], result[1], result[2], " job end ...")
        // fmt.Println("... worker ", result, " job end ...")
    }
    fmt.Println("... app end ...")
}
```

</details>

## Выше + выводом результата в структурированный map

<details><summary>goroutine_005_pool_map_result.go</summary>

```go
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
```

</details>

## Вывод

Родительский процесс отслеживает результаты (представлены map-структурой) работы ограниченного числа "дочерних" горутин.