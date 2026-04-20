package daemon

import (
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/joeliscoding/tuitunes/internal/config"
	"github.com/joeliscoding/tuitunes/internal/daemon/audioplayer"
	"github.com/joeliscoding/tuitunes/internal/daemon/macos"
)

var queue = list.New()

func Run() error {
	if err := os.Remove(config.SocketPath()); err != nil && !os.IsNotExist(err) {
		return err
	}

	socket, err := net.Listen("unix", config.SocketPath()) // listen on Unix domain socket
	if err != nil {
		return err
	}
	defer socket.Close()
	defer os.Remove(config.SocketPath())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(config.SocketPath()) // clean up socket file on exit
		os.Exit(0)
	}()

	for {
		conn, err := socket.Accept() // wait for a client to connect
		if err != nil {
			return err
		}

		// handle client connection
		go func(conn net.Conn) {
			defer conn.Close()

			buf := make([]byte, 4096)

			n, err := conn.Read(buf) // read data from client
			if err != nil {
				log.Fatal(err)
			}

			recievedCommand := string(buf[:n])

			if strings.Contains(recievedCommand, "play") {
				//TODO: make this more robust, maybe use JSON to send commands and data
				fmt.Println("Playing audio..." + string(buf[5:n]))
				err = audioplayer.AddSong(string(buf[5:n]))
				if err != nil {
					log.Fatal(err)
				}
			} else if strings.Contains(recievedCommand, "pause") {
				fmt.Println("Pausing audio...")
				audioplayer.PauseAudio()
			} else if strings.Contains(recievedCommand, "volume") {
				fmt.Println("Changing volume...")
				delta := string(buf[7:n])
				deltaFloat, err := strconv.ParseFloat(delta, 64)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to parse volume delta: %v\n", err)
					return
				}
				audioplayer.ChangeVolume(deltaFloat)
			} else if strings.Contains(recievedCommand, "ping") {
				_, err = conn.Write([]byte("pong")) // send response to client
				if err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Printf("Received: %s", recievedCommand)
			}
		}(conn)
	}
}

func addToQueue(file string) {
	queue.PushBack(file)
}

func updateNowPlaying(file string) {
	// nowplayinghelper is not done yet
	err := macos.UpdateNowPlaying(file)
	if err != nil {
		log.Fatal(err)
	}
}
