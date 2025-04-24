package bot

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ikebastuz/tgn/types"
)

func TestGetDialogState(t *testing.T) {
	t.Run("should create initial dialog state", func(t *testing.T) {
		store := NewInMemoryStore()
		want := types.State(types.STATE_INITIAL)
		got, _ := getDialogState(123, store)
		if got.State != want {
			t.Errorf("expected state %v, got %v", want, got.State)
		}
	})
}

func TestCreateReply(t *testing.T) {
	const USER_ID = 123
	var FROM = types.From{
		ID:       int64(USER_ID),
		USERNAME: "hello",
	}

	t.Run("state - INITIAL and it is not a 'connect' message - should ask user to forward connection message", func(t *testing.T) {
		store := NewInMemoryStore()
		want := []types.ReplyDTO{
			{
				Message: types.ReplyMessage{
					UserID:      FROM.ID,
					Message:     createConnectionMessage(FROM),
					ReplyMarkup: nil,
				},
				NextState: &types.DialogState{
					State: types.WAITING_FOR_CONNECT,
				},
			},
			{
				Message: types.ReplyMessage{
					UserID:      FROM.ID,
					Message:     MESSAGE_FORWARD_CONNECTION_02,
					ReplyMarkup: nil,
				},
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				From: FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})
	t.Run("state - INITIAL and is a 'connect' message to yourself - should show an error", func(t *testing.T) {
		store := NewInMemoryStore()
		want := []types.ReplyDTO{
			{
				Message: types.ReplyMessage{
					UserID:      FROM.ID,
					Message:     MESSAGE_YOU_CANT_CONNECT_TO_YOURSELF,
					ReplyMarkup: nil,
				},
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				Text: fmt.Sprintf("/connect %d", USER_ID),
				From: FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})
	t.Run("state - INITIAL and is a 'connect' message to non-existent user - should show an error", func(t *testing.T) {
		store := NewInMemoryStore()
		want := []types.ReplyDTO{
			{
				Message: types.ReplyMessage{
					UserID:      FROM.ID,
					Message:     MESSAGE_NO_SUCH_USER_IS_AWATING,
					ReplyMarkup: nil,
				},
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				Text: "/connect 100500",
				From: FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})
	t.Run("state - WAITING for connection - should tell about waiting for connection", func(t *testing.T) {
		store := NewInMemoryStore()
		store.SetDialogState(USER_ID, &types.DialogState{State: types.WAITING_FOR_CONNECT})
		var FROM = types.From{
			ID:       int64(USER_ID),
			USERNAME: "hello",
		}
		want := []types.ReplyDTO{
			{
				Message: types.ReplyMessage{
					UserID:      FROM.ID,
					Message:     MESSAGE_WAITING_FOR_CONNECTION,
					ReplyMarkup: nil,
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      "",
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})
}

func TestGetSenderId(t *testing.T) {
	var FROM = &types.From{
		ID:       int64(123),
		USERNAME: "hello",
	}

	t.Run("should return correct sender id from Message", func(t *testing.T) {
		update := types.TelegramUpdate{
			Message: types.Message{
				From: *FROM,
			},
		}
		got, _ := getUserData(update)
		if !reflect.DeepEqual(got, FROM) {
			t.Errorf("expected state %v, got %v", FROM, got)
		}
	})

	t.Run("should return correct sender id from Callback", func(t *testing.T) {
		update := types.TelegramUpdate{
			CallbackQuery: types.CallbackQuery{
				From: *FROM,
			},
		}
		got, _ := getUserData(update)
		if !reflect.DeepEqual(got, FROM) {
			t.Errorf("expected state %v, got %v", FROM, got)
		}
	})

	t.Run("should throw an error if no sender id exists", func(t *testing.T) {
		update := types.TelegramUpdate{}
		_, err := getUserData(update)
		if err == nil || err != ErrorNoSenderIdFound {
			t.Errorf("expected error")
		}
	})
}

func TestGetConnectionId(t *testing.T) {
	t.Run("should not have id if message has no '/connect'", func(t *testing.T) {
		update := types.TelegramUpdate{}
		_, isConnection := getConnectionId(&update)
		if isConnection != false {
			t.Errorf("should be false")
		}
	})

	t.Run("should not have id if message has '/connect' but no id", func(t *testing.T) {
		update := types.TelegramUpdate{
			Message: types.Message{
				Text: " /connect ",
			},
		}
		_, isConnection := getConnectionId(&update)
		if isConnection != false {
			t.Errorf("should be false")
		}
	})

	t.Run("should have id if connection message follows the correct pattern", func(t *testing.T) {
		const USER_ID = 123
		update := types.TelegramUpdate{
			Message: types.Message{
				Text: fmt.Sprintf("    /connect     %d     ", USER_ID),
			},
		}
		id, isConnection := getConnectionId(&update)
		if isConnection != true {
			t.Errorf("should be true")
		}
		if id != USER_ID {
			t.Errorf("wrong user id, wanted %d, got %d", USER_ID, id)
		}
	})
}

func assertNonErrorReply(t testing.TB, got, want []types.ReplyDTO, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("shouldn't have error")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
