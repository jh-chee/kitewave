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

func TestMessageRepository_Pull(t *testing.T) {
	// Create a new instance of sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create the repository with the mock DB
	repo := &messageRepository{
		db: db,
	}

	// Define the test case
	testCases := []struct {
		name          string
		request       *models.Request
		expectedRows  *sqlmock.Rows
		expectedCount int64
		expectedError error
	}{
		{
			name: "Successful pull with cursor",
			request: &models.Request{
				ChatRoom: "a1:a2",
				Limit:    10,
				Cursor:   123456789,
				Reverse:  false,
			},
			expectedRows: sqlmock.NewRows([]string{"id", "chat_room", "sender", "body", "created"}).
				AddRow(1, "a1:a2", "John", "Hello", 1234567890).
				AddRow(2, "a1:a2", "Mary", "Hi", 1234567900),
			expectedCount: 2,
			expectedError: nil,
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Expect a query to be executed
			mock.ExpectQuery("SELECT id, chat_room, sender, body, created FROM messages").
				WillReturnRows(tc.expectedRows)

			// Call the method under test
			messages, _, err := repo.Pull(tc.request)

			// Assert the results
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedCount, int64(len(messages)))
			// Assert other expectations as needed

			// Ensure all expectations were met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestCheckCursorExistence(t *testing.T) {
	// Create a new mock database and get the mock instance
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new messageRepository instance using the mock DB
	repo := &messageRepository{db: db}

	// Define the test cursor value
	cursor := int64(123)

	// Define the expected query and result
	query := `SELECT 1 FROM messages WHERE created = ?`
	expectedRows := sqlmock.NewRows([]string{"1"}).AddRow(1)

	// Expect the query preparation and execution
	mock.ExpectPrepare(query).
		ExpectQuery().
		WithArgs(cursor).
		WillReturnRows(expectedRows)

	// Call the function under test
	err = repo.CheckCursorExistence(cursor)

	// Assert that no error occurred
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verify that all the expected queries were executed
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations: %v", err)
	}
}
