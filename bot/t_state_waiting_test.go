package bot

import (
	"github.com/ikebastuz/tgn/types"
	"testing"
)

func TestCreateReplyWaiting(t *testing.T) {
	t.Run("WAITING state, - tells about waiting for connection", func(t *testing.T) {
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.WaitingForConnectState{
			ConnectionId: &TEST_CONNECTION_ID,
		})
		store.states[TEST_FROM.ID] = &sm

		var FROM = types.From{
			ID:       int64(TEST_USER_ID),
			USERNAME: "hello",
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_WAITING_FOR_CONNECTION,
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
