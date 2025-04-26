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

	for _, reply := range replies {
		for _, message := range reply.Messages {
			if message.MessageID > 0 {
				_, err = client.API().MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
					ID:          message.MessageID,
					Peer:        &tg.InputPeerUser{UserID: reply.UserId},
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
					Peer:        &tg.InputPeerUser{UserID: reply.UserId},
					Message:     message.Message,
					ReplyMarkup: message.ReplyMarkup,
				})
				if err != nil {
					log.Printf("ERROR: Failed to send message: %v", err)
					return err
				}
			}
		}

		if reply.NextState != nil {
			store.SetDialogState(&reply.UserId, reply.NextState)
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

	if isResetMessage(&update) {
		store.ResetUserState(&userData.ID)

		return []types.ReplyDTO{
			{
				UserId: userData.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
			},
		}, nil
	}

	switch dialogState.State {
	case types.STATE_INITIAL:
		incomingConnectionId, isConnectionMessage := getConnectionId(&update)

		if isStartMessage(&update) {
			// Start message
			// TODO: check if not connection exists
			newConnectionId := store.CreateConnectionId(&userData.ID)
			return []types.ReplyDTO{
				{
					UserId: userData.ID,
					Messages: []types.ReplyMessage{
						{
							Message:     createConnectionMessage(userData.USERNAME, newConnectionId),
							ReplyMarkup: nil,
						},
					},
					NextState: &types.DialogState{
						State:        types.WAITING_FOR_CONNECT,
						ConnectionId: &newConnectionId,
					},
				},
				{
					UserId: userData.ID,
					Messages: []types.ReplyMessage{
						{
							Message:     MESSAGE_FORWARD_CONNECTION_02,
							ReplyMarkup: nil,
						},
					},
				},
			}, nil
		} else if isConnectionMessage {
			// Connect message
			// TODO: handle connection
			targetUserId := store.GetConnectionTarget(&incomingConnectionId)
			if targetUserId == nil {
				return []types.ReplyDTO{
					{
						UserId: userData.ID,
						Messages: []types.ReplyMessage{
							{
								Message:     MESSAGE_NO_SUCH_USER_IS_AWATING,
								ReplyMarkup: nil,
							},
						},
					},
				}, nil
			} else if *targetUserId == userData.ID {
				return []types.ReplyDTO{
					{
						UserId: userData.ID,
						Messages: []types.ReplyMessage{
							{
								Message:     MESSAGE_YOU_CANT_CONNECT_TO_YOURSELF,
								ReplyMarkup: nil,
							},
						},
					},
				}, nil
			} else {
				// TODO: update store
				store.DeleteConnectionId(&incomingConnectionId)
				return []types.ReplyDTO{
					{
						UserId: userData.ID,
						Messages: []types.ReplyMessage{
							{
								Message:     MESSAGE_SELECT_YOUR_ROLE,
								ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
							},
						},
						NextState: &types.DialogState{
							State:      types.SELECT_YOUR_ROLE,
							OpponentId: targetUserId,
						},
					},
					{
						UserId: *targetUserId,
						Messages: []types.ReplyMessage{
							{
								Message:     MESSAGE_SELECT_YOUR_ROLE,
								ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
							},
						},
						NextState: &types.DialogState{
							State:      types.SELECT_YOUR_ROLE,
							OpponentId: &userData.ID,
						},
					},
				}, nil
			}
		} else {
			// Irrelevant - show guide
			return []types.ReplyDTO{
				{
					UserId: userData.ID,
					Messages: []types.ReplyMessage{
						{
							Message:     MESSAGE_START_GUIDE,
							ReplyMarkup: nil,
						},
					},
				},
			}, nil
		}

	case types.WAITING_FOR_CONNECT:
		return []types.ReplyDTO{
			{
				UserId: userData.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_WAITING_FOR_CONNECTION,
						ReplyMarkup: nil,
					},
				},
			},
		}, nil

	default:
		// TODO: handle unexpected state
		return []types.ReplyDTO{
			{
				UserId: userData.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_UNEXPECTED_STATE,
						ReplyMarkup: nil,
					},
				},
			},
		}, nil
	}
}
