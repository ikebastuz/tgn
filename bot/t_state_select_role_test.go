package bot

import (
	"github.com/ikebastuz/tgn/actions"
	"github.com/ikebastuz/tgn/types"
	"testing"
)

func TestCreateReplySelectRole(t *testing.T) {
	t.Run("SELECT ROLE state - selected EMPLOYEE - update both users and ask for lower bounds", func(t *testing.T) {
		store := NewInMemoryStore()
		sm1 := types.StateMachine{}
		sm1.SetState(&types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID_2,
		})
		store.states[TEST_USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID,
		})
		store.states[TEST_USER_ID_2] = &sm2

		var nextState1 types.State = &types.SelectLowerBoundsState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
		}
		var nextState2 types.State = &types.SelectLowerBoundsState{
			OpponentId: &TEST_USER_ID,
			Role:       types.ROLE_EMPLOYER,
		}
		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     createSelectLowerBoundsMessage(types.ROLE_EMPLOYEE),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: TEST_FROM_2.ID,
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
				From:      TEST_FROM,
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
			OpponentId: &TEST_USER_ID_2,
		})
		store.states[TEST_USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID,
		})
		store.states[TEST_USER_ID_2] = &sm2

		var nextState1 types.State = &types.SelectLowerBoundsState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYER,
		}
		var nextState2 types.State = &types.SelectLowerBoundsState{
			OpponentId: &TEST_USER_ID,
			Role:       types.ROLE_EMPLOYEE,
		}
		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     createSelectLowerBoundsMessage(types.ROLE_EMPLOYER),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: TEST_FROM_2.ID,
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
				From:      TEST_FROM,
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
			OpponentId: &TEST_USER_ID_2,
		})
		store.states[TEST_USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID,
		})
		store.states[TEST_USER_ID_2] = &sm2

		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
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
				From:      TEST_FROM,
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
			OpponentId: &TEST_USER_ID_2,
		})
		store.states[TEST_USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID,
		})
		store.states[TEST_USER_ID_2] = &sm2

		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
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
				From:      TEST_FROM,
			},
		}

		got, err := createReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}
		assertNonErrorReply(t, got, want, err)
	})

}
