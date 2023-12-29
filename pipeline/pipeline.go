package pipeline

import (
	"errors"
	"github.com/huahuayu/kit/logger"
	"sync"
)

type Pipeline struct {
	root   INode
	once   sync.Once
	doneCh chan struct{} // doneCh is used to notify producer to stop
}

func NewPipeline(root INode) (*Pipeline, error) {
	if root == nil {
		return nil, errors.New("nil root node")
	}
	p := &Pipeline{root: root, once: sync.Once{}, doneCh: make(chan struct{}, 1)}
	if p.hasCycle() {
		return nil, errors.New("cycle not allowed in pipeline")
	}
	return p, nil
}

func (p *Pipeline) JobReceiver(job any) error {
	return p.root.jobReceiver(job)
}

// LatestJob returns the latest job been processed by the pipeline.
func (p *Pipeline) LatestJob() any {
	return p.root.latestJob()
}

// Start starts the pipeline
func (p *Pipeline) Start() chan struct{} {
	p.dfs(func(node INode) {
		logger.Logger.Debug(node.getName(), " start")
		node.start()
	})
	return p.doneCh
}

func (p *Pipeline) Stop() {
	go p.once.Do(func() {
		p.doneCh <- struct{}{}
		p.dfs(func(node INode) {
			logger.Logger.Debug(node.getName(), " stop")
			node.stop()
			node.wait()
		})
	})
}

// Wait waits for all nodes to finish.
func (p *Pipeline) Wait() {
	p.dfs(func(node INode) {
		logger.Logger.Debug(node.getName(), " wait")
		node.wait()
	})
}

// dfs route the node graph and execute fn on each node
func (p *Pipeline) dfs(fn func(node INode)) {
	visited := make(map[INode]bool)
	p.dfsHelper(p.root, visited, fn)
}

func (p *Pipeline) dfsHelper(node INode, visited map[INode]bool, fn func(node INode)) {
	if visited[node] {
		return
	}
	visited[node] = true
	fn(node)
	for _, next := range node.getNext() {
		p.dfsHelper(next, visited, fn)
	}
}

// hasCycle checks if the pipeline has cycle
func (p *Pipeline) hasCycle() bool {
	visited := make(map[INode]bool)
	recStack := make(map[INode]bool)
	for _, node := range p.getAllNodes() {
		if !visited[node] {
			if p.hasCycleHelper(node, visited, recStack) {
				return true
			}
		}
	}
	return false
}

func (p *Pipeline) hasCycleHelper(node INode, visited, recStack map[INode]bool) bool {
	visited[node] = true
	recStack[node] = true

	for _, next := range node.getNext() {
		if !visited[next] && p.hasCycleHelper(next, visited, recStack) {
			return true
		} else if recStack[next] {
			return true
		}
	}
	recStack[node] = false
	return false
}

func (p *Pipeline) getAllNodes() []INode {
	if p.root == nil {
		return nil
	}
	nodes := make([]INode, 0)
	visited := make(map[INode]bool)
	p.dfsHelper(p.root, visited, func(node INode) {
		nodes = append(nodes, node)
	})
	return nodes
}
