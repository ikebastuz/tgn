package bot

import (
	"context"
	"math/rand"

	log "github.com/sirupsen/logrus"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/ikebastuz/tgn/actions"
	"github.com/ikebastuz/tgn/solver"
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
		response := []types.ReplyDTO{}
		log.Infof("debug %v", sm.GetState().GetState())
		switch s := sm.GetState().(type) {
		case *types.SelectUpperBoundsState:
			log.Infof("Resetting opponent %d state aswell", *s.OpponentId)
			response = append(response, resetUserState(s.OpponentId, store))
		case *types.SelectLowerBoundsState:
			log.Infof("Resetting opponent %d state aswell", *s.OpponentId)
			response = append(response, resetUserState(s.OpponentId, store))
		case *types.SelectYourRoleState:
			log.Infof("Resetting opponent %d state aswell", *s.OpponentId)
			response = append(response, resetUserState(s.OpponentId, store))
		default:
			log.Info("default")
		}

		log.Infof("Resetting user %d state", userData.ID)
		response = append(response, resetUserState(&userData.ID, store))
		return response, nil
	}

	// log.Infof("CURRENT STATE: %v", sm.GetState().GetState())
	switch s := sm.GetState().(type) {
	case *types.InitialState:
		incomingConnectionId, isConnectionMessage := getConnectionId(&update)

		if isStartMessage(&update) {
			log.Info("Received START message")

			// Start message
			// TODO: check if not connection exists
			newConnectionId := store.CreateConnectionId(&userData.ID)

			var nextState types.State = &types.WaitingForConnectState{
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

				var nextState1 types.State = &types.SelectYourRoleState{
					OpponentId: targetUserId,
				}
				var nextState2 types.State = &types.SelectYourRoleState{
					OpponentId: &userData.ID,
				}
				return []types.ReplyDTO{
					{
						UserId: userData.ID,
						Messages: []types.ReplyMessage{
							{
								Message:     MESSAGE_SELECT_YOUR_ROLE_CONNECTED,
								ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
							},
						},
						NextState: &nextState1,
					},
					{
						UserId: *targetUserId,
						Messages: []types.ReplyMessage{
							{
								Message:     MESSAGE_SELECT_YOUR_ROLE_CONNECTED,
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
		var nextRole1 types.Role
		var nextRole2 types.Role

		switch update.CallbackQuery.Data {
		case actions.ACTION_SELECT_EMPLOYEE:
			nextRole1 = types.ROLE_EMPLOYEE
			nextRole2 = types.ROLE_EMPLOYER
		case actions.ACTION_SELECT_EMPLOYER:
			nextRole1 = types.ROLE_EMPLOYER
			nextRole2 = types.ROLE_EMPLOYEE
		default:
			return []types.ReplyDTO{
				{
					UserId: userData.ID,
					Messages: []types.ReplyMessage{
						{
							Message:     MESSAGE_SELECT_YOUR_ROLE_UNEXPECTED,
							ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
						},
					},
				},
			}, nil
		}

		var nextState1 types.State = &types.SelectLowerBoundsState{
			OpponentId: s.OpponentId,
			Role:       nextRole1,
		}
		var nextState2 types.State = &types.SelectLowerBoundsState{
			OpponentId: &userData.ID,
			Role:       nextRole2,
		}
		return []types.ReplyDTO{
			{
				UserId: userData.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     createSelectLowerBoundsMessage(nextRole1),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: *s.OpponentId,
				Messages: []types.ReplyMessage{
					{
						Message:     createSelectLowerBoundsMessage(nextRole2),
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
			var nextState types.State = &types.SelectUpperBoundsState{
				OpponentId: s.OpponentId,
				Role:       s.Role,
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
	case *types.SelectUpperBoundsState:
		upper_bound, err := parseSalary(update.Message.Text)

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
			opponentState := store.GetDialogState(s.OpponentId)
			switch os := opponentState.GetState().(type) {
			case *types.WaitingForResultState:
				var employeeRange solver.Range
				var employerRange solver.Range

				switch s.Role {
				case types.ROLE_EMPLOYEE:
					employeeRange = solver.Range{
						Min: *s.LowerBound,
						Max: upper_bound,
					}
					employerRange = solver.Range{
						Min: *os.LowerBound,
						Max: *os.UpperBound,
					}

				case types.ROLE_EMPLOYER:
					employerRange = solver.Range{
						Min: *s.LowerBound,
						Max: upper_bound,
					}
					employeeRange = solver.Range{
						Min: *os.LowerBound,
						Max: *os.UpperBound,
					}
				}

				salary, err := solver.Solve(employeeRange, employerRange)
				if err != nil {
					var nextState1 types.State = &types.ResultErrorState{
						OpponentId: s.OpponentId,
						Role:       s.Role,
					}
					var nextState2 types.State = &types.ResultErrorState{
						OpponentId: &userData.ID,
						Role:       os.Role,
					}
					return []types.ReplyDTO{
						{
							UserId: userData.ID,
							Messages: []types.ReplyMessage{
								{
									Message:     MESSAGE_RESULT_ERROR,
									ReplyMarkup: KEYBOARD_SELECT_YES_NO,
								},
							},
							NextState: &nextState1,
						},
						{
							UserId: *s.OpponentId,
							Messages: []types.ReplyMessage{
								{
									Message:     MESSAGE_RESULT_ERROR,
									ReplyMarkup: KEYBOARD_SELECT_YES_NO,
								},
							},
							NextState: &nextState2,
						},
					}, nil
				}
				var nextState1 types.State = &types.ResultSuccessState{
					OpponentId: s.OpponentId,
					Role:       s.Role,
					LowerBound: s.LowerBound,
					UpperBound: &upper_bound,
					Result:     &salary,
				}
				var nextState2 types.State = &types.ResultSuccessState{
					OpponentId: os.OpponentId,
					Role:       os.Role,
					LowerBound: os.LowerBound,
					UpperBound: os.UpperBound,
					Result:     &salary,
				}

				return []types.ReplyDTO{
					{
						UserId: userData.ID,
						Messages: []types.ReplyMessage{
							{
								Message:     createResultMessage(salary),
								ReplyMarkup: nil,
							},
						},
						NextState: &nextState1,
					},
					{
						UserId: *s.OpponentId,
						Messages: []types.ReplyMessage{
							{
								Message:     createResultMessage(salary),
								ReplyMarkup: nil,
							},
						},
						NextState: &nextState2,
					},
				}, nil

			default:
				var nextState types.State = &types.WaitingForResultState{
					OpponentId: s.OpponentId,
					Role:       s.Role,
					LowerBound: s.LowerBound,
					UpperBound: &upper_bound,
				}
				return []types.ReplyDTO{
					{
						UserId: userData.ID,
						Messages: []types.ReplyMessage{
							{
								Message:     MESSAGE_WAITING_FOR_RESULT,
								ReplyMarkup: nil,
							},
						},
						NextState: &nextState,
					},
				}, nil
			}
		}

	case *types.ResultSuccessState:
		var nextState types.State = &types.ResultSuccessState{
			OpponentId: s.OpponentId,
			Role:       s.Role,
			LowerBound: s.LowerBound,
			UpperBound: s.UpperBound,
			Result:     s.Result,
		}
		return []types.ReplyDTO{
			{
				UserId: userData.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState,
			},
		}, nil
	case *types.ResultErrorState:
		if update.CallbackQuery.Data == actions.ACTION_SELECT_NO {
			var nextState types.State = &types.InitialState{}
			return []types.ReplyDTO{
				{
					UserId: userData.ID,
					Messages: []types.ReplyMessage{
						{
							Message:     MESSAGE_START_GUIDE,
							ReplyMarkup: nil,
						},
					},
					NextState: &nextState,
				},
				{
					UserId: *s.OpponentId,
					Messages: []types.ReplyMessage{
						{
							Message:     MESSAGE_START_GUIDE,
							ReplyMarkup: nil,
						},
					},
					NextState: &nextState,
				},
			}, nil
		} else if update.CallbackQuery.Data == actions.ACTION_SELECT_YES {
			var nextState1 types.State = &types.SelectLowerBoundsState{
				OpponentId: s.OpponentId,
				Role:       s.Role,
			}
			var opponentRole types.Role
			if s.Role == types.ROLE_EMPLOYEE {
				opponentRole = types.ROLE_EMPLOYER
			}
			if s.Role == types.ROLE_EMPLOYER {
				opponentRole = types.ROLE_EMPLOYEE
			}
			var nextState2 types.State = &types.SelectLowerBoundsState{
				OpponentId: &userData.ID,
				Role:       opponentRole,
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
					UserId: *s.OpponentId,
					Messages: []types.ReplyMessage{
						{
							Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
							ReplyMarkup: nil,
						},
					},
					NextState: &nextState2,
				},
			}, nil
		} else {
			return []types.ReplyDTO{}, nil
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
