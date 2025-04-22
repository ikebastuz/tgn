package bot

import (
	"errors"

	"github.com/ikebastuz/tgn/types"
)

var (
	ErrorNoSenderIdFound = errors.New("no sender ID found in message")
)

func HandleMessage(update types.TelegramUpdate, store types.Store) {
	createReply(update, store)
}

func createReply(update types.TelegramUpdate, store types.Store) (types.ReplyDTO, error) {
	return types.ReplyDTO{
		Message: update.Message.Text,
	}, nil
}

func GetDialogState(update types.TelegramUpdate, store types.Store) (*types.DialogState, error) {
	userId, err := getSenderId(update)
	if err != nil {
		return nil, ErrorNoSenderIdFound
	}
	dialogState := store.GetDialogState(userId)
	return dialogState, nil
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
