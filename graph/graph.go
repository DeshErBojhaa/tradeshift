// Package graph is inmemory representation for company hierarchy. All read requests
// are served from this package. It's consumers responsibility to update this package
// after any write call to persistant layer.
package graph

import (
	"errors"
	"fmt"
)

// ErrDuplicateID triggers when given id already exists
var ErrDuplicateID = errors.New("id alreary exists")

// ErrInvalidParentID ...
var ErrInvalidParentID = errors.New("parent not found")

// Node is the building block of Graph. Node ID is chosen as
// the unique identifier for each node. Which is not the best
// practice, but will serve the given problem sufficiently.
type Node struct {
	ID       string `json:"id"`
	ParID    string `json:"pid"`
	Height   int    `json:"height"`
	Children map[string]*Node
}

// Graph is in memory representation of hierarcy.
type Graph struct {
	Root  *Node
	Nodes map[string]*Node
}

// Initialize creates a new graph object
func Initialize(nodes []*Node) (*Graph, error) {
	nodeMap := make(map[string]*Node)
	g := Graph{}
	for _, node := range nodes {
		nodeMap[node.ID] = node
		if node.ParID == "" {
			g.Root = node
		}
	}
	g.Nodes = nodeMap
	return &g, nil
}

// NewEmptyNode ...
func NewEmptyNode() Node {
	return Node{Children: make(map[string]*Node), Height: 0}
}

// EmplaceNode emplaces the given node into the graph. Updates the parent child relationship.
func (g *Graph) EmplaceNode(node *Node) error {
	if _, ok := g.Nodes[node.ID]; ok {
		return ErrDuplicateID
	}

	parNode, ok := g.Nodes[node.ParID]
	if !ok && g.Root != nil {
		return ErrInvalidParentID
	}
	if g.Root == nil {
		g.Root = node
	}
	g.Nodes[node.ID] = node
	if parNode != nil {
		parNode.Children[node.ID] = node
	}
	return nil
}

// UpdateParent sets the parent to parameter 'newPar'
func (g *Graph) UpdateParent(id, newPar string) error {
	var curNode, newParNode *Node
	var ok bool
	if curNode, ok = g.Nodes[id]; !ok {
		return fmt.Errorf("invalid id %s", id)
	}
	if curNode == g.Root {
		return fmt.Errorf("can not update parent of the root")
	}

	if newParNode, ok = g.Nodes[newPar]; !ok {
		return fmt.Errorf("invalid parent id")
	}
	prevParNode := g.Nodes[curNode.ParID]

	// 1. Move children of cur node one level up
	for _, node := range curNode.Children {
		node.Height--
		node.ParID = prevParNode.ID
		prevParNode.Children[node.ID] = node
	}

	// 2. Remove all childs of cur node
	curNode.Children = make(map[string]*Node)

	// 3. Remove cur node from it's prev parent
	delete(prevParNode.Children, id)

	// 4. Add cur node to it's new parent
	newParNode.Children[id] = curNode

	// 5. Set cur nodes parent to right value. Update height
	curNode.ParID = newPar
	curNode.Height = newParNode.Height + 1
	return nil
}

// GetChildren returns all the childrens of a given node
func (g *Graph) GetChildren(id string) ([]*Node, error) {
	curNode, ok := g.Nodes[id]
	if !ok {
		return nil, fmt.Errorf("invalid id %s", id)
	}
	children := make([]*Node, len(curNode.Children))
	i := 0
	for _, v := range curNode.Children {
		children[i] = v
		i++
	}
	return children, nil
}
