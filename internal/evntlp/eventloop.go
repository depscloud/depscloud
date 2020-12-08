package evntlp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Workiva/go-datastructures/queue"
)

// Task represents a unit of work
type Task func(ctx context.Context)

// New constructs an event loop.
func New() *EventLoop {
	return &EventLoop{
		wait:     200 * time.Millisecond,
		mu:       &sync.Mutex{},
		queue:    queue.New(0),
		shutdown: false,
	}
}

// EventLoop is a simple event loop that supports concurrent processing. Tasks
// are processed in the order they're received on the queue. Some tasks may
// require less time to process.
type EventLoop struct {
	wait     time.Duration
	mu       *sync.Mutex
	queue    *queue.Queue
	shutdown bool
}

// Submit adds an item to the end of the queue
func (p *EventLoop) Submit(task Task) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.shutdown {
		return fmt.Errorf("evntlp shutdown, no longer accepting work")
	}

	return p.queue.Put(task)
}

// Once is used to execute a single run of the event loop.
func (p *EventLoop) Once(ctx context.Context) {
	p.mu.Lock()
	data, err := p.queue.Poll(1, p.wait)
	p.mu.Unlock()

	if err == nil && len(data) > 0 {
		task := data[0].(Task)
		task(ctx)
	}
}

// Start runs the event loop.
func (p *EventLoop) Start(parent context.Context) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	for {
		p.Once(ctx)

		p.mu.Lock()
		shutdown := p.shutdown
		queued := p.queue.Len()
		p.mu.Unlock()

		if shutdown && queued == 0 {
			return nil
		}
	}
}

// GracefullyStop waits for the event loop to drain before returning.
func (p *EventLoop) GracefullyStop() error {
	_ = p.Stop()

	for p.queue.Len() > 0 {
		time.Sleep(p.wait)
	}

	return nil
}

// Stop shuts down the server and doesn't wait before returning.
func (p *EventLoop) Stop() error {
	p.mu.Lock()
	p.shutdown = true
	p.mu.Unlock()
	return nil
}
