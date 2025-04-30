package bot

import (
	"fmt"
	"github.com/ikebastuz/tgn/types"
	"regexp"
	"strconv"
	"strings"
)

const CONNECTION_PATTERN = `^\s*/connect\s+(\d+)\s*$`
const START_PATTERN = `(?s).*/start.*`
const RESET_PATTERN = `(?s).*/reset.*`

func createConnectionMessage(userName string, connectionId int64) string {
	return fmt.Sprintf(MESSAGE_FORWARD_CONNECTION_01, userName, connectionId)
}

func createResultMessage(salary int64) string {
	return fmt.Sprintf(MESSAGE_RESULT, salary)
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

func getConnectionId(update *types.TelegramUpdate) (int64, bool) {
	text := strings.TrimSpace(update.Message.Text)

	re := regexp.MustCompile(CONNECTION_PATTERN)

	if re.MatchString(text) {
		matches := re.FindStringSubmatch(text)
		idStr := matches[1]

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return 0, false
		}
		return id, true
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
