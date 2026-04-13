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
	if len(os.Args) < 3 || os.Args[1] != "daemon" {
		fmt.Fprintf(os.Stderr, "usage: go run cmd/cli/main.go daemon /path/to/file.mp3\n")
		os.Exit(1)
	}

	filePath := strings.Join(os.Args[2:], " ")
	if err := sendPlayCommand(filePath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to send play command: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sent play command for %q\n", filePath)

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
