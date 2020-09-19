package main

import (
	"github.com/applenperry-go/config"
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/route"
)

func main() {
	configuration := config.GetConfig()
	err := db.Init(configuration)
	if err != nil {
		panic("Failed to connect to database!")
	}
	r := route.Init(configuration)

	r.Run(":5001")
}
