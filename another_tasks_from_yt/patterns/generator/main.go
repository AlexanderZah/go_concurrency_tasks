package main

import (
	"fmt"
	"sync"
	"time"
)

func makeGenerator(done <-chan struct{}, wg *sync.WaitGroup) <-chan int {
	ch := make(chan int, 1)
	var i = 0

	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				close(ch)
				fmt.Println("done")
				return
			default:
				time.Sleep(time.Microsecond * 250)
				ch <- i
				i++
			}
		}
	}()
	return ch
}

func main() {
	done := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(2)

	ch := makeGenerator(done, &wg)

	go func() {
		defer wg.Done()
		for v := range ch {
			fmt.Println("value:", v)
		}
	}()

	time.Sleep(time.Second)
	close(done)
	wg.Wait()
}
