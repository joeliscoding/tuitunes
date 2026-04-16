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
	"tuitunes/internal/daemon/audioplayer"
	"tuitunes/internal/daemon/macos"
)

var queue = list.New()

func Run() error {
	// TODO: make socket path configurable in global config
	socketPath := "/tmp/tuitunesdaemon.sock"

	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	socket, err := net.Listen("unix", socketPath) // listen on Unix domain socket
	if err != nil {
		return err
	}
	defer socket.Close()
	defer os.Remove(socketPath)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(socketPath) // clean up socket file on exit
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

			_, err = conn.Write([]byte("Command received")) // send response to client
			if err != nil {
				log.Fatal(err)
			}

			if strings.Contains(string(buf[:n]), "play") {
				//TODO: make this more robust, maybe use JSON to send commands and data
				fmt.Println("Playing audio..." + string(buf[5:n]))
				go audioplayer.AddSong(string(buf[5:n])) // extract file path from command
			} else if strings.Contains(string(buf[:n]), "pause") {
				fmt.Println("Pausing audio...")
				go audioplayer.PauseAudio()
			} else if strings.Contains(string(buf[:n]), "volume") {
				fmt.Println("Changing volume...")
				delta := string(buf[7:n])
				deltaFloat, err := strconv.ParseFloat(delta, 64)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to parse volume delta: %v\n", err)
					return
				}
				go audioplayer.ChangeVolume(deltaFloat)
			} else {
				fmt.Printf("Received: %s", string(buf[:n]))
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
