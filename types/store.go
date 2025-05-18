package types

type Store interface {
	GetConnectionTarget(connectionId *int16) *int64
	CreateConnectionId(userId *int64) int16
	DeleteConnectionId(connectionId *int16) error
	GetDialogState(userId *int64) *StateMachine
	ResetUserState(userId *int64) error
	SetDialogState(userId *int64, state State)
	GetUsersCount() int
	SetConnectionId(from int16, to int64)
}
