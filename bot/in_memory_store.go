package bot

import (
	"errors"
	"math/rand"

	"github.com/ikebastuz/tgn/types"
)

type InMemoryStore struct {
	states      map[int64]*types.DialogState
	connections map[int64]int64
}

var DEFAULT_DIALOG_STATE = types.DialogState{
	State: types.STATE_INITIAL,
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		states:      make(map[int64]*types.DialogState),
		connections: make(map[int64]int64),
	}
}

func (s *InMemoryStore) GetDialogState(userId *int64) *types.DialogState {
	if state, exists := s.states[*userId]; exists {
		return state
	}
	s.states[*userId] = &DEFAULT_DIALOG_STATE

	return &DEFAULT_DIALOG_STATE
}

func (s *InMemoryStore) SetDialogState(userId *int64, state types.DialogState) {
	s.states[*userId] = &state
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

	currentState := s.states[*userId]
	currentState.ConnectionId = &id

	return id
}

func (s *InMemoryStore) GetConnectionTarget(connectionId *int64) *int64 {
	if connection, exists := s.connections[*connectionId]; exists {
		return &connection
	}
	return nil
}

func (s *InMemoryStore) DeleteConnectionId(connectionId *int64) error {
	_, exists := s.connections[*connectionId]
	if !exists {
		return errors.New("no such connection")
	} else {
		delete(s.connections, *connectionId)
	}

	for _, state := range s.states {
		if state.ConnectionId != nil && *state.ConnectionId == *connectionId {
			state.ConnectionId = nil
		}
	}

	return nil
}

func (s *InMemoryStore) ResetUserState(userId *int64) error {
	s.states[*userId] = &DEFAULT_DIALOG_STATE

	for connectionId, targetUserId := range s.connections {
		if targetUserId == *userId {
			delete(s.connections, connectionId)
		}
	}
	return nil
}
