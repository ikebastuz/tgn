package bot

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ikebastuz/tgn/actions"
	"github.com/ikebastuz/tgn/types"
)

func TestCreateReply(t *testing.T) {
	var USER_ID int64 = 123
	var USER_ID_2 int64 = 124
	var FROM = types.From{
		ID:       int64(USER_ID),
		USERNAME: "hello",
	}
	var FROM_2 = types.From{
		ID:       int64(USER_ID_2),
		USERNAME: "hello 2",
	}
	var CONNECTION_ID int64 = 100500

	t.Run("INITIAL state, irrelevant message - should show guide", func(t *testing.T) {
		store := NewInMemoryStore()

		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
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

	t.Run("INITIAL state, /start message - create connection to forward", func(t *testing.T) {
		store := NewInMemoryStore()
		store.states[USER_ID] = types.DialogState{
			State:        types.STATE_INITIAL,
			ConnectionId: &CONNECTION_ID,
		}

		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     createConnectionMessage(FROM.USERNAME, CONNECTION_ID),
						ReplyMarkup: nil,
					},
				},
				NextState: &types.DialogState{
					State:        types.WAITING_FOR_CONNECT,
					ConnectionId: &CONNECTION_ID,
				},
			},
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_FORWARD_CONNECTION_02,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				From: FROM,
				Text: " /start ",
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("INITIAL state, /connect message to yourself - should show an error", func(t *testing.T) {
		store := NewInMemoryStore()
		store.connections[CONNECTION_ID] = USER_ID

		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_YOU_CANT_CONNECT_TO_YOURSELF,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				Text: fmt.Sprintf("/connect %d", CONNECTION_ID),
				From: FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("INITIAL state, /connect message to non-existent user - should show an error", func(t *testing.T) {
		store := NewInMemoryStore()
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_NO_SUCH_USER_IS_AWATING,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				Text: "/connect 1337",
				From: FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("INITIAL state, /connect message to existing user - should connect correctly", func(t *testing.T) {
		store := NewInMemoryStore()
		store.states[USER_ID] = types.DialogState{
			State:        types.STATE_INITIAL,
			ConnectionId: &CONNECTION_ID,
		}
		store.connections[CONNECTION_ID] = USER_ID

		want := []types.ReplyDTO{
			{
				UserId: FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_YOUR_ROLE,
						ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
					},
				},
				NextState: &types.DialogState{
					State:      types.SELECT_YOUR_ROLE,
					OpponentId: &FROM.ID,
				},
			},
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_YOUR_ROLE,
						ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
					},
				},
				NextState: &types.DialogState{
					State:      types.SELECT_YOUR_ROLE,
					OpponentId: &FROM_2.ID,
				},
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				Text: fmt.Sprintf("/connect %d", CONNECTION_ID),
				From: FROM_2,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("WAITING state, - tells about waiting for connection", func(t *testing.T) {
		store := NewInMemoryStore()
		store.SetDialogState(&USER_ID, &types.DialogState{State: types.WAITING_FOR_CONNECT})
		var FROM = types.From{
			ID:       int64(USER_ID),
			USERNAME: "hello",
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_WAITING_FOR_CONNECTION,
						ReplyMarkup: nil,
					},
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

	t.Run("SELECT ROLE state - update both users and ask for lower bounds", func(t *testing.T) {
		store := NewInMemoryStore()
		store.SetDialogState(&USER_ID, &types.DialogState{State: types.SELECT_YOUR_ROLE, OpponentId: &USER_ID_2})
		store.SetDialogState(&USER_ID_2, &types.DialogState{State: types.SELECT_YOUR_ROLE, OpponentId: &USER_ID})
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
						ReplyMarkup: nil,
					},
				},
				NextState: &types.DialogState{
					State: types.SELECT_LOWER_BOUNDS,
				},
			},
			{
				UserId: FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
						ReplyMarkup: nil,
					},
				},
				NextState: &types.DialogState{
					State: types.SELECT_LOWER_BOUNDS,
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			CallbackQuery: types.CallbackQuery{
				Data: actions.ACTION_SELECT_EMPLOYEE,
			},
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

	// TODO: cover case when opponentId != null
	t.Run("ANY state with /reset - should set to initial state", func(t *testing.T) {
		store := NewInMemoryStore()
		store.SetDialogState(&USER_ID, &types.DialogState{
			State:      types.SELECT_LOWER_BOUNDS,
			OpponentId: &USER_ID_2,
		})
		store.SetDialogState(&USER_ID_2, &types.DialogState{
			State:      types.SELECT_LOWER_BOUNDS,
			OpponentId: &USER_ID,
		})
		var FROM = types.From{
			ID:       int64(USER_ID),
			USERNAME: "hello",
		}
		wantReply := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
			},
			{
				UserId: FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      " /reset ",
				From:      FROM,
			},
		}

		gotReply, err := createReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}

		if !reflect.DeepEqual(gotReply, wantReply) {
			t.Errorf("expected %v, got %v", wantReply, gotReply)
		}

	})

	t.Run("UNEXPECTED state - should suggest the user to reset the state", func(t *testing.T) {
		store := NewInMemoryStore()
		store.SetDialogState(&USER_ID, &types.DialogState{State: 100500})
		var FROM = types.From{
			ID:       int64(USER_ID),
			USERNAME: "hello",
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_UNEXPECTED_STATE,
						ReplyMarkup: nil,
					},
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
