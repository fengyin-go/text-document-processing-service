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

// Run 解析一批规则 id，返回解析结果列表。
//
// 行为约定：
//   - 任一规则解析失败（例如规则不存在）时立即返回该错误，不再继续后续规则；
//   - 上下文取消时立即返回 ctx.Err()；
//   - 无论以何种方式返回，都通过派生的可取消上下文通知已启动的生产与收集
//     goroutine 退出，确保不残留后台工作（Active() 归零）。
//
// 规则解析仍按原有方式逐个调用 Resolver.Resolve，调用语义保持不变。
func (c *Coordinator) Run(ctx context.Context, ids []string) ([]string, error) {
	// 派生可取消上下文：Run 返回时（正常/出错/取消）统一触发 cancel，
	// 令生产者与收集者解除阻塞并退出。
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan string)
	results := make(chan string)
	// 缓冲为 1：生产者至多上报一个错误，发送永远不会阻塞，
	// 即使主循环已因 ctx 取消而返回也不会泄漏生产者。
	errc := make(chan error, 1)

	// 生产者：逐个解析规则 id，成功则投递到 jobs；失败则上报 errc 并退出。
	c.active.Add(1)
	go func() {
		defer c.active.Add(-1)
		defer close(jobs)
		for _, id := range ids {
			value, err := c.resolver.Resolve(id)
			if err != nil {
				errc <- err
				return
			}
			select {
			case jobs <- value:
			case <-ctx.Done():
				return
			}
		}
	}()

	// 收集者：把 jobs 中的结果转发到 results 供主循环消费。
	c.active.Add(1)
	go func() {
		defer c.active.Add(-1)
		defer close(results)
		for value := range jobs {
			select {
			case results <- value:
			case <-ctx.Done():
				return
			}
		}
	}()

	var collected []string
	for {
		select {
		case value, ok := <-results:
			if !ok {
				// results 已关闭：收集者已退出。优先返回解析错误，
				// 其次返回上下文取消错误，最后返回正常结果。
				select {
				case err := <-errc:
					return nil, err
				default:
				}
				if err := ctx.Err(); err != nil {
					return nil, err
				}
				return collected, nil
			}
			collected = append(collected, value)
		case err := <-errc:
			// 解析失败：立即取消派生上下文，令生产/收集 goroutine 退出后返回错误。
			cancel()
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (c *Coordinator) Active() int64 { return c.active.Load() }
