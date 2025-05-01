package bot

import (
	"fmt"
	"github.com/ikebastuz/tgn/types"
	"testing"
)

func TestCreateReplySelectLB(t *testing.T) {
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

	t.Run("SELECT LOWER BOUNDS state - show error message on invalid value", func(t *testing.T) {
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.SelectLowerBoundsState{
			OpponentId: &FROM_2.ID,
			Role:       types.ROLE_EMPLOYEE,
		})
		store.states[USER_ID] = &sm
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_USE_VALID_POSITIVE_NUMBER,
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
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT LOWER BOUNDS state - proceed further to upper bounds state if number is correct", func(t *testing.T) {
		var lower_bound int64 = 100500
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.SelectLowerBoundsState{
			OpponentId: &USER_ID_2,
			Role:       types.ROLE_EMPLOYEE, // TODO: make it required
		})
		store.states[USER_ID] = &sm

		var nextState types.State = &types.SelectUpperBoundsState{
			OpponentId: &USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
			LowerBound: &lower_bound,
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_UPPER_BOUND,
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
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})
}
