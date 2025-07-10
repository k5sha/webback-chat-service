package store

import (
	"context"
	"database/sql"
	"time"
)

type Message struct {
	ID      string    `json:"id"`
	RoomID  string    `json:"room_id"`
	Sender  string    `json:"sender_id"`
	Content string    `json:"content"`
	Type    string    `json:"type"`
	SentAt  time.Time `json"sent_at"`
}

type MessagesStore struct {
	db *sql.DB
}

func (s *MessagesStore) Create(ctx context.Context, message *Message) (*Message, error) {
	query := `INSERT INTO messages (room_id, sender, content, type) VALUES ($1, $2, $3, $4) RETURNING id, sent_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		message.RoomID,
		message.Sender,
		message.Content,
		message.Type,
	).Scan(
		&message.ID,
		&message.SentAt,
	)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *MessagesStore) ListByRoomID(ctx context.Context, roomID string, limit int) ([]*Message, error) {
	query := `
        SELECT id, room_id, sender, content, type, sent_at 
        FROM messages 
        WHERE room_id = $1 
        ORDER BY sent_at ASC 
        LIMIT $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, roomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(
			&msg.ID,
			&msg.RoomID,
			&msg.Sender,
			&msg.Content,
			&msg.Type,
			&msg.SentAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

// func (s *MessagesStore) Update(ctx context.Context, tx *sql.Tx, user *User) error {
// 	query := `UPDATE users SET username = $1, email = $2, is_activate = $3 WHERE id = $4`

// 	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
// 	defer cancel()

// 	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *MessagesStore) Delete(ctx context.Context, id int64) error {
// 	query := `DELETE FROM users WHERE id = $1`

// 	ctx, cancel := context.WithTimeout(ctx, 3*QueryTimeoutDuration)
// 	defer cancel()

// 	_, err := tx.ExecContext(ctx, query, id)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
