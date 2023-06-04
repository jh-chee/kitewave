package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jh-chee/kitewave/rpc-server/models"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	// Create a new mock DB and get the database and mock object
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a new instance of the messageRepository with the mock DB
	repo := NewMessageRepository(db)

	// Create a sample message to save
	message := &models.Message{
		ChatRoom: "room",
		Sender:   "user",
		Body:     "hello",
	}

	// Set up expectations
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO messages").
		ExpectExec().
		WithArgs(message.ChatRoom, message.Sender, message.Body).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the Save method
	err = repo.Save(message)
	assert.NoError(t, err)

	// Verify that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
