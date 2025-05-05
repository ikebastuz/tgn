package bot_test

import (
	"fmt"
	"github.com/ikebastuz/tgn/bot"
	"github.com/ikebastuz/tgn/types"
	"testing"
)

func TestCreateReplySelectLB(t *testing.T) {
	t.Run("SELECT LOWER BOUNDS state - show error message on invalid value", func(t *testing.T) {
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.SelectLowerBoundsState{
			OpponentId: &TEST_FROM_2.ID,
			Role:       types.ROLE_EMPLOYEE,
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

	t.Run("SELECT LOWER BOUNDS state - proceed further to upper bounds state if number is correct", func(t *testing.T) {
		var lower_bound int64 = 100500
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.SelectLowerBoundsState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
		})

		var nextState types.State = &types.SelectUpperBoundsState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
			LowerBound: &lower_bound,
		}
		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.CreateSelectUpperBoundMessage(types.ROLE_EMPLOYEE),
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
				Text:      fmt.Sprintf("%d", lower_bound),
				From:      TEST_FROM,
			},
		}

		got, err := bot.CreateReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})
}
