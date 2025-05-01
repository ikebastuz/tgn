package bot

import (
	"testing"

	"github.com/ikebastuz/tgn/actions"
	"github.com/ikebastuz/tgn/types"
)

func TestCreateReplyResult(t *testing.T) {
	var USER_ID int64 = 123
	var USER_ID_2 int64 = 124
	var FROM = types.From{
		ID:       int64(USER_ID),
		USERNAME: "hello",
	}

	t.Run("RESULT SUCCESS state - Show guide", func(t *testing.T) {
		var lower_bound int64 = 100
		var upper_bound int64 = 200
		store := NewInMemoryStore()

		sm := types.StateMachine{}
		s := &types.ResultSuccessState{
			OpponentId: &USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
			LowerBound: &lower_bound,
			UpperBound: &upper_bound,
			Result:     &upper_bound,
		}
		sm.SetState(s)
		store.states[USER_ID] = &sm

		var nextState types.State = s
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
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
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("RESULT Error state - Selected No - Set both to initial state", func(t *testing.T) {
		store := NewInMemoryStore()

		sm1 := types.StateMachine{}
		sm1.SetState(&types.ResultErrorState{
			OpponentId: &USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
		})
		store.states[USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.ResultErrorState{
			OpponentId: &USER_ID,
			Role:       types.ROLE_EMPLOYEE,
		})
		store.states[USER_ID_2] = &sm2

		var nextState types.State = &types.InitialState{}

		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState,
			},
			{
				UserId: USER_ID_2,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
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
				From:      FROM,
			},
			CallbackQuery: types.CallbackQuery{
				Data: actions.ACTION_SELECT_NO,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("RESULT Error state - Selected Yes - Move to select lower bounds state", func(t *testing.T) {
		store := NewInMemoryStore()

		sm1 := types.StateMachine{}
		sm1.SetState(&types.ResultErrorState{
			OpponentId: &USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
		})
		store.states[USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.ResultErrorState{
			OpponentId: &USER_ID,
			Role:       types.ROLE_EMPLOYER,
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
						Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: USER_ID_2,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
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
				From:      FROM,
			},
			CallbackQuery: types.CallbackQuery{
				Data: actions.ACTION_SELECT_YES,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})
}
