package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// add timeout to avoid long waiting
func main() {
	rand.Seed(time.Now().Unix())
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	chanForResp := make(chan resp)
	go RPCCall(ctx, chanForResp)
	var response resp
	response = <-chanForResp

	fmt.Println(response.res, response.err)
	// cancel()
}

type resp struct {
	res int
	err error
}

func RPCCall(ctx context.Context, ch chan<- resp) {
	select {
	case <-ctx.Done():
		ch <- resp{
			res: 0,
			err: errors.New("timeout expired"),
		}
	case <-time.After(time.Second * time.Duration(rand.Intn(5))):
		// sleep 0-4 sec
		ch <- resp{
			res: rand.Int(),
		}

	}

}
