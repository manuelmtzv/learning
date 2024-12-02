package workers

import (
	"container/heap"
	"context"
	"order-processing/internal/store"
	"sync"

	"github.com/charmbracelet/log"
)

type Balancer interface {
	Balance(context.Context, <-chan *Request)
}

type BalancerWorker struct {
	logger  *log.Logger
	store   *store.Storage
	pool    Pool
	done    chan *ProcessorWorker
	poolMux sync.Mutex
}

func NewBalancer(logger *log.Logger, store *store.Storage, workers int) Balancer {
	b := &BalancerWorker{
		logger: logger,
		store:  store,
		pool:   make(Pool, 0, workers),
		done:   make(chan *ProcessorWorker),
	}

	for i := 0; i < workers; i++ {
		w := NewProcessor(logger, store)
		go w.Work(context.Background())
		heap.Push(&b.pool, w)
	}

	return b
}

func (b *BalancerWorker) Balance(ctx context.Context, workStream <-chan *Request) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case req := <-workStream:
				b.dispatch(ctx, req)
			case w := <-b.done:
				b.completed(w)
			}
		}
	}()
}

func (b *BalancerWorker) dispatch(ctx context.Context, req *Request) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				b.poolMux.Lock()
				w := heap.Pop(&b.pool).(*ProcessorWorker)
				w.requests <- req
				w.pending++
				heap.Push(&b.pool, w)
				b.poolMux.Unlock()
			}
		}
	}()
}

func (b *BalancerWorker) completed(w *ProcessorWorker) {
	b.poolMux.Lock()
	defer b.poolMux.Unlock()
	w.pending--
	heap.Remove(&b.pool, w.index)
	heap.Push(&b.pool, w)
}

type Pool []*ProcessorWorker

func (p Pool) Len() int {
	return len(p)
}

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p Pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}

func (p *Pool) Push(x interface{}) {
	*p = append(*p, x.(*ProcessorWorker))
}

func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[0 : n-1]
	return x
}
