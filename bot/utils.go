package bot

import (
	"fmt"
)

func createConnectionMessage(userId int64) string {
	return fmt.Sprintf("%s %v", FORWARD_CONNECTION_MESSAGE, 1)
}
