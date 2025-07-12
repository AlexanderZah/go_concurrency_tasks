package main

import (
	"fmt"
	"sync"
)

// merge two channels
func main() {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 20)

	ch1 <- 1
	ch2 <- 2
	ch2 <- 4
	close(ch1)
	close(ch2)

	// ch3 := syncMerge(ch1, ch2) // bad sync way of merging
	ch3 := merge[int](ch1, ch2)

	for val := range ch3 {
		fmt.Println(val)
	}
}

func merge[T any](chans ...chan T) chan T {
	returnCh := make(chan T)
	go func() {
		var wg sync.WaitGroup
		for _, ch := range chans {
			wg.Add(1)
			go func(c chan T) {
				defer wg.Done()
				for val := range c {
					returnCh <- val
				}
			}(ch)
		}

		go func() {
			wg.Wait()
			close(returnCh)
		}()
	}()

	return returnCh
}

//TODO можно сделать mergeCh через воркер пул
// go func() {
// 		for _, url := range urls {
// 			ch <- url
// 		}
// 		close(ch)
// 	}()
// 	numGorutins := 30
// 	wg.Add(numGorutins)
// 	for range numGorutins {
// 		go func() {
// 			defer wg.Done()

// 			for u := range ch {
// 				code := sendRequest(u)
// 				locker.Lock()
// 				codes[code]++
// 				locker.Unlock()
// 			}
// 		}()
// 	}
