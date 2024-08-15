package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/srumut/todofullstack/server"
	"github.com/srumut/todofullstack/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading.env file")
	}

	pstore, err := storage.NewPostgres()
	if err != nil {
		panic(err)
	}
	s := server.NewServer(pstore)
	log.Fatal(s.Start())
}
