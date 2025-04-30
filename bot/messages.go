package bot

import (
	"errors"
)

var (
	ErrorNoSenderIdFound = errors.New("no sender ID found in message")
)

const (
	MESSAGE_FORWARD_CONNECTION_01        = "%s invites you to negotiate\nTo join - send the following message to @NegotiMateBot\n\n/connect %v"
	MESSAGE_FORWARD_CONNECTION_02        = "Forward this ☝️ message to person\nyou want to negotiate with"
	MESSAGE_WAITING_FOR_CONNECTION       = "Waiting for connection"
	MESSAGE_YOU_CANT_CONNECT_TO_YOURSELF = "You can't connect to yourself"
	MESSAGE_NO_SUCH_USER_IS_AWATING      = "No such user is awaiting for connection"
	MESSAGE_START_GUIDE                  = "To use this bot use:\n/start - to initiate negotitaion\n/connect <ID> - to connect to person\n/reset - to reset"
	MESSAGE_SELECT_YOUR_ROLE_CONNECTED   = "You are connected!\nSelect your role"
	MESSAGE_SELECT_YOUR_ROLE_UNEXPECTED  = "Not sure what happened...\nPlease, select your role"
	MESSAGE_UNEXPECTED_STATE             = "Wow, I'm not sure how we ended up being here\nTry to reset with\n/reset"
	MESSAGE_SELECT_SALARY_LOWER_BOUND    = "Select your lower salary bounds"
	MESSAGE_SELECT_SALARY_UPPER_BOUND    = "Select your upper salary bounds"
	MESSAGE_WAITING_FOR_RESULT           = "Waiting for result..."
	MESSAGE_USE_VALID_POSITIVE_NUMBER    = "Use valid positive number"
	MESSAGE_RESULT                       = "Congratulations. You can agree on the amount of: %d"
)
