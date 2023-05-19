package utils

import (
	"context"
	"sync"
)

type CounterContext interface {
	CurrentContext() context.Context
	//Next wait for next number
	AwaitNext() <-chan int
	Increment()
}

func NewCounterContext() CounterContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &counterContext{
		count:  0,
		ctx:    ctx,
		cancel: cancel,
	}
}

type counterContext struct {
	mutex  sync.Mutex
	count  int
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *counterContext) CurrentContext() context.Context {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.ctx
}

var _ CounterContext = &counterContext{}

func (c *counterContext) AwaitNext() <-chan int {
	c.mutex.Lock()
	done := c.ctx.Done()
	currentCount := c.count
	c.mutex.Unlock()
	ch := make(chan int, 1)
	go func() {
		<-done
		ch <- currentCount + 1
	}()
	return ch
}

func (c *counterContext) Increment() {
	c.cancel()
	c.mutex.Lock()
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.count++
	c.mutex.Unlock()
}
