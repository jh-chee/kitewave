package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jh-chee/kitewave/rpc-server/models"
	"github.com/rs/zerolog/log"
)

const (
	timeout = 5 * time.Minute
)

type MessageRepository interface {
	Save(message *models.Message) error
}

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (r *messageRepository) Save(message *models.Message) error {
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Start a new transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}

	// Prepare the INSERT statement within the transaction
	stmt, err := tx.Prepare("INSERT INTO messages (chat_room, sender, body) VALUES (?, ?, ?)")
	if err != nil {
		rollback(tx)
		return fmt.Errorf("fail to prepare save query")
	}
	defer stmt.Close()

	// Execute the INSERT statement with parameter values
	_, err = stmt.ExecContext(ctx, message.ChatRoom, message.Sender, message.Body)
	if err != nil {
		rollback(tx)
		return fmt.Errorf("fail to exec save query")
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

func rollback(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		log.Error().Err(err).Msg("rollback unsuccessful")
	}
}
