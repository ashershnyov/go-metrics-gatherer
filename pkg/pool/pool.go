package pool

import "sync"

// Resetable describes a type that has an ability to be reset.
type Resetable interface {
	Reset()
}

// Pool is a generic-using wrapper around sync.Pool.
type Pool[T Resetable] struct {
	pool sync.Pool
}

// New creates a new Pool with the given `new` function.
func New[T Resetable](new func() T) *Pool[T] {
	p := &Pool[T]{}
	p.pool = sync.Pool{
		New: func() interface{} {
			return new()
		},
	}
	return p
}

// Get retrieves an object from the pool.
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put returns an object to the pool for reuse.
func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}
