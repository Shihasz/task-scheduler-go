package shared

import "task-scheduler-go/internal/storage"

// Global storage instance (for development/demo purposes)
// In production, we'd use dependency injection or a database.
var (
	StorageInstance storage.Storage
)

func init() {
	StorageInstance = storage.NewMemoryStorage()
}
