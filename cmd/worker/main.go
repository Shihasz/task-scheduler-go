package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"task-scheduler-go/internal/storage"
	"task-scheduler-go/internal/worker"
	"task-scheduler-go/internal/worker/executors"
)

func main() {
	// Initialize storage (same as scheduler for now)
	store := storage.NewMemoryStorage()

	// Create worker
	workerID := "worker-1"
	if len(os.Args) > 1 {
		workerID = os.Args[1]
	}

	w := worker.NewWorker(workerID, store)

	// Initialize executors
	registry := worker.NewExecutorRegistry()
	registry.Register(executors.NewPrintMessageExecutor())

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Printf("Shutdown signal received, stopping worker %s...", workerID)
		w.Stop()
	}()

	// Start the worker
	log.Printf("Starting worker %s...", workerID)
	w.Start(registry.GetExecutors())

	log.Printf("Worker %s stopped.", workerID)
}
