package model

import (
	"bytelyon-functions/pkg/db"
	"sync"
)

type NodeTree struct {
	mu    sync.Mutex
	wg    sync.WaitGroup
	s3    db.S3
	graph map[string]map[string][]*Node
}

func NewNodeTree() *NodeTree {
	return &NodeTree{
		s3:    db.NewS3(),
		graph: make(map[string]map[string][]*Node),
	}
}
