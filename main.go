package main

import (
	"log"
	"os"

	"github.com/DeshErBojhaa/tradeshift/api"
	"gopkg.in/go-playground/validator.v8"
)

func main() {
	v := validator.New(&validator.Config{
		TagName:      "validate",
		FieldNameTag: "json",
	})

	if err := api.Serve(":8080", os.Getenv("MYSQL_CONN"), v); err != nil {
		log.Fatal(err)
	}
}

// $ CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -installsuffix cgo -ldflags '-s' -o tradeshift
// Build with ^. This creats a static binary
