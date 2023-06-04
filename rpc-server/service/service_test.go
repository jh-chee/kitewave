package service

import (
	"fmt"
	"testing"
)

func TestGetChatRoom(t *testing.T) {
	tests := []struct {
		chat           string
		expectedResult string
		expectedError  error
	}{
		{
			chat:           "user1:user2",
			expectedResult: "user1:user2",
			expectedError:  nil,
		},
		{
			chat:           "user2:user1",
			expectedResult: "user1:user2",
			expectedError:  nil,
		},
		{
			chat:           "User2:User1",
			expectedResult: "User1:User2",
			expectedError:  nil,
		},
		{
			chat:           "uSEr1:UsEr2",
			expectedResult: "UsEr2:uSEr1",
			expectedError:  nil,
		},
		{
			chat:           "user1",
			expectedResult: "",
			expectedError:  fmt.Errorf("received invalid chat format user1, expecting user1:user2"),
		},
		{
			chat:           "user1:user2:user3",
			expectedResult: "",
			expectedError:  fmt.Errorf("received invalid chat format user1:user2:user3, expecting user1:user2"),
		},
	}

	for _, test := range tests {
		result, err := getChatRoom(test.chat)

		if result != test.expectedResult {
			t.Errorf("Expected result: %s, but got: %s", test.expectedResult, result)
		}

		if (err == nil && test.expectedError != nil) || (err != nil && err.Error() != test.expectedError.Error()) {
			t.Errorf("Expected error: %v, but got: %v", test.expectedError, err)
		}
	}
}
