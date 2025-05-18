package bot

import (
	"errors"
	"math/rand"

	"github.com/ikebastuz/tgn/types"
)

type InMemoryStore struct {
	states      map[int64]*types.StateMachine
	connections map[int16]int64
}

var DEFAULT_DIALOG_STATE = types.InitialState{}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		states:      make(map[int64]*types.StateMachine),
		connections: make(map[int16]int64),
	}
}

func (s *InMemoryStore) GetDialogState(userId *int64) *types.StateMachine {
	if sm, exists := s.states[*userId]; exists {
		return sm
	}
	sm := &types.StateMachine{}
	sm.SetState(&DEFAULT_DIALOG_STATE)
	s.states[*userId] = sm

	return sm
}

func (s *InMemoryStore) SetDialogState(userId *int64, state types.State) {
	if sm, exists := s.states[*userId]; exists {
		sm.SetState(state)
		return
	}

	sm := &types.StateMachine{}
	sm.SetState(state)
	s.states[*userId] = sm
}

func (s *InMemoryStore) CreateConnectionId(userId *int64) int16 {
	// Check if connection exists
	if sm, exists := s.states[*userId]; exists {
		s := sm.GetState()
		switch state := s.(type) {
		case *types.WaitingForConnectState:
			return *state.ConnectionId
		}
	}

	var id int16
	for {
		id = int16(rand.Int63n(10000) + 1)
		if _, exists := s.connections[id]; !exists {
			break
		}
	}
	s.connections[id] = *userId
	return id
}

func (s *InMemoryStore) GetConnectionTarget(connectionId *int16) *int64 {
	if connection, exists := s.connections[*connectionId]; exists {
		return &connection
	}
	return nil
}

func (s *InMemoryStore) DeleteConnectionId(connectionId *int16) error {
	_, exists := s.connections[*connectionId]
	if !exists {
		return errors.New("no such connection")
	} else {
		delete(s.connections, *connectionId)
	}

	return nil
}

func (s *InMemoryStore) ResetUserState(userId *int64) error {
	sm := s.states[*userId]
	sm.SetState(&DEFAULT_DIALOG_STATE)

	for connectionId, targetUserId := range s.connections {
		if targetUserId == *userId {
			delete(s.connections, connectionId)
		}
	}
	return nil
}

func (s *InMemoryStore) GetUsersCount() int {
	return len(s.states)
}

func (s *InMemoryStore) SetConnectionId(from int16, to int64) {
	s.connections[from] = to
}
