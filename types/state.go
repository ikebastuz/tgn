package types

type State string

const (
	STATE_INITIAL       = "initial"
	WAITING_FOR_CONNECT = "waiting-for-connect"
	SELECT_YOUR_ROLE    = "select-your-role"
	SELECT_LOWER_BOUNDS = "select-lower-bounds"
	SELECT_UPPER_BOUNDS = "select-upper-bounds"
	RESULT              = "result"
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
