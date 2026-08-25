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
		for attempt := 0; attempt < 5; attempt++ {
			// 取消信号已到，停止重试，不再发起任何新的后端调用。
			if ctx.Err() != nil {
				return
			}
			// 将可取消的 ctx 透传给操作，使其内部调用也能及时收住。
			if err := r.op(ctx); err == nil {
				return
			}
			// 重试间隔必须可被取消打断，否则关闭阶段会死等 delay 跑完整条重试链。
			select {
			case <-ctx.Done():
				return
			case <-time.After(r.delay):
			}
		}
	}()
	return done
}

func (r *Runner) Active() int64 { return r.active.Load() }
