package bot

import (
	"github.com/ikebastuz/tgn/types"
)

func HandleMessage(update types.TelegramUpdate, store *Store) {
	store.mu.Lock()
	store.count++
	store.mu.Unlock()
}
