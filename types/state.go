package types

type State int

const (
	STATE_INITIAL = iota
	WAITING_FOR_CONNECT
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
