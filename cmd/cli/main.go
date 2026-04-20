package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/joeliscoding/tuitunes/internal/config"
)

func main() {
	config.SetEnvDefaults() // Load environment variables from env.go file

	// Check if daemon is already running by trying to connect to the socket
	if !checkDaemonStatus() {
		fmt.Println("Daemon is not running. Starting daemon...")
		cmd := exec.Command(config.DaemonPath())
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to start daemon: %v\n", err)
		}
		if cmd.Process != nil {
			_ = cmd.Process.Release()
		}

		// Wait for the daemon to start
		waitForDaemon()
	} else {
		fmt.Println("Daemon is already running.")
	}

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
		fmt.Fprintf(os.Stderr, "usage: go run ./cmd/cli play /path/to/file.mp3\n")
		os.Exit(1)
	}

	fmt.Printf("Sent %s command for %q\n", os.Args[1], filePath)

}

func checkDaemonStatus() bool {
	conn, err := net.Dial("unix", config.SocketPath())
	if err != nil {
		return false
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("ping")); err != nil {
		return false
	}

	_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return false
	}

	response := strings.TrimSpace(string(buf[:n]))
	return strings.Contains(response, "pong")
}

func waitForDaemon() {
	timeout := time.After(5 * time.Second)
	tick := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-timeout:
			fmt.Fprintln(os.Stderr, "timed out waiting for daemon to start")
			os.Exit(1)
		case <-tick:
			if checkDaemonStatus() {
				return
			}
		}
	}
}

func sendPlayCommand(filePath string) error {
	conn, err := net.Dial("unix", config.SocketPath())
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
	conn, err := net.Dial("unix", config.SocketPath())
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
	conn, err := net.Dial("unix", config.SocketPath())
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
