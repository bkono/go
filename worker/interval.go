package worker

import (
	"time"

	"github.com/bkono/go/ds"
)

type IntervalWorker[T any] struct {
	Interval  time.Duration
	TakeCount int

	q         *ds.ConcurrentQueue[T]
	t         *time.Timer
	onEventFn func(T)
}

func NewIntervalWorker[T any](interval time.Duration, takeCount int, onEventFn func(T)) *IntervalWorker[T] {
	return &IntervalWorker[T]{
		Interval:  interval,
		TakeCount: takeCount,

		q:         ds.NewConcurrentQueue[T](),
		onEventFn: onEventFn,
	}
}

func (w *IntervalWorker[T]) Push(t T) {
	w.q.Push(t)
	if w.t == nil {
		w.t = time.AfterFunc(w.Interval, w.Run)
	}
}

func (w *IntervalWorker[T]) Run() {
	if w.q.IsEmpty() {
		return
	}

	for _, evt := range w.q.Take(w.TakeCount) {
		w.onEventFn(evt)
	}

	// schedule another run in w.Interval
	if w.q.IsEmpty() {
		w.t = nil
	} else {
		w.t.Reset(w.Interval)
	}
}
