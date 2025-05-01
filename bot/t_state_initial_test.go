package bot

import (
	"fmt"
	"testing"

	"github.com/ikebastuz/tgn/types"
)

func TestCreateReplyInitial(t *testing.T) {
	t.Run("INITIAL state, irrelevant message - should show guide", func(t *testing.T) {
		store := NewInMemoryStore()

		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
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
				From: TEST_FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("INITIAL state, /start message - create connection to forward", func(t *testing.T) {
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.InitialState{})
		store.states[TEST_USER_ID] = &sm

		var nextState types.State = &types.WaitingForConnectState{
			ConnectionId: &TEST_CONNECTION_ID,
		}
		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     createConnectionMessage(TEST_FROM.USERNAME, TEST_CONNECTION_ID),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState,
			},
			{
				UserId: TEST_FROM.ID,
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
				From: TEST_FROM,
				Text: " /start ",
			},
		}

		got, err := createReply(update, store)

		// TODO: test with correct connection id
		assertSingleReply(t, got[1], want[1], err)
	})

	t.Run("INITIAL state, /connect message to yourself - should show an error", func(t *testing.T) {
		store := NewInMemoryStore()
		store.connections[TEST_CONNECTION_ID] = TEST_USER_ID

		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
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
				Text: fmt.Sprintf("/connect %d", TEST_CONNECTION_ID),
				From: TEST_FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("INITIAL state, /connect message to non-existent user - should show an error", func(t *testing.T) {
		store := NewInMemoryStore()
		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
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
				From: TEST_FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("INITIAL state, /connect message to existing user - should connect correctly", func(t *testing.T) {
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.WaitingForConnectState{
			ConnectionId: &TEST_CONNECTION_ID,
		})
		store.states[TEST_USER_ID] = &sm
		store.connections[TEST_CONNECTION_ID] = TEST_USER_ID

		var nextState1 types.State = &types.SelectYourRoleState{
			OpponentId: &TEST_FROM_2.ID,
		}
		var nextState2 types.State = &types.SelectYourRoleState{
			OpponentId: &TEST_FROM.ID,
		}

		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_YOUR_ROLE_CONNECTED,
						ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
					},
				},
				NextState: &nextState2,
			},
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_YOUR_ROLE_CONNECTED,
						ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
					},
				},
				NextState: &nextState1,
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				Text: fmt.Sprintf("/connect %d", TEST_CONNECTION_ID),
				From: TEST_FROM_2,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})
}
