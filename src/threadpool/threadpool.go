package threadpool

import "sync"

type Pool struct {
	wg sync.WaitGroup
}

func New() *Pool {
	return &Pool{sync.WaitGroup{}}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) AddTask(task func()) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		task()
	}()
}

// add chan

func MakeChan[T any](task func() T) chan T {
	done := make(chan T)

	go func() {
		done <- task()
	}()
	return done
}
