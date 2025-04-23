package bot

import (
	"context"
	"errors"
	"log"
	"math/rand"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/ikebastuz/tgn/types"
)

var (
	ErrorNoSenderIdFound = errors.New("no sender ID found in message")
)

const (
	FORWARD_CONNECTION_MESSAGE = "Forward this message to the person\nyou want to negotiate with\n\nTo join - type\n/connect"
)

func HandleMessage(ctx context.Context, client *telegram.Client, update types.TelegramUpdate, store types.Store) error {
	replies, err := createReply(update, store)

	if err != nil {
		return err
	}

	for _, reply := range replies {
		if reply.MessageID > 0 {
			_, err = client.API().MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
				ID:          update.CallbackQuery.Message.MessageID,
				Peer:        &tg.InputPeerUser{UserID: reply.UserID},
				Message:     reply.Message,
				ReplyMarkup: reply.ReplyMarkup,
			})

			if err != nil {
				log.Printf("ERROR: Failed to edit message: %v", err)
				return err
			}
			// TODO: advance state
		} else {
			_, err = client.API().MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
				RandomID:    rand.Int63(),
				Peer:        &tg.InputPeerUser{UserID: reply.UserID},
				Message:     reply.Message,
				ReplyMarkup: reply.ReplyMarkup,
			})
			if err != nil {
				log.Printf("ERROR: Failed to send message: %v", err)
				return err
			}
			// TODO: advance state
		}
	}

	return nil
}

func createReply(update types.TelegramUpdate, store types.Store) ([]types.ReplyDTO, error) {
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
		return []types.ReplyDTO{
			{
				UserID:      userId,
				Message:     createConnectionMessage(userId),
				ReplyMarkup: nil,
			},
		}, nil
	default:
		return nil, nil
	}
}
