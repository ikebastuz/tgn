package bot

import (
	"github.com/ikebastuz/tgn/types"
	"math/rand"
)

type InMemoryStore struct {
	states      map[int64]types.DialogState
	connections map[int64]int64
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		states:      make(map[int64]types.DialogState),
		connections: make(map[int64]int64),
	}
}

func (s *InMemoryStore) GetDialogState(userId *int64) *types.DialogState {
	if state, exists := s.states[*userId]; exists {
		return &state
	}
	return &types.DialogState{
		State: types.STATE_INITIAL,
	}
}

func (s *InMemoryStore) SetDialogState(userId *int64, state *types.DialogState) {
	s.states[*userId] = *state
}

func (s *InMemoryStore) CreateConnectionId(userId *int64) int64 {
	// Check if connection exists
	if state, exists := s.states[*userId]; exists {
		if state.ConnectionId != nil {
			return *state.ConnectionId
		}
	}

	var id int64
	for {
		id = rand.Int63n(10000) + 1
		if _, exists := s.connections[id]; !exists {
			break
		}
	}
	s.connections[id] = *userId
	return id
}

func (s *InMemoryStore) GetConnectionTarget(connectionId *int64) *int64 {
	if connection, exists := s.connections[*connectionId]; exists {
		return &connection
	}
	return nil
}
