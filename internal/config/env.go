package config

import "os"

const (
	defaultSocketPath = "/tmp/tuitunesdaemon.sock"
	defaultDaemonPath = "/usr/local/bin/tuitunes-daemon"
)

func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func SetEnvDefaults() {
	if os.Getenv("TUITUNES_SOCK") == "" {
		_ = os.Setenv("TUITUNES_SOCK", defaultSocketPath)
	}

	if os.Getenv("TUITUNES_DAEMON") == "" {
		_ = os.Setenv("TUITUNES_DAEMON", defaultDaemonPath)
	}
}

func SocketPath() string {
	return getEnvOrDefault("TUITUNES_SOCK", defaultSocketPath)
}

func DaemonPath() string {
	return getEnvOrDefault("TUITUNES_DAEMON", defaultDaemonPath)
}
