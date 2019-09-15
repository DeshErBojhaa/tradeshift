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

// Node is the building block of Graph. Node ID is chosen as
// the unique identifier for each node. Which is not the best
// practice, but will serve the given problem sufficiently.
type Node struct {
	ID       string // Unique
	ParID    string
	Children map[string]*Node
	Height   int
}

// Graph is in memory representation of hierarcy.
type Graph struct {
	UpdateDB bool
	Root     *Node
	Nodes    map[string]*Node
}

// Initialize creates a new graph object
func Initialize(nodes []*Node) (*Graph, error) {
	nodeMap := make(map[string]*Node)
	root := &Node{}

	for _, node := range nodes {
		nodeMap[node.ID] = node
		if node.ParID == "" {
			root = node
		}
	}
	return &Graph{
		Root:  root,
		Nodes: nodeMap,
	}, nil
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
func (g *Graph) CreateNode(id, parent string) (*Node, error) {
	if _, ok := g.Nodes[id]; ok {
		return nil, ErrDuplicateID
	}

	parNode, ok := g.Nodes[parent]
	if !ok {
		return nil, ErrInvalidParentID
	}

	newNode := &Node{
		ID:       id,
		ParID:    parent,
		Children: make(map[string]*Node),
	}
	// Update the graph
	parNode.Children[id] = newNode

	// Update the persistance
	// err := db.AddNewNode(parNode, newNode)  // Node pathaye update korbo, naki valu gula pathaye node return koerbo??
	return newNode, nil
}

// UpdateParent sets the parent to parameter 'newPar'
func (g *Graph) UpdateParent(id, newPar string) error {
	var curNode, parNode *Node
	var ok bool
	if curNode, ok = g.Nodes[id]; !ok {
		return fmt.Errorf("invalid id %s", id)
	}
	if curNode == g.Root {
		return fmt.Errorf("can not update parent of the root")
	}

	if parNode, ok = g.Nodes[newPar]; !ok {
		return fmt.Errorf("invalid parent id")
	}
	// Remove reference from the old parent
	prevPar := g.Nodes[curNode.ParID]
	delete(prevPar.Children, id)
	// Add reference to the new parent
	parNode.Children[id] = curNode

	// Update persistance

	return nil
}

// GetChildren returns all the childrens of a given node
func (g *Graph) GetChildren(id string) ([]*Node, error) {
	var curNode *Node
	var ok bool

	if curNode, ok = g.Nodes[id]; !ok {
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

// $ CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -installsuffix cgo -ldflags '-s' -o server
// Build with ^. This creats a static binary
