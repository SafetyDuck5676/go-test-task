package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	port := flag.String("port", "8080", "port to listen on")
	maxQueues := flag.Int("max-queues", 10, "maximum number of queues")
	defaultSize := flag.Int("default-size", 100, "default queue size")
	timeout := flag.Int("timeout", 5, "default timeout in seconds")
	flag.Parse()

	queueManager := NewQueueManager(*maxQueues, *defaultSize)
	server := NewServer(queueManager, *port, time.Duration(*timeout)*time.Second)

	log.Printf("Starting server on port %s", *port)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
