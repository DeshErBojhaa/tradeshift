package storage

import "github.com/DeshErBojhaa/tradeshift/graph"

// Persister exposes behaviour for underlying static types.
type Persister interface {
	GetNodes() ([]*graph.Node, error)
	InsertNode(node *graph.Node) error
	UpdateParent(curNode, targetNode *graph.Node) error
}
