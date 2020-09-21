package main

import (
	"github.com/applenperry-go/config"
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/route"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	configuration := config.GetConfig()
	if err := db.Init(configuration); err != nil {
		panic("Failed to connect to database!")
	}
	r := route.Init(configuration)

	r.Run(":5001")
}
