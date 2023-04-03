# Горутины

Развитие примера запуска ограниченного числа работающих одновременно горутин с учетом персональной обработки каждой из них определенных входных параметров.

## Пример не отрабатывающих в фоне горутин

__Комментарий__: родительский процесс завершится, не дождавшись окончания фоновых

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

[goroutine_000.go](./goroutine_000.go)

```bash
$ go run goroutine_000.go 
... app start ...
... app end ...
```

## Пример отрабатывающих в фоне горутин

__Комментарий__: родительский процесс ждет необходимое для отработки всех горутин время

```go
package main

import (
    "fmt"
    "time"
)

const secondsCount, iterationsCount = 3, 5

func worker() {
    fmt.Println("... job start ...")
    time.Sleep(time.Duration(secondsCount) * time.Second)
    fmt.Println("... job end ...")
}

func main() {
    fmt.Println("... app start ...")
    for i := 0; i < iterationsCount; i++ {
        go worker()
    }
    time.Sleep(time.Duration(iterationsCount*secondsCount)*time.Second + 1) // <-----
    fmt.Println("... app end ...")
}
```

[goroutine_001.go](./goroutine_001.go)

```bash
$ go run goroutine_001.go 
... app start ...
... job start ...
... job start ...
... job start ...
... job start ...
... job start ...
... job end ...
... job end ...
... job end ...
... job end ...
... job end ...
... app end ...
```

## Выше + с эффектом разной длительности исполнения

__Комментарий__: введена задержка исполнения задач случайного характера - исполняются разное время по длительности.

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

const secondsCount, iterationsCount = 3, 5

func worker() {
    fmt.Println("... job start ...")
    // Добавим стохастичности в длительность работы
    secondsToSleep := rand.Intn(secondsCount) // <-----
    time.Sleep(time.Duration(secondsToSleep) * time.Second) // <-----
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

[goroutine_002_rand.go](./goroutine_002_rand.go)

```bash
$ go run goroutine_002_rand.go 
... app start ...
... job start ...
... job start ...
... job start ...
... job worked 0 sec. and end ...
... job start ...
... job start ...
... job worked 1 sec. and end ...
... job worked 2 sec. and end ...
... job worked 2 sec. and end ...
... job worked 2 sec. and end ...
... app end ...
```

## Пример отработки горутин ровно столько, сколько им надо по времени

__Комментарий__: канал следит за фактом исполнения горутин, в случе отработки всех - сразу завершается. Также видно, что порядок завершения задач идет не по очереди их запуска.

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
    c <- id // <----- // Отправляет значение обратно к main
}

func main() {
    fmt.Println("... app start ...")
    c := make(chan int) // <----- // Делает канал для связи
    for i := 0; i < iterationsCount; i++ {
        go worker(i, c)
    }
    // если сделать i < iterationsCount+1, то будет
    // fatal error: all goroutines are asleep - deadlock!
    for i := 0; i < iterationsCount; i++ {
        workerID := <-c // <----- // Получает значение от канала (блокирующее ожидание)
        fmt.Println("worker ID:", workerID, "finished")
    }
    fmt.Println("... app end ...")
}
```

[goroutine_003_chan.go](./goroutine_003_chan.go)

```bash
$ go run goroutine_003_chan.go 
... app start ...
... worker ID: 4 will start ...
... worker ID: 1 will start ...
... worker ID: 1 worked 0 sec. and end ...
worker ID: 1 finished
... worker ID: 3 will start ...
... worker ID: 2 will start ...
... worker ID: 0 will start ...
... worker ID: 0 worked 1 sec. and end ...
worker ID: 0 finished
... worker ID: 3 worked 2 sec. and end ...
worker ID: 3 finished
... worker ID: 2 worked 2 sec. and end ...
worker ID: 2 finished
... worker ID: 4 worked 2 sec. and end ...
worker ID: 4 finished
... app end ...
```

## Выше + с ограниченным числом исполнителей на определенное число задач

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
    jobsChan := make(chan int, jobsCount)
    resultsChan := make(chan [3]int, jobsCount)

    for w := 1; w <= workerCount; w++ { // <----- 
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
```

[goroutine_004_pool.go](./goroutine_004_pool.go)

```bash
$ go run goroutine_004_pool.go 
... app start ...
... worker ID: 3 will start ...
... worker ID: 1 will start ...
... worker ID: 1 worked 0 sec. and end ...
... worker ID: 1 will start ...
... worker  1 2 4  job end ...
... worker ID: 2 will start ...
... worker ID: 3 worked 2 sec. and end ...
... worker ID: 3 will start ...
... worker ID: 2 worked 2 sec. and end ...
... worker  3 1 2  job end ...
... worker ID: 1 worked 2 sec. and end ...
... worker  2 4 8  job end ...
... worker  1 3 6  job end ...
... worker ID: 3 worked 1 sec. and end ...
... worker  3 5 10  job end ...
... app end ...
```

## Выше + с выводом результата в структурированный map

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
        result := map[string]int{ // <----- 
            "id":     id,
            "input":  input,
            "output": output,
        }
        resultsChan <- result
    }
}

func main() {

    fmt.Println("... app start ...")
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
        result := <-resultsChan
        fmt.Println("... worker ID:", result["id"],
            "start with INPUT:", result["input"],
            "and end with OUTPUT:", result["output"], " ...")
    }
    fmt.Println("... app end ...")
}
```

[goroutine_005_pool_map_result.go](./goroutine_005_pool_map_result.go)

```bash
$ go run goroutine_005_pool_map_result.go 
... app start ...
... worker ID: 3 will start ...
... worker ID: 1 will start ...
... worker ID: 1 worked 0 sec. and end ...
... worker ID: 2 will start ...
... worker ID: 1 start with INPUT: 2 and end with OUTPUT: 4  ...
... worker ID: 1 will start ...
... worker ID: 1 worked 2 sec. and end ...
... worker ID: 1 will start ...
... worker ID: 3 worked 2 sec. and end ...
... worker ID: 2 worked 2 sec. and end ...
... worker ID: 1 start with INPUT: 4 and end with OUTPUT: 8  ...
... worker ID: 2 start with INPUT: 3 and end with OUTPUT: 6  ...
... worker ID: 3 start with INPUT: 1 and end with OUTPUT: 2  ...
... worker ID: 1 worked 1 sec. and end ...
... worker ID: 1 start with INPUT: 5 and end with OUTPUT: 10  ...
... app end ...
```

## Вывод

В итоге данный подход можно масштабировать на относительно "одноразовые" задачи, когда необходимо:

* ограничить пул одновременно работающих воркеров-горутин;
* привести структуру результата к формату ключ-значение;
* отразить в результате персонализированно для какого воркера на какой запрос, что он ответил.
