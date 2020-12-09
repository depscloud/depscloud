package eventlp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
)

// Task represents a unit of work
type Task func(ctx context.Context)

// New constructs an event loop.
func New() *EventLoop {
	return &EventLoop{
		clock:    clockwork.NewRealClock(),
		wait:     200 * time.Millisecond,
		mu:       &sync.Mutex{},
		queue:    &LinkedList{},
		shutdown: false,
	}
}

// EventLoop is a simple event loop that supports concurrent processing. Tasks
// are processed in the order they're received on the queue. Some tasks may
// require less time to process.
type EventLoop struct {
	clock    clockwork.Clock
	wait     time.Duration
	mu       *sync.Mutex
	queue    *LinkedList
	shutdown bool
}

// Clock overrides the clock used for the event loop.
func (p *EventLoop) WithClock(clock clockwork.Clock) *EventLoop {
	p.clock = clock
	return p
}

// Submit adds an item to the end of the queue
func (p *EventLoop) Submit(task Task) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.shutdown {
		return fmt.Errorf("event loop shutdown, no longer accepting work")
	}

	p.queue.PushBack(task)
	return nil
}

// Once is used to execute a single run of the event loop.
func (p *EventLoop) Once(ctx context.Context) {
	p.mu.Lock()
	task, ok := p.queue.PopFront().(Task)
	p.mu.Unlock()

	if ok && task != nil {
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
		queued := p.queue.Size()
		p.mu.Unlock()

		if shutdown && queued == 0 {
			return nil
		}

		p.clock.Sleep(p.wait)
	}
}

// GracefullyStop waits for the event loop to drain before returning.
func (p *EventLoop) GracefullyStop() error {
	_ = p.Stop()

	for p.queue.Size() > 0 {
		p.clock.Sleep(p.wait)
	}

	return nil
}

// Stop shuts down the server and doesn't wait before returning.
func (p *EventLoop) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.shutdown = true
	return nil
}
