package bot

import (
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

	// t.Run("should update store by reference", func(t *testing.T) {
	// 	store := Store{
	// 		count: 0,
	// 	}
	// 	update := types.TelegramUpdate{}
	// 	want := int64(1)
	// 	HandleMessage(update, &store)
	// 	if store.count != want {
	// 		t.Errorf("expected %v, got %v", want, store.count)
	// 	}
	// })
}

func TestCreateReply(t *testing.T) {
	t.Run("should ask user to forward connection message", func(t *testing.T) {
		store := NewInMemoryStore()
		var FROM = types.From{
			ID:       int64(123),
			USERNAME: "hello",
		}
		want := []types.ReplyDTO{
			{
				Message: types.ReplyMessage{
					UserID:      FROM.ID,
					Message:     createConnectionMessage(FROM),
					ReplyMarkup: nil,
				},
			},
			{
				Message: types.ReplyMessage{
					UserID:      FROM.ID,
					Message:     FORWARD_CONNECTION_MESSAGE_02,
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
