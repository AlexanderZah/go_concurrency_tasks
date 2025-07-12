package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// junior решение
func urlsMain() map[int]int {
	rand.Seed(time.Now().UnixNano())
	urls := make([]string, 1000)
	codes := make(map[int]int)
	wg := &sync.WaitGroup{}
	locker := &sync.Mutex{}

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			code := sendRequest(url)
			locker.Lock()
			codes[code]++
			locker.Unlock()
		}(url)
	}
	wg.Wait()

	// wg.Wait()
	return codes
}

// middle решение
func urlsMainWithChannels() map[int]int {
	rand.Seed(time.Now().UnixNano())
	urls := make([]string, 1000)
	codes := make(map[int]int)
	wg := &sync.WaitGroup{}
	locker := &sync.Mutex{}
	ch := make(chan string)

	go func() {
		for _, url := range urls {
			ch <- url
		}
		close(ch)
	}()
	numGorutins := 30
	wg.Add(numGorutins)
	for range numGorutins {
		go func() {
			defer wg.Done()

			for u := range ch {
				code := sendRequest(u)
				locker.Lock()
				codes[code]++
				locker.Unlock()
			}
		}()
	}
	fmt.Println(runtime.NumGoroutine())
	wg.Wait()

	return codes
}

func urlsMainSemaphora() map[int]int {
	rand.Seed(time.Now().UnixNano())
	urls := make([]string, 1000)
	codes := make(map[int]int)
	locker := &sync.Mutex{}

	sema := make(chan struct{}, 5)

	for _, url := range urls {
		sema <- struct{}{}
		go func(url string) {
			defer func() {
				<-sema
			}()

			code := sendRequest(url)

			locker.Lock()
			codes[code]++
			locker.Unlock()
		}(url)
	}
	return codes
}

func sendRequest(url string) (code int) {
	time.Sleep(time.Microsecond * 20)

	if rand.Intn(2) == 0 {
		return 200
	}
	return 500
}

func main() {
	// codes := urlsMain()
	// codes := urlsMainWithChannels()
	codes := urlsMainSemaphora()
	fmt.Println(codes)
	for code, num := range codes {
		fmt.Println(code, num)
	}
}
