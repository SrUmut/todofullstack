package main

import (
	"log"

	"github.com/srumut/todofullstack/server"
	"github.com/srumut/todofullstack/storage"
)

func main() {
	pstore, err := storage.NewPostgres()
	if err != nil {
		panic(err)
	}

	s := server.NewServer(pstore)
	log.Fatal(s.Start())
}

// TODO: add input check
// TODO: add fading and little bit latency to the deleting a todo
