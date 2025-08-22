package repository

import (
	"testing"

	"github.com/nathaliaoliveira/goapp/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGenerateContentHash(t *testing.T) {
	repo := &eventRepository{}

	event := &domain.EmailEvent{
		Type:      "sent",
		Email:     "user@example.com",
		Site:      "site-a.com",
		Timestamp: "2025-08-20T10:30:00Z",
	}

	hash := repo.generateContentHash(event)

	assert.NotEmpty(t, hash)
	assert.Len(t, hash, 64)

	hash2 := repo.generateContentHash(event)
	assert.Equal(t, hash, hash2)
}

func TestGenerateContentHash_DifferentEvents(t *testing.T) {
	repo := &eventRepository{}

	event1 := &domain.EmailEvent{
		Type:      "sent",
		Email:     "user@example.com",
		Site:      "site-a.com",
		Timestamp: "2025-08-20T10:30:00Z",
	}

	event2 := &domain.EmailEvent{
		Type:      "open",
		Email:     "user@example.com",
		Site:      "site-a.com",
		Timestamp: "2025-08-20T10:30:00Z",
	}

	hash1 := repo.generateContentHash(event1)
	hash2 := repo.generateContentHash(event2)

	assert.NotEqual(t, hash1, hash2)
}

func TestGenerateContentHash_EmptyFields(t *testing.T) {
	repo := &eventRepository{}

	event := &domain.EmailEvent{
		Type:      "",
		Email:     "",
		Site:      "",
		Timestamp: "",
	}

	hash := repo.generateContentHash(event)

	assert.NotEmpty(t, hash)
	assert.Len(t, hash, 64)
	
	event2 := &domain.EmailEvent{
		Type:      "sent",
		Email:     "user@example.com",
		Site:      "site-a.com",
		Timestamp: "2025-08-20T10:30:00Z",
	}
	hash2 := repo.generateContentHash(event2)
	assert.NotEqual(t, hash, hash2)
} 