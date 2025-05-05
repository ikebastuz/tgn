package bot_test

import (
	"testing"

	"github.com/ikebastuz/tgn/bot"
	"github.com/ikebastuz/tgn/types"
)

func TestCreateReplyEdge(t *testing.T) {
	t.Run("ANY state with /reset - should set to initial state", func(t *testing.T) {
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.SelectLowerBoundsState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
		})
		store.SetDialogState(&TEST_USER_ID_2, &types.SelectLowerBoundsState{
			OpponentId: &TEST_USER_ID,
			Role:       types.ROLE_EMPLOYER,
		})

		want := []types.ReplyDTO{
			{
				UserId: TEST_USER_ID_2,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
			},
			{
				UserId: TEST_USER_ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.MESSAGE_START_GUIDE,
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
				From:      TEST_FROM,
			},
		}

		got, err := bot.CreateReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}

		assertNonErrorReply(t, got, want, err)

		wantState := &types.StateMachine{}
		wantState.SetState(&types.InitialState{})

		gotUserState := store.GetDialogState(&TEST_USER_ID)
		gotUser2State := store.GetDialogState(&TEST_USER_ID_2)

		assertState(t, *gotUserState, *wantState)
		assertState(t, *gotUser2State, *wantState)
	})

	t.Run("UNEXPECTED state - should suggest the user to reset the state", func(t *testing.T) {
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.UnexpectedState{})

		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.MESSAGE_UNEXPECTED_STATE,
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
