package dispatch

import (
	"context"
	"sync/atomic"
	"time"
)

type Operation func(context.Context) error

type Runner struct {
	op     Operation
	delay  time.Duration
	active atomic.Int64
}

func NewRunner(delay time.Duration, op Operation) *Runner {
	return &Runner{op: op, delay: delay}
}

func (r *Runner) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	r.active.Add(1)
	go func() {
		defer close(done)
		defer r.active.Add(-1)
		workCtx := context.Background()
		for attempt := 0; attempt < 5; attempt++ {
			if err := r.op(workCtx); err == nil {
				return
			}
			time.Sleep(r.delay)
		}
	}()
	return done
}

func (r *Runner) Active() int64 { return r.active.Load() }
