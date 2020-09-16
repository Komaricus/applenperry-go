package main

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/route"
)

func main() {
	db := db.Init()
	e := route.Init(db)

	e.Logger.Fatal(e.Start(":5001"))
}
