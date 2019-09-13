package graph

import (
	"errors"
	"fmt"
)

// TODO: Cover tests

// ErrDuplicateID triggers when given id already exists
var ErrDuplicateID = errors.New("id alreary exists")

// ErrInvalidParentID ...
var ErrInvalidParentID = errors.New("parent not found")

// Node is the building block of Graph.
type Node struct {
	ID       string // Unique
	Parent   *Node
	Children map[string]*Node
	Height   int
}

// Graph is the in memory representation of the company jierarcy.
type Graph struct {
	UpdateDB bool
	Root     *Node
	Nodes    map[string]*Node
}

var hiararchy Graph

func init() {
	// Make it from the database
	hiararchy = Graph{
		UpdateDB: false,
		Root:     nil,
		Nodes:    make(map[string]*Node),
	}
}

// CreateNode creates a new in memory node. Make it persist in the database
func CreateNode(id, parent string) (*Node, error) {
	if _, ok := hiararchy.Nodes[id]; ok {
		return nil, ErrDuplicateID
	}

	var parNode *Node
	var ok bool
	if parNode, ok = hiararchy.Nodes[parent]; !ok {
		return nil, ErrInvalidParentID
	}

	newNode := &Node{
		ID:       id,
		Parent:   parNode,
		Children: make(map[string]*Node),
	}
	// Update the graph
	parNode.Children[id] = newNode

	// Update the persistance
	// err := db.AddNewNode(parNode, newNode)  // Node pathaye update korbo, naki valu gula pathaye node return koerbo??
	return newNode, nil
}

// UpdateParent sets the parent to parameter 'newPar'
func UpdateParent(id, newPar string) error {
	var curNode, parNode *Node
	var ok bool
	if curNode, ok = hiararchy.Nodes[id]; !ok {
		return fmt.Errorf("invalid id %s", id)
	}
	if curNode == hiararchy.Root {
		return fmt.Errorf("can not update parent of the root")
	}

	if parNode, ok = hiararchy.Nodes[newPar]; !ok {
		return fmt.Errorf("invalid parent id")
	}
	// Remove reference from the old parent
	delete(curNode.Parent.Children, id)
	// Add reference to the new parent
	parNode.Children[id] = curNode

	// Update persistance

	return nil
}

// GetChildren returns all the childrens of a given node
func GetChildren(id string) ([]*Node, error) {
	var curNode *Node
	var ok bool

	if curNode, ok = hiararchy.Nodes[id]; !ok {
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
