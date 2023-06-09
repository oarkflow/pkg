package pool

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

var _ Pool = new(limitedPool)

// limitedPool contains all information for a limited pool instance.
type limitedPool struct {
	workers uint
	work    chan *workUnit
	cancel  chan struct{}
	closed  bool
	m       sync.RWMutex
	// the max duration goroutine alive
	currworkers uint // current workers
	minWorkers  uint
	qsize       uint // queue size
	timeToLive  time.Duration
}

// NewLimited returns a new limited pool instance
func NewLimited(workers uint) Pool {
	if workers == 0 {
		panic("invalid workers '0'")
	}

	timeToLive := 1 * time.Minute
	return NewExtLimited((workers*10)/8, workers, 2*workers, timeToLive)
}

// NewExtLimited returns a new limited pool instance with args
func NewExtLimited(minworkers, maxworkers uint, qsize uint, ttl time.Duration) Pool {
	if maxworkers == 0 || minworkers == 0 {
		panic("invalid workers '0'")
	}

	p := &limitedPool{
		workers:    maxworkers,
		minWorkers: minworkers,
		timeToLive: ttl,
		qsize:      qsize,
	}
	if p.minWorkers <= 0 {
		p.minWorkers = maxworkers * 5 / 10
		if p.minWorkers <= 0 {
			p.minWorkers = 1
		}
	}

	if p.qsize < maxworkers*2 {
		p.qsize = maxworkers * 2
	}

	p.initialize()

	return p
}

// current workers
func (p *limitedPool) CurrWorkers() uint {
	return p.currworkers
}

func (p *limitedPool) MaxWorkers() uint {
	return p.workers
}

// Incomplete Task
func (p *limitedPool) IncompleteTasks() uint {
	return uint(len(p.work))
}

func (p *limitedPool) initialize() {
	p.work = make(chan *workUnit, p.qsize)
	p.cancel = make(chan struct{})
	p.closed = false

	p.currworkers = 0
	// fire up min workers here
	for i := 0; i < int(p.minWorkers); i++ {
		p.newWorker(p.work, p.cancel)
		p.currworkers++
	}
}

// passing work and cancel channels to newWorker() to avoid any potential race condition
// betweeen p.work read & write
func (p *limitedPool) newWorker(work chan *workUnit, cancel chan struct{}) {
	go func(p *limitedPool) {
		var wu *workUnit

		defer func(p *limitedPool) {
			if err := recover(); err != nil {

				trace := make([]byte, 1<<8)
				n := runtime.Stack(trace, false)
				s := fmt.Sprintf(errRecovery, err, string(trace[:int(math.Min(float64(n), float64(7000)))]))

				iwu := wu
				iwu.err = &ErrRecovery{s: s}
				close(iwu.done)

				// need to fire up new worker to replace this one as this one is exiting
				p.newWorker(p.work, p.cancel)
			}
		}(p)

		var value interface{}
		var err error

		for {
			if p.timeToLive <= 0 {
				select {
				case wu = <-work:

					// possible for one more nilled out value to make it
					// through when channel closed, don't quite understad the why
					if wu == nil {
						continue
					}

					// support for individual WorkUnit cancellation
					// and batch job cancellation
					if wu.cancelled.Load() == nil {
						value, err = wu.fn(wu)

						wu.writing.Store(struct{}{})

						// need to check again in case the WorkFunc cancelled this unit of work
						// otherwise we'll have a race condition
						if wu.cancelled.Load() == nil && wu.cancelling.Load() == nil {
							wu.value, wu.err = value, err

							// who knows where the Done channel is being listened to on the other end
							// don't want this to block just because caller is waiting on another unit
							// of work to be done first so we use close
							close(wu.done)
						}
					}
				case <-cancel:
					return
				}
			} else {
				select {
				case wu = <-work:

					// possible for one more nilled out value to make it
					// through when channel closed, don't quite understad the why
					if wu == nil {
						continue
					}

					// support for individual WorkUnit cancellation
					// and batch job cancellation
					if wu.cancelled.Load() == nil {
						value, err = wu.fn(wu)

						wu.writing.Store(struct{}{})

						// need to check again in case the WorkFunc cancelled this unit of work
						// otherwise we'll have a race condition
						if wu.cancelled.Load() == nil && wu.cancelling.Load() == nil {
							wu.value, wu.err = value, err

							// who knows where the Done channel is being listened to on the other end
							// don't want this to block just because caller is waiting on another unit
							// of work to be done first so we use close
							close(wu.done)
						}
					}
				case <-time.After(p.timeToLive):
					// too long idle exit
					p.m.Lock()
					if p.currworkers >= p.minWorkers {
						p.currworkers--
						p.m.Unlock()
						return
					}
					p.m.Unlock()

				case <-cancel:
					return
				}
			}
		}
	}(p)
}

// Queue queues the work to be run, and starts processing immediately
func (p *limitedPool) Queue(fn WorkFunc) WorkUnit {
	w := &workUnit{
		done: make(chan struct{}),
		fn:   fn,
	}

	go func() {
		p.m.RLock()
		if p.closed {
			w.err = &ErrPoolClosed{s: errClosed}
			if w.cancelled.Load() == nil {
				close(w.done)
			}
			p.m.RUnlock()
			return
		}

		needMoreWorks := false
		//
		workLen := len(p.work)
		if workLen > 0 {
			// need more work
			if p.currworkers < p.workers {
				needMoreWorks = true
			}
		}

		p.work <- w

		p.m.RUnlock()

		p.m.Lock()
		if needMoreWorks {
			//
			workLen = len(p.work)
			if workLen > 0 {
				// create work
				if p.currworkers < p.workers {
					p.newWorker(p.work, p.cancel)
					p.currworkers++
				}
			}
		}
		p.m.Unlock()
	}()

	return w
}

// Reset reinitializes a pool that has been closed/cancelled back to a working state.
// if the pool has not been closed/cancelled, nothing happens as the pool is still in
// a valid running state
func (p *limitedPool) Reset() {
	p.m.Lock()

	if !p.closed {
		p.m.Unlock()
		return
	}

	// cancelled the pool, not closed it, pool will be usable after calling initialize().
	p.initialize()
	p.m.Unlock()
}

func (p *limitedPool) closeWithError(err error) {
	p.m.Lock()

	if !p.closed {
		close(p.cancel)
		close(p.work)
		p.closed = true
	}

	for wu := range p.work {
		wu.cancelWithError(err)
	}

	p.m.Unlock()
}

// Cancel cleans up the pool workers and channels and cancels and pending
// work still yet to be processed.
// call Reset() to reinitialize the pool for use.
func (p *limitedPool) Cancel() {
	err := &ErrCancelled{s: errCancelled}
	p.closeWithError(err)
}

// Close cleans up the pool workers and channels and cancels any pending
// work still yet to be processed.
// call Reset() to reinitialize the pool for use.
func (p *limitedPool) Close() {
	err := &ErrPoolClosed{s: errClosed}
	p.closeWithError(err)
}

// Batch creates a new Batch object for queueing Work Units separate from any others
// that may be running on the pool. Grouping these Work Units together allows for individual
// Cancellation of the Batch Work Units without affecting anything else running on the pool
// as well as outputting the results on a channel as they complete.
// NOTE: Batch is not reusable, once QueueComplete() has been called it's lifetime has been sealed
// to completing the Queued items.
func (p *limitedPool) Batch() Batch {
	return newBatch(p)
}
