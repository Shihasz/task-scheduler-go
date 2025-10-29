# Distributed Task Scheduler

A distributed task scheduling system built in Go that allows asynchronous execution of tasks across multiple workers.

## Architecture

- **Server**: Combined scheduler and worker with HTTP API
- **Storage**: In-memory task storage (PostgreSQL coming soon)
- **Executors**: Plugins for different task types

## Running the project

```bash
# Run the combined server (scheduler + worker)
go run cmd/server/main.go
```
