package bot

import (
	"fmt"
	"github.com/ikebastuz/tgn/types"
)

func createConnectionMessage(userData types.From) string {
	return fmt.Sprintf(FORWARD_CONNECTION_MESSAGE_01, userData.USERNAME, userData.ID)
}

func getUserData(update types.TelegramUpdate) (*types.From, error) {
	if update.CallbackQuery.From.ID > 0 {
		return &update.CallbackQuery.From, nil
	}
	if update.Message.From.ID > 0 {
		return &update.Message.From, nil
	}

	return nil, ErrorNoSenderIdFound
}

func getDialogState(userId int64, store types.Store) (*types.DialogState, error) {
	dialogState := store.GetDialogState(userId)
	return dialogState, nil
}
