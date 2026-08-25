package requestmeta

import "sync"

type Meta struct {
	Identity string
	Labels   []string
}

type Pool struct {
	pool sync.Pool
}

func NewPool() *Pool {
	return &Pool{pool: sync.Pool{New: func() any { return &Meta{} }}}
}

func (p *Pool) Acquire(identity string) *Meta {
	meta := p.pool.Get().(*Meta)
	meta.Identity = identity
	return meta
}

func (p *Pool) Release(meta *Meta) {
	p.pool.Put(meta)
}
