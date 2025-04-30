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
			store.SetDialogState(&reply.UserId, *reply.NextState)
		}
	}

	return nil
}

func createReply(update types.TelegramUpdate, store types.Store) ([]types.ReplyDTO, error) {
	userData, err := getUserData(update)
	if err != nil {
		return nil, err
	}

	sm, err := getDialogState(userData.ID, store)
	if err != nil {
		return nil, err
	}

	if isResetMessage(&update) {
		log.Infof("Received RESET message from USER %v", userData.ID)

		// TODO: handle case when already connected to another user
		// need to reset that user as well
		log.Infof("Resetting user %d state aswell", userData.ID)
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
		switch s := sm.GetState().(type) {
		case *types.SelectUpperBoundsState:
		case *types.SelectLowerBoundsState:
		case *types.SelectYourRoleState:
			log.Infof("Resetting opponent %d state aswell", *s.OpponentId)

			store.ResetUserState(s.OpponentId)
			response = append(response, types.ReplyDTO{
				UserId: *s.OpponentId,
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

	switch s := sm.GetState().(type) {
	case *types.InitialState:
		incomingConnectionId, isConnectionMessage := getConnectionId(&update)

		if isStartMessage(&update) {
			log.Info("Received START message")

			// Start message
			// TODO: check if not connection exists
			newConnectionId := store.CreateConnectionId(&userData.ID)

			var nextState types.State_NG = &types.WaitingForConnectState{
				ConnectionId: &newConnectionId,
			}

			return []types.ReplyDTO{
				{
					UserId: userData.ID,
					Messages: []types.ReplyMessage{
						{
							Message:     createConnectionMessage(userData.USERNAME, newConnectionId),
							ReplyMarkup: nil,
						},
					},
					NextState: &nextState,
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
				err := store.DeleteConnectionId(&incomingConnectionId)
				if err != nil {
					log.Warnf("No connection %d to delete", incomingConnectionId)
				}

				var nextState1 types.State_NG = &types.SelectYourRoleState{
					OpponentId: targetUserId,
				}
				var nextState2 types.State_NG = &types.SelectYourRoleState{
					OpponentId: &userData.ID,
				}
				return []types.ReplyDTO{
					{
						UserId: userData.ID,
						Messages: []types.ReplyMessage{
							{
								Message:     MESSAGE_SELECT_YOUR_ROLE,
								ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
							},
						},
						NextState: &nextState1,
					},
					{
						UserId: *targetUserId,
						Messages: []types.ReplyMessage{
							{
								Message:     MESSAGE_SELECT_YOUR_ROLE,
								ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
							},
						},
						NextState: &nextState2,
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

	case *types.WaitingForConnectState:
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

	case *types.SelectYourRoleState:
		opponentId := s.OpponentId

		var nextState1 types.State_NG = &types.SelectLowerBoundsState{
			OpponentId: opponentId,
		}
		var nextState2 types.State_NG = &types.SelectLowerBoundsState{
			OpponentId: &userData.ID,
		}
		return []types.ReplyDTO{
			{
				UserId: userData.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: *opponentId,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState2,
			},
		}, nil

	case *types.SelectLowerBoundsState:
		lower_bound, err := parseSalary(update.Message.Text)

		if err != nil {
			return []types.ReplyDTO{
				{
					UserId: userData.ID,
					Messages: []types.ReplyMessage{
						{
							Message:     MESSAGE_USE_VALID_POSITIVE_NUMBER,
							ReplyMarkup: nil,
						},
					},
				},
			}, nil
		} else {
			var nextState types.State_NG = &types.SelectUpperBoundsState{
				OpponentId: &userData.ID,
				LowerBound: &lower_bound,
			}
			return []types.ReplyDTO{
				{
					UserId: userData.ID,
					Messages: []types.ReplyMessage{
						{
							Message:     MESSAGE_SELECT_SALARY_UPPER_BOUND,
							ReplyMarkup: nil,
						},
					},
					NextState: &nextState,
				},
			}, nil
		}

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
