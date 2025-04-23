package bot

import (
	"fmt"
	"github.com/ikebastuz/tgn/types"
)

func createConnectionMessage(userId int64) string {
	return fmt.Sprintf("%s %v", FORWARD_CONNECTION_MESSAGE, userId)
}

func getSenderId(update types.TelegramUpdate) (int64, error) {
	if update.CallbackQuery.From.ID > 0 {
		userId := update.CallbackQuery.From.ID
		return userId, nil
	}
	if update.Message.From.ID > 0 {
		userId := update.Message.From.ID
		return userId, nil
	}

	return 0, ErrorNoSenderIdFound
}

func getDialogState(userId int64, store types.Store) (*types.DialogState, error) {
	dialogState := store.GetDialogState(userId)
	return dialogState, nil
}
