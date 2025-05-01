package bot

import (
	"github.com/ikebastuz/tgn/types"
	"testing"
)

func TestCreateReplyWaiting(t *testing.T) {
	var USER_ID int64 = 123
	var FROM = types.From{
		ID:       int64(USER_ID),
		USERNAME: "hello",
	}
	var CONNECTION_ID int64 = 100500

	t.Run("WAITING state, - tells about waiting for connection", func(t *testing.T) {
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.WaitingForConnectState{
			ConnectionId: &CONNECTION_ID,
		})
		store.states[FROM.ID] = &sm

		var FROM = types.From{
			ID:       int64(USER_ID),
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
