package api

import (
	"fmt"
	"net/http"

	"github.com/DeshErBojhaa/tradeshift/graph"
	"github.com/DeshErBojhaa/tradeshift/storage"
	"github.com/DeshErBojhaa/tradeshift/webber/core"
	"gopkg.in/go-playground/validator.v8"
)

const (
	pathParamID       = "id"
	pathParanParentID = "parid"
	responseKeyNode   = "node"
	responseKeyErrors = "errors"
	errorBadBody      = "invalid request body"
)

// Controller ...
type Controller struct {
	store    storage.Persister
	validate *validator.Validate
	g        *graph.Graph
}

// GetChildren returns all children of a given node. For fast response time
// we first try to return from the in memory cache. i.e. 'graph'
func (c Controller) GetChildren(req core.Request) core.ResponseWriter {
	id, ok := req.PathParam(pathParamID)
	if !ok {
		return NewResponse(http.StatusBadRequest, core.MediaTypeJSON).
			Data(responseKeyErrors, errorBadBody).Writer
	}
	children, err := c.g.GetChildren(id)
	if err != nil {
		return NewResponse(http.StatusInternalServerError, core.MediaTypeJSON).
			Data(responseKeyErrors, err.Error()).Writer
	}
	return NewResponse(http.StatusCreated, core.MediaTypeJSON).Data(responseKeyNode, children).Writer
}

// UpdateParent changes parent of a given node. First chenge the underlying
// persistence storage. If that succeeds, update the in memory cache.
func (c Controller) UpdateParent(req core.Request) core.ResponseWriter {
	id, ok := req.PathParam(pathParamID)
	if !ok {
		return NewResponse(http.StatusBadRequest, core.MediaTypeJSON).
			Data(responseKeyErrors, fmt.Sprintln("id not found in request")).Writer
	}
	parID, ok := req.PathParam(pathParanParentID)
	if !ok {
		return NewResponse(http.StatusBadRequest, core.MediaTypeJSON).
			Data(responseKeyErrors, fmt.Sprintln("parent id not found in request")).Writer
	}
	if id == parID {
		return NewResponse(http.StatusBadRequest, core.MediaTypeJSON).
			Data(responseKeyErrors, fmt.Sprintln("self cycle is not allowed")).Writer
	}

	curNode, newPar := c.g.Nodes[id], c.g.Nodes[parID]
	if curNode == nil || newPar == nil {
		return NewResponse(http.StatusBadRequest, core.MediaTypeJSON).
			Data(responseKeyErrors, fmt.Sprintf("invalid id: %s or parent id: %s", id, parID)).Writer
	}

	if err := c.store.UpdateParent(curNode, newPar); err != nil {
		return NewResponse(http.StatusInternalServerError, core.MediaTypeJSON).
			Data(responseKeyErrors, err.Error()).Writer
	}

	if err := c.g.UpdateParent(id, parID); err != nil {
		return NewResponse(http.StatusInternalServerError, core.MediaTypeJSON).
			Data(responseKeyErrors, err.Error()).Writer
	}
	return NewResponse(http.StatusCreated, core.MediaTypeJSON).Data(responseKeyNode, curNode).Writer
}

// Create adds an node to storage, updates the in-memory cache and returns the node.
func (c Controller) Create(req core.Request) core.ResponseWriter {
	node := graph.NewEmptyNode()
	if err := req.JSON(&node); err != nil {
		return NewResponse(http.StatusBadRequest, core.MediaTypeJSON).
			Data(responseKeyErrors, errorBadBody).Writer
	}

	// child node. Get height from it's parent.
	if node.ParID != "" {
		parent := c.g.Nodes[node.ParID]
		node.Height = parent.Height + 1
	}

	if err := c.store.InsertNode(&node); err != nil {
		return NewResponse(http.StatusInternalServerError, core.MediaTypeJSON).
			Data(responseKeyErrors, err.Error()).Writer
	}

	if err := c.g.EmplaceNode(&node); err != nil {
		return NewResponse(http.StatusInternalServerError, core.MediaTypeJSON).
			Data(responseKeyErrors, err.Error()).Writer
	}

	return NewResponse(http.StatusCreated, core.MediaTypeJSON).Data(responseKeyNode, node).Writer
}
