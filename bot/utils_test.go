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

func TestIsStartMessage(t *testing.T) {
	t.Run("is not a start message", func(t *testing.T) {
		update := types.TelegramUpdate{}
		got := isStartMessage(&update)
		if got != false {
			t.Errorf("should be false")
		}
	})
	t.Run("is a start message", func(t *testing.T) {
		update := types.TelegramUpdate{
			Message: types.Message{
				Text: " /start ",
			},
		}
		got := isStartMessage(&update)
		if got == false {
			t.Errorf("should be true")
		}
	})
}

func TestIsResetMessage(t *testing.T) {
	t.Run("is not a reset message", func(t *testing.T) {
		update := types.TelegramUpdate{}
		got := isResetMessage(&update)
		if got != false {
			t.Errorf("should be false")
		}
	})
	t.Run("is a reset message", func(t *testing.T) {
		update := types.TelegramUpdate{
			Message: types.Message{
				Text: " /reset ",
			},
		}
		got := isResetMessage(&update)
		if got == false {
			t.Errorf("should be true")
		}
	})
}
