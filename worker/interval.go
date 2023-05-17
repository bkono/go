package worker

import (
	"context"
	"sync"
	"time"

	"github.com/bkono/go/ds"
)

const (
	DefaultInterval    = 5 * time.Second
	DefaultTakeCount   = 5
	DefaultWorkerCount = 1
)

// IntervalWorker is a worker that runs at a given interval and processes a batch of items at a time.
// It allows concurrent pushing of new events, with a controlled, throttled processing.
// The worker will be idle when the queue is empty, starting processing only once the queue has items, and stopping when it reaches empty.
type IntervalWorker[T any] struct {
	Interval  time.Duration
	TakeCount int

	q              *ds.ConcurrentQueue[T]
	t              *time.Timer
	onEventFn      func(T)
	workCh         chan T
	workerCount    int
	done           chan struct{}
	workerCtx      context.Context
	workerCancelFn context.CancelFunc
	wg             *sync.WaitGroup
}

// IntervalWorkerOptions are options for creating a new IntervalWorker. For any value left empty, a sensible default will be used.
type IntervalWorkerOptions struct {
	// Default is 5 seconds.
	Interval time.Duration
	// Default is 5.
	TakeCount int
	// Default is 1.
	WorkerCount int
}

// NewIntervalWorker creates a new IntervalWorker with the given options and onEventFn.
func NewIntervalWorker[T any](opts *IntervalWorkerOptions, onEventFn func(T)) *IntervalWorker[T] {
	if opts.Interval == 0 {
		opts.Interval = DefaultInterval
	}
	if opts.TakeCount == 0 {
		opts.TakeCount = DefaultTakeCount
	}
	if opts.WorkerCount == 0 {
		opts.WorkerCount = DefaultWorkerCount
	}
	ch := make(chan T, opts.WorkerCount)
	ctx, cancelFn := context.WithCancel(context.Background())

	w := &IntervalWorker[T]{
		Interval:  opts.Interval,
		TakeCount: opts.TakeCount,

		q:              ds.NewConcurrentQueue[T](),
		onEventFn:      onEventFn,
		workCh:         ch,
		workerCount:    opts.WorkerCount,
		workerCtx:      ctx,
		workerCancelFn: cancelFn,
		wg:             &sync.WaitGroup{},
	}

	// spin up WorkerCount functions that will call onEvent with the work channel
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case t := <-ch:
					onEventFn(t)
				}
			}
		}()
	}

	return w
}

// Push adds a new item to the queue. If the queue is empty, it will start the interval timer for processing.
func (w *IntervalWorker[T]) Push(t T) {
	w.q.Push(t)
	if w.t == nil {
		w.t = time.AfterFunc(w.Interval, w.Run)
	}
}

// Run runs the worker, taking the next TakeCount items from the queue and distributing them to the w.WorkerCount for processing.
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

func (w *IntervalWorker[T]) Close() {
	w.workerCancelFn()
	close(w.workCh)
	w.wg.Wait()
}
