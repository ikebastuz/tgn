package bot

import (
	"context"
	log "github.com/sirupsen/logrus"
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
					log.Errorf("Failed to edit message: %v", err)
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
					log.Errorf("Failed to send message: %v", err)
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
		log.Infof("Received RESET message from USER %v", userData.ID)

		// TODO: handle case when already connected to another user
		// need to reset that user as well
		store.ResetUserState(&userData.ID)

		response := []types.ReplyDTO{
			{
				UserId: userData.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
			},
		}
		opponentId := dialogState.OpponentId
		if opponentId != nil {
			log.Infof("Resetting opponent %d state aswell", &opponentId)

			store.ResetUserState(opponentId)
			response = append(response, types.ReplyDTO{
				UserId: *opponentId,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
			})
		}

		return response, nil
	}

	log.Infof("User state is %s", dialogState.State)

	switch dialogState.State {
	case types.STATE_INITIAL:
		incomingConnectionId, isConnectionMessage := getConnectionId(&update)

		if isStartMessage(&update) {
			log.Info("Received START message")

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
			log.Info("Received CONNECT message")
			// Connect message
			// TODO: handle connection
			targetUserId := store.GetConnectionTarget(&incomingConnectionId)
			if targetUserId == nil {
				log.Warnf("Connection id %d does not exist", incomingConnectionId)
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
				log.Warn("Trying to connect to yourself")
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
				log.Infof("Connecting USER %d to USER %d", userData.ID, *targetUserId)
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
			log.Warn("unknown command, showing guide")
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

	case types.SELECT_YOUR_ROLE:
		opponentId := dialogState.OpponentId

		return []types.ReplyDTO{
			{
				UserId: userData.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
						ReplyMarkup: nil,
					},
				},
				NextState: &types.DialogState{
					State: types.SELECT_LOWER_BOUNDS,
				},
			},
			{
				UserId: *opponentId,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
						ReplyMarkup: nil,
					},
				},
				NextState: &types.DialogState{
					State: types.SELECT_LOWER_BOUNDS,
				},
			},
		}, nil

	default:
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
