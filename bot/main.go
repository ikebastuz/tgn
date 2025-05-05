package bot

import (
	"context"
	"math/rand"

	log "github.com/sirupsen/logrus"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/ikebastuz/tgn/actions"
	"github.com/ikebastuz/tgn/metrics"
	"github.com/ikebastuz/tgn/solver"
	"github.com/ikebastuz/tgn/types"
)

const UPPER_BOUND_MULTIPLIER int64 = 3

func HandleMessage(ctx context.Context, client *telegram.Client, update types.TelegramUpdate, store types.Store) error {
	metrics.RequestCounter.Inc()
	replies, err := CreateReply(update, store)
	if err != nil {
		metrics.ErrorCounter.Inc()
		log.Errorf("❌ Failed to create reply: %v", err)
		return err
	}

	for _, reply := range replies {
		metrics.ReplyCounter.Inc()
		for _, message := range reply.Messages {
			if message.MessageID > 0 {
				_, err = client.API().MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
					ID:          message.MessageID,
					Peer:        &tg.InputPeerUser{UserID: reply.UserId},
					Message:     message.Message,
					ReplyMarkup: message.ReplyMarkup,
				})

				if err != nil {
					metrics.ErrorCounter.Inc()
					log.Errorf("❌ Message update failed: Could not edit message: %v", err)
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
					metrics.ErrorCounter.Inc()
					log.Errorf("❌ Message delivery failed: Could not send message: %v", err)
					return err
				}
			}
		}

		if reply.NextState != nil {
			metrics.StateTransitionCounter.Inc()
			log.Infof("👤 User %d updates state: %T", reply.UserId, *reply.NextState)
			store.SetDialogState(&reply.UserId, *reply.NextState)

		}
	}

	metrics.UsersStoreCounter.Set(float64(store.GetUsersCount()))

	return nil
}

func CreateReply(update types.TelegramUpdate, store types.Store) ([]types.ReplyDTO, error) {
	userData, err := getUserData(update)
	if err != nil {
		metrics.ErrorCounter.Inc()
		return nil, err
	}

	sm, err := getDialogState(userData.ID, store)
	if err != nil {
		metrics.ErrorCounter.Inc()
		return nil, err
	}

	if isResetMessage(&update) {
		metrics.ResetCounter.Inc()
		log.Infof("Received RESET message from USER %v", userData.ID)

		response := []types.ReplyDTO{}
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
		}

		log.Infof("Resetting user %d state", userData.ID)
		response = append(response, resetUserState(&userData.ID, store))
		return response, nil
	}

	log.Infof("👤 User %d in state: %T", userData.ID, sm.GetState())
	switch s := sm.GetState().(type) {
	case *types.InitialState:
		incomingConnectionId, isConnectionMessage := getConnectionId(&update)

		if isStartMessage(&update) {
			reply := createStartReply(store, userData)
			return reply, nil
		} else if isConnectionMessage {
			metrics.ConnectAttemptCounter.Inc()
			log.Info("Received CONNECT message")
			targetUserId := store.GetConnectionTarget(&incomingConnectionId)
			if targetUserId == nil {
				metrics.ErrorCounter.Inc()
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
				metrics.ErrorCounter.Inc()
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
			metrics.ErrorCounter.Inc()
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
			metrics.RoleSelectCounter.Inc()
		case actions.ACTION_SELECT_EMPLOYER:
			nextRole1 = types.ROLE_EMPLOYER
			nextRole2 = types.ROLE_EMPLOYEE
			metrics.RoleSelectCounter.Inc()
		default:
			metrics.ErrorCounter.Inc()
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
						Message:     CreateSelectLowerBoundsMessage(nextRole1),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: *s.OpponentId,
				Messages: []types.ReplyMessage{
					{
						Message:     CreateSelectLowerBoundsMessage(nextRole2),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState2,
			},
		}, nil
	case *types.SelectLowerBoundsState:
		lower_bound, err := parseSalary(update.Message.Text)

		if err != nil {
			metrics.SalaryParseErrorCounter.Inc()
			metrics.ErrorCounter.Inc()
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
							Message:     CreateSelectUpperBoundMessage(s.Role),
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
			metrics.SalaryParseErrorCounter.Inc()
			metrics.ErrorCounter.Inc()
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
			if upper_bound > *s.LowerBound*UPPER_BOUND_MULTIPLIER {
				metrics.ErrorCounter.Inc()
				return []types.ReplyDTO{
					{
						UserId: userData.ID,
						Messages: []types.ReplyMessage{
							{
								Message:     CreateUseValidUpperBoundMessage(),
								ReplyMarkup: nil,
							},
						},
					},
				}, nil
			}

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
					metrics.ResultErrorCounter.Inc()
					metrics.ErrorCounter.Inc()
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
				metrics.ResultSuccessCounter.Inc()
				var nextState1 types.State = &types.InitialState{}
				var nextState2 types.State = &types.InitialState{}

				return []types.ReplyDTO{
					{
						UserId: userData.ID,
						Messages: []types.ReplyMessage{
							{
								Message:     CreateResultMessage(salary),
								ReplyMarkup: nil,
							},
						},
						NextState: &nextState1,
					},
					{
						UserId: *s.OpponentId,
						Messages: []types.ReplyMessage{
							{
								Message:     CreateResultMessage(salary),
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
	case *types.ResultErrorState:
		switch update.CallbackQuery.Data {
		case actions.ACTION_SELECT_NO:
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
		case actions.ACTION_SELECT_YES:
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
							Message:     CreateSelectLowerBoundsMessage(s.Role),
							ReplyMarkup: nil,
						},
					},
					NextState: &nextState1,
				},
				{
					UserId: *s.OpponentId,
					Messages: []types.ReplyMessage{
						{
							Message:     CreateSelectLowerBoundsMessage(opponentRole),
							ReplyMarkup: nil,
						},
					},
					NextState: &nextState2,
				},
			}, nil
		default:
			metrics.ErrorCounter.Inc()
			return []types.ReplyDTO{}, nil
		}
	default:
		metrics.ErrorCounter.Inc()
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
