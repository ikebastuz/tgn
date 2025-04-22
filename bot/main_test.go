package bot

import (
	"testing"

	"github.com/ikebastuz/tgn/types"
)

func TestGetDialogState(t *testing.T) {
	t.Run("should create initial dialog state", func(t *testing.T) {
		store := NewInMemoryStore()
		update := types.TelegramUpdate{
			Message: types.Message{
				From: types.From{
					ID: 123,
				},
			},
		}
		want := types.State(types.STATE_INITIAL)
		got, _ := GetDialogState(update, store)
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

// func TestCreateReply(t *testing.T) {
// 	t.Run("should return telegram ID if user is not in store", func(t *testing.T) {
// 		store := Store{
// 			count: 0,
// 		}
// 		want := "hello"
// 		update := types.TelegramUpdate{
// 			UpdateID: 1,
// 			Message: types.Message{
// 				MessageID: 1,
// 				Text:      want,
// 				From: types.From{
// 					USERNAME: "",
// 					IS_BOT:   false,
// 					ID:       1,
// 				},
// 			},
// 		}
//
// 		reply, err := createReply(update, &store)
// 		if err != nil {
// 			t.Errorf("shouldn't have error")
// 		}
//
// 		if reply.Message != want {
// 			t.Errorf("expected %v, got %v", want, reply.Message)
// 		}
// 	})
// }

func TestGetSenderId(t *testing.T) {
	const USER_ID = int64(123)

	t.Run("should return correct sender id from Message", func(t *testing.T) {
		update := types.TelegramUpdate{
			Message: types.Message{
				From: types.From{
					ID: USER_ID,
				},
			},
		}
		got, _ := getSenderId(update)
		if got != USER_ID {
			t.Errorf("expected state %v, got %v", USER_ID, got)
		}
	})

	t.Run("should return correct sender id from Callback", func(t *testing.T) {
		update := types.TelegramUpdate{
			CallbackQuery: types.CallbackQuery{
				From: types.From{
					ID: USER_ID,
				},
			},
		}
		got, _ := getSenderId(update)
		if got != USER_ID {
			t.Errorf("expected state %v, got %v", USER_ID, got)
		}
	})

	t.Run("should throw an error if no sender id exists", func(t *testing.T) {
		update := types.TelegramUpdate{}
		_, err := getSenderId(update)
		if err == ErrorNoSenderIdFound {
			t.Errorf("expected error")
		}
	})
}
