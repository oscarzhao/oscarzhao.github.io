package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	tasks := []func(){
		func() { time.Sleep(time.Second); fmt.Println("1 sec later") },
		func() { time.Sleep(time.Second * 2); fmt.Println("2 sec later") },
	}
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	for _, task := range tasks {
		task := task
		go func() {
			defer wg.Done()
			task()
		}()
	}
	wg.Wait()
	fmt.Println("finish")
}
