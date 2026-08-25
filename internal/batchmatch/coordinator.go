package batchmatch

import (
	"context"
	"sync/atomic"
)

type Resolver interface {
	Resolve(string) (string, error)
}

type ResolverFunc func(string) (string, error)

func (f ResolverFunc) Resolve(id string) (string, error) { return f(id) }

type Coordinator struct {
	resolver Resolver
	active   atomic.Int64
}

func NewCoordinator(resolver Resolver) *Coordinator {
	return &Coordinator{resolver: resolver}
}

func (c *Coordinator) Run(ctx context.Context, ids []string) ([]string, error) {
	jobs := make(chan string)
	results := make(chan string)
	errs := make(chan error)
	c.active.Add(1)
	go func() {
		defer c.active.Add(-1)
		for _, id := range ids {
			value, err := c.resolver.Resolve(id)
			if err != nil {
				errs <- err
				return
			}
			jobs <- value
		}
		close(jobs)
	}()
	c.active.Add(1)
	go func() {
		defer c.active.Add(-1)
		defer close(results)
		for value := range jobs {
			results <- value
		}
	}()
	var collected []string
	for {
		select {
		case value, ok := <-results:
			if !ok {
				return collected, nil
			}
			collected = append(collected, value)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (c *Coordinator) Active() int64 { return c.active.Load() }
