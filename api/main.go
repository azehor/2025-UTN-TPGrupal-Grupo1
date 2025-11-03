package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"quepc/api/cmd/server"
	"quepc/api/config"
)

func main() {
	// Centralize config: only load .env in development via config.LoadConfig
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	// Build address safely using resolved host/port
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	fmt.Printf("Starting API on %s (env: %s)\n", addr, cfg.Env)

	if err := server.Start(addr); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
