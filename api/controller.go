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
	pathParanParentID = "par_id"
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
			Data(responseKeyErrors, err).Writer
	}
	return NewResponse(http.StatusCreated, core.MediaTypeJSON).Data(responseKeyNode, children).Writer
}

// UpdateParent changes parent of a given node. First chenge the underlying 
// persistance storage. If that succeeds, update the in memory cache.
func (c Controller) UpdateParent(req core.Request) core.ResponseWriter {
	id, ok := req.PathParam(pathParamID)
	if !ok {
		return NewResponse(http.StatusBadRequest, core.MediaTypeJSON).
			Data(responseKeyErrors, fmt.Sprintf("id: %s not found", id)).Writer
	}
	parID, ok := req.PathParam(pathParanParentID)
	if !ok {
		return NewResponse(http.StatusBadRequest, core.MediaTypeJSON).
			Data(responseKeyErrors, fmt.Sprintf("parent id: %s not found", parID)).Writer
	}
	curNode := c.g.Nodes[id]
	newPar := c.g.Nodes[parID]
	if curNode != nil || newPar != nil {
		return NewResponse(http.StatusBadRequest, core.MediaTypeJSON).
			Data(responseKeyErrors, fmt.Sprintf("invalid id: %s or parent id: %s", id, parID)).Writer
	}

	if err := c.store.UpdateParent(curNode, newPar); err != nil {
		return NewResponse(http.StatusInternalServerError, core.MediaTypeJSON).
			Data(responseKeyErrors, err).Writer
	}

	if err := c.g.UpdateParent(id, parID); err != nil {
		return NewResponse(http.StatusInternalServerError, core.MediaTypeJSON).
			Data(responseKeyErrors, err).Writer
	}
	return nil
}

// Create adds an app to storage and returns it with its unique identifier
func (c Controller) Create(req core.Request) core.ResponseWriter {
	node := graph.Node{}
	if err := req.JSON(&node); err != nil {
		return NewResponse(http.StatusBadRequest, core.MediaTypeJSON).Data(responseKeyErrors, errorBadBody).Writer
	}
	// Add validation

	// ok, messages, err := app.Validate(c.validate)
	// if err != nil {
	// 	log.Println(err)
	// 	return NewResponse(http.StatusInternalServerError, core.MediaTypeJSON).Writer
	// }

	// if !ok {
	// 	return NewResponse(http.StatusBadRequest, core.MediaTypeJSON).Data(responseKeyErrors, messages).Writer
	// }
	if node.ParID != "" {
		parent := c.g.Nodes[node.ParID]
		node.Height = parent.Height + 1
	}
	if err := c.store.InsertNode(&node); err != nil {
		return NewResponse(http.StatusInternalServerError, core.MediaTypeJSON).Writer
	}

	return NewResponse(http.StatusCreated, core.MediaTypeJSON).Data(responseKeyNode, node).Writer
}
