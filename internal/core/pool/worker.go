package pool

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
)

type Task struct {
	ID       string
	Func     func() error
	Priority int
}

type WorkerPool struct {
	workers       int
	taskQueue     chan Task
	wg            sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
	activeWorkers int64
}

func NewWorkerPool(workers int) *WorkerPool {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:   workers,
		taskQueue: make(chan Task, workers*2),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.taskQueue)
	wp.wg.Wait()
}

func (wp *WorkerPool) Submit(task Task) error {
	select {
	case wp.taskQueue <- task:
		return nil
	case <-wp.ctx.Done():
		return context.Canceled
	}
}

func (wp *WorkerPool) SubmitWithContext(ctx context.Context, task Task) error {
	select {
	case wp.taskQueue <- task:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-wp.ctx.Done():
		return context.Canceled
	}
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case task, ok := <-wp.taskQueue:
			if !ok {
				return
			}

			atomic.AddInt64(&wp.activeWorkers, 1)
			if err := task.Func(); err != nil {
				// Log error or handle it appropriately
			}
			atomic.AddInt64(&wp.activeWorkers, -1)
		case <-wp.ctx.Done():
			return
		}
	}
}

func (wp *WorkerPool) AvailableWorkers() int64 {
	return int64(wp.workers) - atomic.LoadInt64(&wp.activeWorkers)
}

func (wp *WorkerPool) QueueSize() int {
	return len(wp.taskQueue)
}
