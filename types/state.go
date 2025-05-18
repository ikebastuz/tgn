package types

//// Partial migration to state machine

type StateType string
type Role string

type State interface {
	GetState() StateType
}

const (
	STATE_INITIAL             StateType = "initial"
	STATE_WAITING_FOR_CONNECT StateType = "waiting-for-connect"
	STATE_SELECT_YOUR_ROLE    StateType = "select-your-role"
	STATE_SELECT_LOWER_BOUNDS StateType = "select-lower-bounds"
	STATE_SELECT_UPPER_BOUNDS StateType = "select-upper-bounds"
	STATE_WAITING_FOR_RESULT  StateType = "waiting-for-result"
	STATE_RESULT_ERROR        StateType = "result-error"
	STATE_UNEXPECTED          StateType = "unexpected"
)

const (
	ROLE_EMPLOYER Role = "employer"
	ROLE_EMPLOYEE Role = "employee"
)

type StateMachine struct {
	current State
}

func (sm *StateMachine) SetState(s State) {
	sm.current = s
}

func (sm *StateMachine) GetState() State {
	return sm.current
}

type InitialState struct{}

func (s *InitialState) GetState() StateType {
	return STATE_INITIAL
}

type WaitingForConnectState struct {
	ConnectionId *int16
}

func (s *WaitingForConnectState) GetState() StateType {
	return STATE_WAITING_FOR_CONNECT
}

type SelectYourRoleState struct {
	OpponentId *int64
}

func (s *SelectYourRoleState) GetState() StateType {
	return STATE_SELECT_YOUR_ROLE
}

type SelectLowerBoundsState struct {
	OpponentId *int64
	Role       Role
}

func (s *SelectLowerBoundsState) GetState() StateType {
	return STATE_SELECT_LOWER_BOUNDS
}

type SelectUpperBoundsState struct {
	OpponentId *int64
	Role       Role
	LowerBound *int64
}

func (s *SelectUpperBoundsState) GetState() StateType {
	return STATE_SELECT_UPPER_BOUNDS
}

type WaitingForResultState struct {
	OpponentId *int64
	Role       Role
	LowerBound *int64
	UpperBound *int64
}

func (s *WaitingForResultState) GetState() StateType {
	return STATE_WAITING_FOR_RESULT
}

type ResultErrorState struct {
	OpponentId *int64
	Role       Role
}

func (s *ResultErrorState) GetState() StateType {
	return STATE_RESULT_ERROR
}

type UnexpectedState struct{}

func (s *UnexpectedState) GetState() StateType {
	return STATE_UNEXPECTED
}
