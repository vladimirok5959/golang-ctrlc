package main

import (
	"context"
	"fmt"
	"time"
)

type TestObject struct {
	ctx    context.Context
	cancel context.CancelFunc
	chDone chan bool
}

func Run() *TestObject {
	ctx, cancel := context.WithCancel(context.Background())
	w := &TestObject{ctx: ctx, cancel: cancel, chDone: make(chan bool)}

	go func() {
		for {
			select {
			case <-w.ctx.Done():
				w.chDone <- true
				return
			case <-time.After(1 * time.Second):
				fmt.Printf("[TestObject]: I do something every second!\n")
				fmt.Printf("[TestObject]: Press CTRL + C for shutdown...\n")
			}
		}
	}()

	return w
}

func (this *TestObject) Shutdown(ctx context.Context) error {
	fmt.Printf("[TestObject]: OK! I will shutdown!\n")

	this.cancel()

	select {
	case <-this.chDone:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
