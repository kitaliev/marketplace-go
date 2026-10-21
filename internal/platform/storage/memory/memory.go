package memory

import (
	"context"
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

type Validatable[T any] interface {
	GetID() int64
}

type Data[T Validatable[T]] struct {
	mu    sync.Mutex
	items map[int64]T
}

func New[T Validatable[T]](items []T) *Data[T] {
	d := Data[T]{}
	d.items = make(map[int64]T, len(items))
	for _, p := range items {
		d.items[p.GetID()] = p
	}
	return &d
}

func (f *Data[T]) GetByID(ctx context.Context, id int64) (T, error) {
	var t T
	if err := ctx.Err(); err != nil {
		return t, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.items[id]
	if !ok {
		return t, ErrNotFound
	}
	return p, nil
}

func (f *Data[T]) Save(ctx context.Context, item T) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.items[item.GetID()] = item
	return nil
}
