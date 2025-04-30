package types

type State string

const (
	STATE_INITIAL             = "initial"
	STATE_WAITING_FOR_CONNECT = "waiting-for-connect"
	STATE_SELECT_YOUR_ROLE    = "select-your-role"
	STATE_SELECT_LOWER_BOUNDS = "select-lower-bounds"
	STATE_SELECT_UPPER_BOUNDS = "select-upper-bounds"
	STATE_RESULT              = "result"
)

type DialogState struct {
	State        State
	ConnectionId *int64
	OpponentId   *int64
	LowerBound   *int64
	UpperBound   *int64
}

//// Partial migration to state machine

type StateType string
type Role string

type State_NG interface {
	GetState() StateType
}

const (
	STATE_INITIAL_NG             StateType = "initial"
	STATE_WAITING_FOR_CONNECT_NG StateType = "waiting-for-connect"
	STATE_SELECT_YOUR_ROLE_NG    StateType = "select-your-role"
	STATE_SELECT_LOWER_BOUNDS_NG StateType = "select-lower-bounds"
	STATE_SELECT_UPPER_BOUNDS_NG StateType = "select-upper-bounds"
	STATE_RESULT_NG              StateType = "result"
)

const (
	ROLE_EMPLOYER Role = "employer"
	ROLE_EMPLOYEE Role = "employee"
)

type StateMachine struct {
	current State_NG
}

func (sm *StateMachine) SetState(s State_NG) {
	sm.current = s
}

func (sm *StateMachine) GetState() State_NG {
	return sm.current
}

type InitialState struct{}

func (s *InitialState) GetState() StateType {
	return STATE_INITIAL_NG
}

type WaitingForConnectState struct {
	ConnectionId *int64
}

func (s *WaitingForConnectState) GetState() StateType {
	return STATE_WAITING_FOR_CONNECT_NG
}

type SelectYourRoleState struct {
	OpponentId *int64
}

func (s *SelectYourRoleState) GetState() StateType {
	return STATE_SELECT_YOUR_ROLE_NG
}

type SelectLowerBoundsState struct {
	OpponentId *int64
	Role       Role
}

func (s *SelectLowerBoundsState) GetState() StateType {
	return STATE_SELECT_LOWER_BOUNDS_NG
}

type SelectUpperBoundsState struct {
	OpponentId *int64
	Role       *Role
	LowerBound *int64
}

func (s *SelectUpperBoundsState) GetState() StateType {
	return STATE_SELECT_UPPER_BOUNDS_NG
}

type ResultState struct {
	OpponentId *int64
	Role       *Role
	LowerBound *int64
	UpperBound *int64
}

func (s *ResultState) GetState() StateType {
	return STATE_RESULT_NG
}
