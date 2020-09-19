package main

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/route"
)

func main() {
	err := db.Init()
	if err != nil {
		panic("Failed to connect to database!")
	}
	r := route.Init()

	r.Run(":5001")
}
