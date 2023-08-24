package main

import (
	"log"
	router "spw/api"
	"spw/config"
)

func main() {
	config.Get()
	log.Println("hello world")

	router.Init()
}
