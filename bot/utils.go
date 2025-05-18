package bot

import (
	"fmt"
	"github.com/ikebastuz/tgn/types"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

const CONNECTION_PATTERN = `^\s*/connect\s+(\d+)\s*$`
const START_PATTERN = `(?s).*/start.*`
const RESET_PATTERN = `(?s).*/reset.*`

func CreateSelectUpperBoundMessage(role types.Role) string {
	var message string
	switch role {
	case types.ROLE_EMPLOYEE:
		message = MESSAGE_SELECT_SALARY_UPPER_BOUND_EMPLOYEE
	case types.ROLE_EMPLOYER:
		message = MESSAGE_SELECT_SALARY_UPPER_BOUND_EMPLOYER
	}
	return fmt.Sprintf("%s\n%s", message, CreateUseValidUpperBoundMessage())
}

func CreateUseValidUpperBoundMessage() string {
	return fmt.Sprintf(MESSAGE_USE_VALID_UPPER_BOUND, UPPER_BOUND_MULTIPLIER)
}

func CreateConnectionMessage(userName string, connectionId int16) string {
	return fmt.Sprintf(MESSAGE_FORWARD_CONNECTION_01, userName, connectionId)
}

func CreateResultMessage(salary int64) string {
	return fmt.Sprintf(MESSAGE_RESULT_SUCCESS, salary)
}

func CreateSelectLowerBoundsMessage(role types.Role) string {
	var message string
	switch role {
	case types.ROLE_EMPLOYEE:
		message = MESSAGE_SELECT_SALARY_LOWER_BOUND_EMPLOYEE
	case types.ROLE_EMPLOYER:
		message = MESSAGE_SELECT_SALARY_LOWER_BOUND_EMPLOYER
	}
	return fmt.Sprintf(message, role)
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

func getDialogState(userId int64, store types.Store) (*types.StateMachine, error) {
	dialogState := store.GetDialogState(&userId)
	return dialogState, nil
}

func getConnectionId(update *types.TelegramUpdate) (int16, bool) {
	text := strings.TrimSpace(update.Message.Text)

	re := regexp.MustCompile(CONNECTION_PATTERN)

	if re.MatchString(text) {
		matches := re.FindStringSubmatch(text)
		idStr := matches[1]

		id, err := strconv.ParseInt(idStr, 10, 16)
		if err != nil {
			return 0, false
		}
		return int16(id), true
	}

	return 0, false
}

func isStartMessage(update *types.TelegramUpdate) bool {
	text := strings.TrimSpace(update.Message.Text)

	re := regexp.MustCompile(START_PATTERN)

	return re.MatchString(text)
}

func isResetMessage(update *types.TelegramUpdate) bool {
	text := strings.TrimSpace(update.Message.Text)

	re := regexp.MustCompile(RESET_PATTERN)

	return re.MatchString(text)
}

func parseSalary(message string) (int64, error) {
	salaryStr := strings.TrimSpace(message)

	value, err := strconv.ParseInt(salaryStr, 10, 64)

	if err != nil {
		return 0, err
	}

	return value, nil
}

func resetUserState(userId *int64, store types.Store) types.ReplyDTO {
	store.ResetUserState(userId)

	return types.ReplyDTO{
		UserId: *userId,
		Messages: []types.ReplyMessage{
			{
				Message:     MESSAGE_START_GUIDE,
				ReplyMarkup: nil,
			},
		},
	}
}

func createStartReply(store types.Store, userData *types.From) []types.ReplyDTO {
	log.Info("Received START message")
	newConnectionId := store.CreateConnectionId(&userData.ID)

	var nextState types.State = &types.WaitingForConnectState{
		ConnectionId: &newConnectionId,
	}

	return []types.ReplyDTO{
		{
			UserId: userData.ID,
			Messages: []types.ReplyMessage{
				{
					Message:     CreateConnectionMessage(userData.USERNAME, newConnectionId),
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
	}
}
