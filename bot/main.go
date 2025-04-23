package bot

import (
	"errors"

	"github.com/ikebastuz/tgn/types"
)

var (
	ErrorNoSenderIdFound = errors.New("no sender ID found in message")
)

const (
	FORWARD_CONNECTION_MESSAGE = "Forward this message\n/connect"
)

func HandleMessage(update types.TelegramUpdate, store types.Store) {
	createReply(update, store)
}

func createReply(update types.TelegramUpdate, store types.Store) (*types.ReplyDTO, error) {
	userId, err := getSenderId(update)
	if err != nil {
		return nil, err
	}
	dialogState, err := getDialogState(userId, store)

	if err != nil {
		return nil, err
	}

	switch dialogState.State {
	case types.STATE_INITIAL:
		return &types.ReplyDTO{
			UserID:      userId,
			Message:     createConnectionMessage(userId),
			ReplyMarkup: nil,
		}, nil
	default:
		return &types.ReplyDTO{}, nil
	}
}

func getDialogState(userId int64, store types.Store) (*types.DialogState, error) {
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
