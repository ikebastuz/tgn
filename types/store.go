package types

type Store interface {
	GetDialogState(userId int64) *DialogState
	SetDialogState(userId int64, state *DialogState)
}
