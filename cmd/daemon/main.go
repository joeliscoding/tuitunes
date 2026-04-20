package main

import (
	"log"

	"github.com/joeliscoding/tuitunes/internal/daemon"
)

func main() {
	if err := daemon.Run(); err != nil {
		log.Fatal(err)
	}
}
