package main

import "time"

type StorageDriver interface {
	// Set caches the item.
	Set(string, interface{}, time.Duration) error
	// Get retrieves the object for the given key.
	Get(string) (interface{}, error)
	// Delete deletes the object and key for the given key.
	Delete(string) (interface{}, error)
	// Exists reports whether object for the given key exists.
	Exists(string) (int64, error)
	// Close closes the storage client, releasing any open resources.
	Close() error
}
