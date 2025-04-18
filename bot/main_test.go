package bot

import (
	"testing"

	"github.com/ikebastuz/tgn/types"
)

func TestHandleMessage(t *testing.T) {
	t.Run("should update store by reference", func(t *testing.T) {
		store := Store{
			count: 0,
		}
		update := types.TelegramUpdate{}
		want := int64(1)
		HandleMessage(update, &store)
		if store.count != want {
			t.Errorf("expected %v, got %v", want, store.count)
		}
	})
}

func TestCreateReply(t *testing.T) {
	t.Run("should return telegram ID if user is not in store", func(t *testing.T) {
		store := Store{
			count: 0,
		}
		want := "hello"
		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      want,
				From: types.From{
					USERNAME: "",
					IS_BOT:   false,
					ID:       1,
				},
			},
		}

		reply, err := createReply(update, &store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}

		if reply.Message != want {
			t.Errorf("expected %v, got %v", want, reply.Message)
		}
	})
}
