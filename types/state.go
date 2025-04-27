package types

type State string

const (
	STATE_INITIAL             = "initial"
	STATE_WAITING_FOR_CONNECT = "waiting-for-connect"
	STATE_SELECT_YOUR_ROLE    = "select-your-role"
	STATE_SELECT_LOWER_BOUNDS = "select-lower-bounds"
	STATE_SELECT_UPPER_BOUNDS = "select-upper-bounds"
	STATE_RESULT              = "result"
)

type DialogState struct {
	State           State
	ConnectionId    *int64
	AnchorMessageId *int
	OpponentId      *int64
	LowerBound      *int64
	UpperBound      *int64
}

type ConnectionState struct {
	RequesterId *int64
	ResponderId *int64
}
