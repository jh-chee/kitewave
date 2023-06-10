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
	Pull(req *models.Request) ([]*models.Message, int64, error)
	CheckCursorExistence(cursor int64) error
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

func (r *messageRepository) Pull(req *models.Request) ([]*models.Message, int64, error) {
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	sortOrder := "ASC"
	cursorCondition := fmt.Sprintf(" AND created >= %d", req.Cursor)
	if req.Reverse {
		sortOrder = "DESC"
		cursorCondition = fmt.Sprintf(" AND created <= %d", req.Cursor)
	}

	// No cursor given
	if req.Cursor == 0 {
		cursorCondition = ""
	}

	query := fmt.Sprintf(`SELECT id, chat_room, sender, body, created FROM messages WHERE chat_room='%s'%s ORDER BY created %s, id %s LIMIT %d`,
		req.ChatRoom, cursorCondition, sortOrder, sortOrder, req.Limit+1,
	)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	messages := []*models.Message{}
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(&message.Id, &message.ChatRoom, &message.Sender, &message.Body, &message.Created)
		if err != nil {
			return nil, 0, err
		}
		messages = append(messages, message)
	}

	nextCursor := int64(0)
	if len(messages) > int(req.Limit) {
		lastMessage := messages[req.Limit]
		nextCursor = lastMessage.Created
		messages = messages[:req.Limit] // Remove the extra message beyond the limit
	}

	return messages, nextCursor, nil
}

func (r *messageRepository) CheckCursorExistence(cursor int64) error {
	// cursor 0 is a default start cursor
	if cursor == 0 {
		return nil
	}

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Prepare the query statement
	query := `SELECT 1 FROM messages WHERE created = ? LIMIT 1`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("fail to prepare exist query: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.QueryContext(ctx, cursor)
	if err != nil {
		return fmt.Errorf("fail to exec exist query: %w", err)
	}
	defer rows.Close()

	// Check if a row exists
	if rows.Next() {
		return nil
	}

	return fmt.Errorf("cursor does not exist")
}

func rollback(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		log.Error().Err(err).Msg("rollback unsuccessful")
	}
}
