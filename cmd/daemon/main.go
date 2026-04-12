package main

import (
	"log"

	"tuitunes/internal/daemon"
)

func main() {
	if err := daemon.Run(); err != nil {
		log.Fatal(err)
	}
}
