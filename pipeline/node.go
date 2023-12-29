package pipeline

import (
	"errors"
	"fmt"
	"github.com/huahuayu/kit/logger"
	"runtime"
	"sync"
)

type INode interface {
	getName() string
	SetNext(node INode)
	getNext() []INode
	jobReceiver(job any) error
	start()
	stop()
	wait()
	latestJob() any
}

// Node defines a node in a pipeline.
// name is the node identifier.
// workFn is the function actually do the work.
// ingestCh is the channel to receive job from previous node.
// jobCh is the channel to receive job from ingestCh, it's a buffered channel.
// doneCh is the channel to notify the node to stop.
// workerPoolSize is the size of the worker pool.
// lastJob is the last job been processed by the node, it can be used as breakpoint.
// next is the next node in the pipeline.
type Node[T any] struct {
	name           string
	workFn         func(T) (any, error)
	ingestCh       chan T
	jobCh          chan T
	doneCh         chan struct{}
	workerPoolSize int
	lastJob        T
	once           sync.Once
	wg             sync.WaitGroup
	next           []INode
}

var _ INode = (*Node[any])(nil)

// NewNode creates a new node
func NewNode[T any](name string, workFunc func(T) (any, error), jobPoolSize int, workerPoolSize int) INode {
	var node = &Node[T]{
		name:           name,
		workFn:         workFunc,
		ingestCh:       make(chan T),
		jobCh:          make(chan T, jobPoolSize),
		doneCh:         make(chan struct{}),
		workerPoolSize: workerPoolSize,
	}
	return node
}

// NewDefaultNode creates a new node with default jobPoolSize and workerPoolSize
func NewDefaultNode[T any](name string, workFn func(T) (any, error)) INode {
	cpuNum := runtime.NumCPU()
	return NewNode(name, workFn, cpuNum, cpuNum)
}

// getName returns the node name
func (n *Node[T]) getName() string {
	return n.name
}

// SetNext set the next node in the pipeline
func (n *Node[T]) SetNext(node INode) {
	if n.next == nil {
		n.next = make([]INode, 0)
	}
	n.next = append(n.next, node)
}

// getNext returns the next node in the pipeline
func (n *Node[T]) getNext() []INode {
	return n.next
}

// jobReceiver receive job from previous node
func (n *Node[T]) jobReceiver(job any) error {
	if job == nil {
		return nil
	}
	switch t := job.(type) {
	case T:
		n.ingestCh <- job.(T)
	default:
		return errors.New(fmt.Sprintf("job type %T is not accept by [%s]", t, n.name))
	}
	return nil
}

// start listen the job from previous node, start all workers
func (n *Node[T]) start() {
	go func() {
		for {
			select {
			case job := <-n.ingestCh:
				n.jobCh <- job
			case <-n.doneCh:
				close(n.jobCh)
				logger.Logger.Infof("[%s] received cancel signal, stop receive new job...", n.name)
				return
			}
		}
	}()
	// Start [workerPoolSize] of worker to do the job
	n.wg.Add(n.workerPoolSize)
	for i := 0; i < n.workerPoolSize; i++ {
		go n.work(i)
	}
}

func (n *Node[T]) stop() {
	n.doneCh <- struct{}{}
}

func (n *Node[T]) work(workerId int) {
	logger.Logger.Debugf("[%s] worker %d starting...", n.name, workerId)
	defer n.wg.Done()
	for job := range n.jobCh {
		// Atomic store the last job (brings lock)
		//atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&n.lastJob)), unsafe.Pointer(&job))
		n.lastJob = job
		// Working on the job
		result, err := n.workFn(job)
		if err != nil {
			logger.Logger.Errorf("[%s] worker %d failed with: %s", n.name, workerId, err)
			continue
		}
		logger.Logger.Debugf("[%s] worker %d completed the job successful", n.name, workerId)
		// Pipe to next node
		if n.next != nil {
			for _, next := range n.next {
				logger.Logger.Debugf("[%s] worker %d passed job to %s", n.name, workerId, next.getName())
				next.jobReceiver(result)
			}
		}
	}
	logger.Logger.Debugf("[%s] worker %d get off work", n.name, workerId)
}

func (n *Node[T]) wait() {
	n.once.Do(func() {
		n.wg.Wait()
	})
}

func (n *Node[T]) latestJob() any {
	//return *(*T)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&n.lastJob))))
	return n.lastJob
}
