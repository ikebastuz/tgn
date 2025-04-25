package bot

import (
	"reflect"
	"testing"

	"github.com/ikebastuz/tgn/types"
)

func TestCreateReply(t *testing.T) {
	var USER_ID int64 = 123
	var FROM = types.From{
		ID:       int64(USER_ID),
		USERNAME: "hello",
	}
	var CONNECTION_ID int64 = 100500

	// t.Run("INITIAL state, irrelevant message", func(t *testing.T) {
	// 	store := NewInMemoryStore()
	//
	// 	want := []types.ReplyDTO{
	// 		{
	// 			Message: types.ReplyMessage{
	// 				UserID:      FROM.ID,
	// 				Message:     MESSAGE_START_GUIDE,
	// 				ReplyMarkup: nil,
	// 			},
	// 		},
	// 	}
	//
	// 	update := types.TelegramUpdate{
	// 		Message: types.Message{
	// 			From: FROM,
	// 		},
	// 	}
	//
	// 	got, err := createReply(update, store)
	// 	assertNonErrorReply(t, got, want, err)
	// })

	t.Run("state - INITIAL and it is not a 'connect' message - should ask user to forward connection message", func(t *testing.T) {
		store := NewInMemoryStore()
		store.states[USER_ID] = types.DialogState{
			State:        types.STATE_INITIAL,
			ConnectionId: &CONNECTION_ID,
		}

		want := []types.ReplyDTO{
			{
				Message: types.ReplyMessage{
					UserID:      FROM.ID,
					Message:     createConnectionMessage(FROM.USERNAME, CONNECTION_ID),
					ReplyMarkup: nil,
				},
				NextState: &types.DialogState{
					State:        types.WAITING_FOR_CONNECT,
					ConnectionId: &CONNECTION_ID,
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
	// t.Run("state - INITIAL and is a 'connect' message to yourself - should show an error", func(t *testing.T) {
	// 	store := NewInMemoryStore()
	// 	want := []types.ReplyDTO{
	// 		{
	// 			Message: types.ReplyMessage{
	// 				UserID:      FROM.ID,
	// 				Message:     MESSAGE_YOU_CANT_CONNECT_TO_YOURSELF,
	// 				ReplyMarkup: nil,
	// 			},
	// 		},
	// 	}
	//
	// 	update := types.TelegramUpdate{
	// 		Message: types.Message{
	// 			Text: fmt.Sprintf("/connect %d", USER_ID),
	// 			From: FROM,
	// 		},
	// 	}
	//
	// 	got, err := createReply(update, store)
	// 	assertNonErrorReply(t, got, want, err)
	// })
	// t.Run("state - INITIAL and is a 'connect' message to non-existent user - should show an error", func(t *testing.T) {
	// 	store := NewInMemoryStore()
	// 	want := []types.ReplyDTO{
	// 		{
	// 			Message: types.ReplyMessage{
	// 				UserID:      FROM.ID,
	// 				Message:     MESSAGE_NO_SUCH_USER_IS_AWATING,
	// 				ReplyMarkup: nil,
	// 			},
	// 		},
	// 	}
	//
	// 	update := types.TelegramUpdate{
	// 		Message: types.Message{
	// 			Text: "/connect 100500",
	// 			From: FROM,
	// 		},
	// 	}
	//
	// 	got, err := createReply(update, store)
	// 	assertNonErrorReply(t, got, want, err)
	// })
	t.Run("state - WAITING for connection - should tell about waiting for connection", func(t *testing.T) {
		store := NewInMemoryStore()
		store.SetDialogState(&USER_ID, &types.DialogState{State: types.WAITING_FOR_CONNECT})
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

func assertNonErrorReply(t testing.TB, got, want []types.ReplyDTO, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("shouldn't have error")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected \n%v, \ngot \n%v", want, got)
	}
}
