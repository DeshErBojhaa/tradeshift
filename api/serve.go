package api

import (
	"log"

	"github.com/DeshErBojhaa/tradeshift/graph"
	"github.com/DeshErBojhaa/tradeshift/storage/mysql"
	"github.com/DeshErBojhaa/tradeshift/webber"
	"github.com/DeshErBojhaa/tradeshift/webber/core"
	"gopkg.in/go-playground/validator.v8"
)

// Serve the API server
func Serve(listenAddress, connString string, v *validator.Validate) error {
	db, err := mysql.NewMySQLStore(connString)
	if err != nil {
		log.Fatal(err)
	}
	nodes, err := db.GetNodes()
	if err != nil {
		log.Fatal(err)
	}

	gp, err := graph.Initialize(nodes)
	if err != nil {
		log.Fatal(err)
	}
	controller := Controller{
		store:    db,
		validate: v,
		g:        gp,
	}

	s := webber.NewServer(listenAddress, core.MediaTypeJSON)
	s.POST("/node", controller.Create)
	s.UPDATE("/{id}?par={parid}", controller.UpdateParent)
	s.GET("/children/{id}", controller.GetChildren)

	return s.Serve()
}
