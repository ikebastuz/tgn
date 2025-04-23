package bot

import (
	"errors"
)

var (
	ErrorNoSenderIdFound = errors.New("no sender ID found in message")
)

const (
	MESSAGE_FORWARD_CONNECTION_01  = "To join %s - type\n/connect %v"
	MESSAGE_FORWARD_CONNECTION_02  = "Forward this ☝️ message to person\nyou want to negotiate with"
	MESSAGE_WAITING_FOR_CONNECTION = "Waiting for connection"
)
