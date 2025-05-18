package bot_test

import (
	"testing"

	"github.com/ikebastuz/tgn/bot"
	"github.com/ikebastuz/tgn/bot/actions"
	"github.com/ikebastuz/tgn/types"
)

func TestCreateReplyResult(t *testing.T) {
	t.Run("RESULT Error state - Selected No - Set both to initial state", func(t *testing.T) {
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.ResultErrorState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
		})
		store.SetDialogState(&TEST_USER_ID_2, &types.ResultErrorState{
			OpponentId: &TEST_USER_ID,
			Role:       types.ROLE_EMPLOYEE,
		})

		var nextState types.State = &types.InitialState{}

		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState,
			},
			{
				UserId: TEST_USER_ID_2,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState,
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      "anything",
				From:      TEST_FROM,
			},
			CallbackQuery: types.CallbackQuery{
				Data: actions.ACTION_SELECT_NO,
			},
		}

		got, err := bot.CreateReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("RESULT Error state - Selected Yes - Move to select lower bounds state", func(t *testing.T) {
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.ResultErrorState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
		})
		store.SetDialogState(&TEST_USER_ID_2, &types.ResultErrorState{
			OpponentId: &TEST_USER_ID,
			Role:       types.ROLE_EMPLOYER,
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
				UserId: TEST_USER_ID_2,
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
			Message: types.Message{
				MessageID: 1,
				Text:      "anything",
				From:      TEST_FROM,
			},
			CallbackQuery: types.CallbackQuery{
				Data: actions.ACTION_SELECT_YES,
			},
		}

		got, err := bot.CreateReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})
}
