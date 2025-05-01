package bot

import (
	"testing"

	"github.com/ikebastuz/tgn/types"
)

func TestCreateReplyEdge(t *testing.T) {
	t.Run("ANY state with /reset - should set to initial state", func(t *testing.T) {
		store := NewInMemoryStore()
		sm1 := types.StateMachine{}
		sm1.SetState(&types.SelectLowerBoundsState{
			OpponentId: &TEST_USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
		})
		store.states[TEST_USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectLowerBoundsState{
			OpponentId: &TEST_USER_ID,
			Role:       types.ROLE_EMPLOYER,
		})
		store.states[TEST_USER_ID_2] = &sm2

		want := []types.ReplyDTO{
			{
				UserId: TEST_USER_ID_2,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
			},
			{
				UserId: TEST_USER_ID,
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
				From:      TEST_FROM,
			},
		}

		got, err := createReply(update, store)
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
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.UnexpectedState{})
		store.states[TEST_USER_ID] = &sm

		want := []types.ReplyDTO{
			{
				UserId: TEST_FROM.ID,
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
