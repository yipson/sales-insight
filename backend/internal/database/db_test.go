package database

import (
	"context"
	"testing"
	"time"
)

func TestNewDB(t *testing.T) {
	cfg := DefaultConfig()

	db, err := NewDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.Health(ctx); err != nil {
		t.Fatalf("health check failed: %v", err)
	}

	t.Log("database connection successful")
}
