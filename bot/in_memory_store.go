package bot

import (
	"github.com/ikebastuz/tgn/types"
)

type InMemoryStore struct {
	states map[int64]types.DialogState
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		states: make(map[int64]types.DialogState),
	}
}

func (s *InMemoryStore) GetDialogState(userId int64) *types.DialogState {
	if state, exists := s.states[userId]; exists {
		return &state
	}
	return &types.DialogState{
		State: types.STATE_INITIAL,
	}
}

func (s *InMemoryStore) SetDialogState(userId int64, state *types.DialogState) {
	s.states[userId] = *state
}
