package bot

import (
	"fmt"
	"github.com/ikebastuz/tgn/types"
	"regexp"
	"strconv"
	"strings"
)

const CONNECTION_PATTERN = `^\s*/connect\s+(\d+)\s*$`

func createConnectionMessage(userData types.From) string {
	return fmt.Sprintf(MESSAGE_FORWARD_CONNECTION_01, userData.USERNAME, userData.ID)
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

func getDialogState(userId int64, store types.Store) (*types.DialogState, error) {
	dialogState := store.GetDialogState(userId)
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
