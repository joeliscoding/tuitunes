package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// TODO: make socket path configurable in global config
const socketPath = "/tmp/tuitunesdaemon.sock"

func main() {

	filePath := strings.Join(os.Args[2:], " ")
	switch os.Args[1] {
	case "play":
		if err := sendPlayCommand(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to send play command: %v\n", err)
			os.Exit(1)
		}
	case "pause":
		if err := sendPauseCommand(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to send pause command: %v\n", err)
			os.Exit(1)
		}
	case "volume":
		delta := os.Args[2]
		if err := sendVolumeCommand(delta); err != nil {
			fmt.Fprintf(os.Stderr, "failed to send volume command: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "usage: go run cmd/cli/main.go play /path/to/file.mp3\n")
		os.Exit(1)
	}

	fmt.Printf("Sent %s command for %q\n", os.Args[1], filePath)

}

func sendPlayCommand(filePath string) error {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("play " + filePath)); err != nil {
		return err
	}

	_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err == nil && n > 0 {
		fmt.Printf("Daemon: %s\n", strings.TrimSpace(string(buf[:n])))
	}

	return nil
}

func sendPauseCommand() error {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("pause")); err != nil {
		return err
	}

	_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err == nil && n > 0 {
		fmt.Printf("Daemon: %s\n", strings.TrimSpace(string(buf[:n])))
	}

	return nil
}

func sendVolumeCommand(delta string) error {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("volume " + delta)); err != nil {
		return err
	}

	_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err == nil && n > 0 {
		fmt.Printf("Daemon: %s\n", strings.TrimSpace(string(buf[:n])))
	}

	return nil
}
