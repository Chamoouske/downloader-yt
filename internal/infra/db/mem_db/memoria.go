package memoria

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

type MemoriaDatabase[T any] struct {
	data map[string]T
	mu   sync.RWMutex
}

func NewMemoriaDatabase[T any]() *MemoriaDatabase[T] {
	return &MemoriaDatabase[T]{data: map[string]T{}}
}

func (r *MemoriaDatabase[T]) Save(id string, v T) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[id] = v
	return nil
}

func (r *MemoriaDatabase[T]) Get(id string) (T, error) {
	v, ok := r.data[id]
	var zero T
	if !ok {
		return zero, ErrNotFound
	}

	return v, nil
}

func (r *MemoriaDatabase[T]) Remove(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return ErrNotFound
	}
	delete(r.data, id)
	return nil
}
