package bot_test

import (
	"fmt"
	"testing"

	"github.com/ikebastuz/tgn/bot"
	"github.com/ikebastuz/tgn/types"
)

func TestCreateReplySelectUP(t *testing.T) {
	t.Run("SELECT UPPER BOUNDS state - show error message on invalid value", func(t *testing.T) {
		var lower_bound int64 = 100500
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.SelectUpperBoundsState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
			LowerBound: &lower_bound,
		})

		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.MESSAGE_USE_VALID_POSITIVE_NUMBER,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      "qweasd",
				From:      TEST_FROM,
			},
		}

		got, err := bot.CreateReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT UPPER BOUNDS state - Opponent is still selecting - Waiting for result state", func(t *testing.T) {
		var lower_bound int64 = 100500
		var upper_bound int64 = 100500
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.SelectUpperBoundsState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
			LowerBound: &lower_bound,
		})
		store.SetDialogState(&TEST_USER_ID_2, &types.SelectUpperBoundsState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYER,
			LowerBound: &lower_bound,
		})

		var nextState types.State = &types.WaitingForResultState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
			LowerBound: &lower_bound,
			UpperBound: &lower_bound,
		}
		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.MESSAGE_WAITING_FOR_RESULT,
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
				Text:      fmt.Sprintf("%d", upper_bound),
				From:      TEST_FROM,
			},
		}

		got, err := bot.CreateReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT UPPER BOUNDS state - Opponent is ready and there is shared number - show it", func(t *testing.T) {
		var lower_bound int64 = 100500
		var upper_bound int64 = 100500
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.SelectUpperBoundsState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
			LowerBound: &lower_bound,
		})
		store.SetDialogState(&TEST_USER_ID_2, &types.WaitingForResultState{
			OpponentId: &TEST_USER_ID,
			Role:       types.ROLE_EMPLOYER,
			LowerBound: &lower_bound,
			UpperBound: &upper_bound,
		})

		var nextState types.State = &types.InitialState{}
		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.CreateResultMessage(upper_bound),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState,
			},
			{
				UserId: TEST_FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.CreateResultMessage(upper_bound),
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
				Text:      fmt.Sprintf("%d", upper_bound),
				From:      TEST_FROM,
			},
		}

		got, err := bot.CreateReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT UPPER BOUNDS state - Opponent is ready and there is no shared number - suggest restarting", func(t *testing.T) {
		var lower_bound1 int64 = 100
		var upper_bound1 int64 = 200
		var lower_bound2 int64 = 10
		var upper_bound2 int64 = 20
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.SelectUpperBoundsState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
			LowerBound: &lower_bound1,
		})
		store.SetDialogState(&TEST_USER_ID_2, &types.WaitingForResultState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYER,
			LowerBound: &lower_bound2,
			UpperBound: &upper_bound2,
		})

		var nextState1 types.State = &types.ResultErrorState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
		}
		var nextState2 types.State = &types.ResultErrorState{
			OpponentId: &TEST_USER_ID,
			Role:       types.ROLE_EMPLOYER,
		}
		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.MESSAGE_RESULT_ERROR,
						ReplyMarkup: bot.KEYBOARD_SELECT_YES_NO,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: TEST_FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.MESSAGE_RESULT_ERROR,
						ReplyMarkup: bot.KEYBOARD_SELECT_YES_NO,
					},
				},
				NextState: &nextState2,
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      fmt.Sprintf("%d", upper_bound1),
				From:      TEST_FROM,
			},
		}

		got, err := bot.CreateReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT UPPER BOUNDS state - upper bound is above 3x then lower bound - ", func(t *testing.T) {
		var lower_bound1 int64 = 100
		var upper_bound1 int64 = 500
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.SelectUpperBoundsState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
			LowerBound: &lower_bound1,
		})

		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.CreateUseValidUpperBoundMessage(),
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      fmt.Sprintf("%d", upper_bound1),
				From:      TEST_FROM,
			},
		}

		got, err := bot.CreateReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})
}
