package types

type State int

const (
	STATE_INITIAL = iota
	WAITING_FOR_CONNECT
)

type DialogState struct {
	State           State
	AnchorMessageId *int
	OpponentId      *int64
	LowerBound      *int64
	UpperBound      *int64
}
