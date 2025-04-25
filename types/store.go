package types

type Store interface {
	GetConnectionTarget(connectionId *int64) *int64
	CreateConnectionId(userId *int64) int64
	DeleteConnectionId(connectionId *int64) error
	GetDialogState(userId *int64) *DialogState
	SetDialogState(userId *int64, state *DialogState)
}
