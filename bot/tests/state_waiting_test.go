package bot_test

import (
	"github.com/ikebastuz/tgn/bot"
	"github.com/ikebastuz/tgn/types"
	"testing"
)

func TestCreateReplyWaiting(t *testing.T) {
	t.Run("WAITING state, - tells about waiting for connection", func(t *testing.T) {
		store := bot.NewInMemoryStore()
		store.SetDialogState(&TEST_USER_ID, &types.WaitingForConnectState{
			ConnectionId: &TEST_CONNECTION_ID,
		})

		var FROM = types.From{
			ID:       int64(TEST_USER_ID),
			USERNAME: "hello",
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     bot.MESSAGE_WAITING_FOR_CONNECTION,
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

		got, err := bot.CreateReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}

		assertNonErrorReply(t, got, want, err)
	})
}
