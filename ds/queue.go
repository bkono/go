package ds

import (
	"sync"
)

type Event struct {
	Op string
}

type ConcurrentQueue[T any] struct {
	underlying *[]T
	lock       sync.RWMutex
}

func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
	return &ConcurrentQueue[T]{underlying: &[]T{}, lock: sync.RWMutex{}}
}

func (q *ConcurrentQueue[T]) Push(n T) {
	q.lock.Lock()
	defer q.lock.Unlock()
	Push(q.underlying, n)
}

func (q *ConcurrentQueue[T]) Pop() (n T) {
	q.lock.Lock()
	defer q.lock.Unlock()
	return Pop(q.underlying)
}

func (q *ConcurrentQueue[T]) Len() int {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return Len(q.underlying)
}

func (q *ConcurrentQueue[T]) Take(n int) []T {
	q.lock.Lock()
	defer q.lock.Unlock()
	return Take(q.underlying, n)
}

func Push[T any](q *[]T, n T) {
	*q = append(*q, n)
}

func Pop[T any](q *[]T) (n T) {
	if len(*q) == 0 {
		return
	}
	n = (*q)[0]
	*q = (*q)[1:]
	return
}

func Take[T any](q *[]T, n int) []T {
	if n > len(*q) {
		n = len(*q)
	}
	results := make([]T, n)
	for i := 0; i < n; i++ {
		results[i] = Pop(q)
	}
	return results
}

func Len[T any](q *[]T) int {
	return len(*q)
}
