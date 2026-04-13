package daemon

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
				playAudio(string(buf[5:n])) // extract file path from command
			} else {
				fmt.Printf("Received: %s", string(buf[:n]))
			}
		}(conn)
	}
}

func playAudio(file string) {
	f, err := os.Open(file)
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

	// nowplayinghelper is not done yet
	// err = macos.UpdateNowPlaying("TestTitle", "TestArtist", "TestAlbum")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	select {}
}
