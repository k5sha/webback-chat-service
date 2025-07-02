package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("record already exists")

	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Messages interface {
		ListByRoomID(context.Context, string, int) ([]*Message, error)
		Create(context.Context, *Message) (*Message, error)
		// Delete(context.Context, int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Messages: &MessagesStore{db},
	}
}
