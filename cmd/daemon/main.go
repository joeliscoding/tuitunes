package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

func main() {
	socketPath := "/tmp/tuitunesdaemon.sock"

	socket, err := net.Listen("unix", socketPath) // listen on Unix domain socket
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(socketPath) // clean up socket file on exit
		os.Exit(1)
	}()

	for {
		conn, err := socket.Accept() // wait for a client to connect
		if err != nil {
			log.Fatal(err)
		}

		// handle client connection
		go func(conn net.Conn) {
			defer conn.Close()

			buf := make([]byte, 4096)

			n, err := conn.Read(buf) // read data from client
			if err != nil {
				log.Fatal(err)
			}

			if strings.Contains(string(buf[:n]), "play") {
				fmt.Println("Playing audio...")
				playAudio()
			} else {
				fmt.Printf("Received: %s", string(buf[:n]))
			}

			_, err = conn.Write([]byte("Command received")) // send response to client
			if err != nil {
				log.Fatal(err)
			}
		}(conn)
	}
}

func playAudio() {
	f, err := os.Open("test.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	speaker.Play(streamer)
	select {}
}
