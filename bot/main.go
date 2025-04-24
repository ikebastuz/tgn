package bot

import (
	"context"
	"log"
	"math/rand"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/ikebastuz/tgn/types"
)

func HandleMessage(ctx context.Context, client *telegram.Client, update types.TelegramUpdate, store types.Store) error {
	replies, err := createReply(update, store)
	if err != nil {
		return err
	}

	userData, err := getUserData(update)
	if err != nil {
		return err
	}

	for _, reply := range replies {
		message := reply.Message
		nextState := reply.NextState

		if message.MessageID > 0 {
			_, err = client.API().MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
				ID:          message.MessageID,
				Peer:        &tg.InputPeerUser{UserID: message.UserID},
				Message:     message.Message,
				ReplyMarkup: message.ReplyMarkup,
			})

			if err != nil {
				log.Printf("ERROR: Failed to edit message: %v", err)
				return err
			}
		} else {
			_, err = client.API().MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
				RandomID:    rand.Int63(),
				Peer:        &tg.InputPeerUser{UserID: message.UserID},
				Message:     message.Message,
				ReplyMarkup: message.ReplyMarkup,
			})
			if err != nil {
				log.Printf("ERROR: Failed to send message: %v", err)
				return err
			}
		}

		if nextState != nil {
			store.SetDialogState(userData.ID, nextState)
		}
	}

	return nil
}

func createReply(update types.TelegramUpdate, store types.Store) ([]types.ReplyDTO, error) {
	userData, err := getUserData(update)
	if err != nil {
		return nil, err
	}

	dialogState, err := getDialogState(userData.ID, store)
	if err != nil {
		return nil, err
	}

	switch dialogState.State {
	case types.STATE_INITIAL:
		connectionId, isConnectionMessage := getConnectionId(&update)

		if isConnectionMessage {
			if connectionId == userData.ID {
				return []types.ReplyDTO{
					{
						Message: types.ReplyMessage{
							UserID:      userData.ID,
							Message:     MESSAGE_YOU_CANT_CONNECT_TO_YOURSELF,
							ReplyMarkup: nil,
						},
					},
				}, nil
			}
			if store.GetDialogState(connectionId).State != types.WAITING_FOR_CONNECT {
				return []types.ReplyDTO{
					{
						Message: types.ReplyMessage{
							UserID:      userData.ID,
							Message:     MESSAGE_NO_SUCH_USER_IS_AWATING,
							ReplyMarkup: nil,
						},
					},
				}, nil
			}
			// TODO: handle connection here
			return []types.ReplyDTO{}, nil
		}
		return []types.ReplyDTO{
			{
				Message: types.ReplyMessage{
					UserID:      userData.ID,
					Message:     createConnectionMessage(*userData),
					ReplyMarkup: nil,
				},
				NextState: &types.DialogState{
					State: types.WAITING_FOR_CONNECT,
				},
			},
			{
				Message: types.ReplyMessage{
					UserID:      userData.ID,
					Message:     MESSAGE_FORWARD_CONNECTION_02,
					ReplyMarkup: nil,
				},
			},
		}, nil
	case types.WAITING_FOR_CONNECT:
		return []types.ReplyDTO{
			{
				Message: types.ReplyMessage{
					UserID:      userData.ID,
					Message:     MESSAGE_WAITING_FOR_CONNECTION,
					ReplyMarkup: nil,
				},
			},
		}, nil

	default:
		return nil, nil
	}
}
