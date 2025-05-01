package bot

import (
	"github.com/ikebastuz/tgn/actions"
	"github.com/ikebastuz/tgn/types"
	"testing"
)

func TestCreateReplySelectRole(t *testing.T) {
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

	t.Run("SELECT ROLE state - selected EMPLOYEE - update both users and ask for lower bounds", func(t *testing.T) {
		store := NewInMemoryStore()
		sm1 := types.StateMachine{}
		sm1.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID_2,
		})
		store.states[USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID,
		})
		store.states[USER_ID_2] = &sm2

		var nextState1 types.State = &types.SelectLowerBoundsState{
			OpponentId: &USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
		}
		var nextState2 types.State = &types.SelectLowerBoundsState{
			OpponentId: &USER_ID,
			Role:       types.ROLE_EMPLOYER,
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     createSelectLowerBoundsMessage(types.ROLE_EMPLOYEE),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     createSelectLowerBoundsMessage(types.ROLE_EMPLOYER),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState2,
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
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT ROLE state - selected EMPLOYER - update both users and ask for lower bounds", func(t *testing.T) {
		store := NewInMemoryStore()
		sm1 := types.StateMachine{}
		sm1.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID_2,
		})
		store.states[USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID,
		})
		store.states[USER_ID_2] = &sm2

		var nextState1 types.State = &types.SelectLowerBoundsState{
			OpponentId: &USER_ID_2,
			Role:       types.ROLE_EMPLOYER,
		}
		var nextState2 types.State = &types.SelectLowerBoundsState{
			OpponentId: &USER_ID,
			Role:       types.ROLE_EMPLOYEE,
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     createSelectLowerBoundsMessage(types.ROLE_EMPLOYER),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     createSelectLowerBoundsMessage(types.ROLE_EMPLOYEE),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState2,
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			CallbackQuery: types.CallbackQuery{
				Data: actions.ACTION_SELECT_EMPLOYER,
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
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT ROLE state - received unexpected data - prompt for role again", func(t *testing.T) {
		store := NewInMemoryStore()
		sm1 := types.StateMachine{}
		sm1.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID_2,
		})
		store.states[USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID,
		})
		store.states[USER_ID_2] = &sm2

		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_YOUR_ROLE_UNEXPECTED,
						ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			CallbackQuery: types.CallbackQuery{
				Data: " unexpected value ",
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
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT ROLE state - received text message - prompt for role again", func(t *testing.T) {
		store := NewInMemoryStore()
		sm1 := types.StateMachine{}
		sm1.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID_2,
		})
		store.states[USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID,
		})
		store.states[USER_ID_2] = &sm2

		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_YOUR_ROLE_UNEXPECTED,
						ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      " some text ",
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}
		assertNonErrorReply(t, got, want, err)
	})

}
