package bot

import (
	"github.com/ikebastuz/tgn/types"
)

func HandleMessage(update types.TelegramUpdate, store *Store) {
	store.mu.Lock()
	store.count++
	store.mu.Unlock()

	createReply(update, store)
}

func createReply(update types.TelegramUpdate, store *Store) (types.ReplyDTO, error) {
	return types.ReplyDTO{
		Message: update.Message.Text,
	}, nil
}
