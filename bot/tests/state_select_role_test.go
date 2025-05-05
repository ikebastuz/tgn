package bot_test

import (
	"github.com/ikebastuz/tgn/actions"
	"github.com/ikebastuz/tgn/bot"
	"github.com/ikebastuz/tgn/types"
	"testing"
)

func TestCreateReplySelectRole(t *testing.T) {
	t.Run("SELECT ROLE state - selected EMPLOYEE - update both users and ask for lower bounds", func(t *testing.T) {
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID_2,
		})
		store.SetDialogState(&TEST_USER_ID_2, &types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID,
		})

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
						Message:     bot.CreateSelectLowerBoundsMessage(types.ROLE_EMPLOYEE),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: TEST_FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.CreateSelectLowerBoundsMessage(types.ROLE_EMPLOYER),
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

		got, err := bot.CreateReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT ROLE state - selected EMPLOYER - update both users and ask for lower bounds", func(t *testing.T) {
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID_2,
		})
		store.SetDialogState(&TEST_USER_ID_2, &types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID,
		})

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
						Message:     bot.CreateSelectLowerBoundsMessage(types.ROLE_EMPLOYER),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: TEST_FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.CreateSelectLowerBoundsMessage(types.ROLE_EMPLOYEE),
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

		got, err := bot.CreateReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT ROLE state - received unexpected data - prompt for role again", func(t *testing.T) {
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID_2,
		})
		store.SetDialogState(&TEST_USER_ID_2, &types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID,
		})

		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.MESSAGE_SELECT_YOUR_ROLE_UNEXPECTED,
						ReplyMarkup: bot.KEYBOARD_SELECT_YOUR_ROLE,
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

		got, err := bot.CreateReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT ROLE state - received text message - prompt for role again", func(t *testing.T) {
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID_2,
		})
		store.SetDialogState(&TEST_USER_ID_2, &types.SelectYourRoleState{
			OpponentId: &TEST_USER_ID,
		})

		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.MESSAGE_SELECT_YOUR_ROLE_UNEXPECTED,
						ReplyMarkup: bot.KEYBOARD_SELECT_YOUR_ROLE,
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

		got, err := bot.CreateReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}
		assertNonErrorReply(t, got, want, err)
	})

}
