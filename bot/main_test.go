package bot

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ikebastuz/tgn/actions"
	"github.com/ikebastuz/tgn/types"
)

func TestCreateReply(t *testing.T) {
	var USER_ID int64 = 123
	var USER_ID_2 int64 = 124
	var FROM = types.From{
		ID:       int64(USER_ID),
		USERNAME: "hello",
	}
	var FROM_2 = types.From{
		ID:       int64(USER_ID_2),
		USERNAME: "hello 2",
	}
	var CONNECTION_ID int64 = 100500

	t.Run("INITIAL state, irrelevant message - should show guide", func(t *testing.T) {
		store := NewInMemoryStore()

		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				From: FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("INITIAL state, /start message - create connection to forward", func(t *testing.T) {
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.InitialState{})
		store.states[USER_ID] = &sm

		var nextState types.State_NG = &types.WaitingForConnectState{
			ConnectionId: &CONNECTION_ID,
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     createConnectionMessage(FROM.USERNAME, CONNECTION_ID),
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState,
			},
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_FORWARD_CONNECTION_02,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				From: FROM,
				Text: " /start ",
			},
		}

		got, err := createReply(update, store)

		assertSingleReply(t, got[1], want[1], err)
	})

	t.Run("INITIAL state, /connect message to yourself - should show an error", func(t *testing.T) {
		store := NewInMemoryStore()
		store.connections[CONNECTION_ID] = USER_ID

		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_YOU_CANT_CONNECT_TO_YOURSELF,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				Text: fmt.Sprintf("/connect %d", CONNECTION_ID),
				From: FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("INITIAL state, /connect message to non-existent user - should show an error", func(t *testing.T) {
		store := NewInMemoryStore()
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_NO_SUCH_USER_IS_AWATING,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				Text: "/connect 1337",
				From: FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("INITIAL state, /connect message to existing user - should connect correctly", func(t *testing.T) {
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.WaitingForConnectState{
			ConnectionId: &CONNECTION_ID,
		})
		store.states[USER_ID] = &sm
		store.connections[CONNECTION_ID] = USER_ID

		var nextState1 types.State_NG = &types.SelectYourRoleState{
			OpponentId: &FROM_2.ID,
		}
		var nextState2 types.State_NG = &types.SelectYourRoleState{
			OpponentId: &FROM.ID,
		}

		want := []types.ReplyDTO{
			{
				UserId: FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_YOUR_ROLE_CONNECTED,
						ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
					},
				},
				NextState: &nextState2,
			},
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_YOUR_ROLE_CONNECTED,
						ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
					},
				},
				NextState: &nextState1,
			},
		}

		update := types.TelegramUpdate{
			Message: types.Message{
				Text: fmt.Sprintf("/connect %d", CONNECTION_ID),
				From: FROM_2,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})
	//
	t.Run("WAITING state, - tells about waiting for connection", func(t *testing.T) {
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.WaitingForConnectState{
			ConnectionId: &CONNECTION_ID,
		})
		store.states[FROM.ID] = &sm

		var FROM = types.From{
			ID:       int64(USER_ID),
			USERNAME: "hello",
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_WAITING_FOR_CONNECTION,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      "",
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}

		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT ROLE state - selected EMPLOYEE - update both users and ask for lower bounds", func(t *testing.T) {
		store := NewInMemoryStore()
		sm1 := types.StateMachine{}
		sm1.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID_2,
		})
		store.states[USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID,
		})
		store.states[USER_ID_2] = &sm2

		var nextState1 types.State_NG = &types.SelectLowerBoundsState{
			OpponentId: &USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
		}
		var nextState2 types.State_NG = &types.SelectLowerBoundsState{
			OpponentId: &USER_ID,
			Role:       types.ROLE_EMPLOYER,
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState2,
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			CallbackQuery: types.CallbackQuery{
				Data: actions.ACTION_SELECT_EMPLOYEE,
			},
			Message: types.Message{
				MessageID: 1,
				Text:      "",
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT ROLE state - selected EMPLOYER - update both users and ask for lower bounds", func(t *testing.T) {
		store := NewInMemoryStore()
		sm1 := types.StateMachine{}
		sm1.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID_2,
		})
		store.states[USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID,
		})
		store.states[USER_ID_2] = &sm2

		var nextState1 types.State_NG = &types.SelectLowerBoundsState{
			OpponentId: &USER_ID_2,
			Role:       types.ROLE_EMPLOYER,
		}
		var nextState2 types.State_NG = &types.SelectLowerBoundsState{
			OpponentId: &USER_ID,
			Role:       types.ROLE_EMPLOYEE,
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState1,
			},
			{
				UserId: FROM_2.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_LOWER_BOUND,
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState2,
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			CallbackQuery: types.CallbackQuery{
				Data: actions.ACTION_SELECT_EMPLOYER,
			},
			Message: types.Message{
				MessageID: 1,
				Text:      "",
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT ROLE state - received unexpected data - prompt for role again", func(t *testing.T) {
		store := NewInMemoryStore()
		sm1 := types.StateMachine{}
		sm1.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID_2,
		})
		store.states[USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID,
		})
		store.states[USER_ID_2] = &sm2

		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_YOUR_ROLE_UNEXPECTED,
						ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			CallbackQuery: types.CallbackQuery{
				Data: " unexpected value ",
			},
			Message: types.Message{
				MessageID: 1,
				Text:      "",
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT ROLE state - received text message - prompt for role again", func(t *testing.T) {
		store := NewInMemoryStore()
		sm1 := types.StateMachine{}
		sm1.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID_2,
		})
		store.states[USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectYourRoleState{
			OpponentId: &USER_ID,
		})
		store.states[USER_ID_2] = &sm2

		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_YOUR_ROLE_UNEXPECTED,
						ReplyMarkup: KEYBOARD_SELECT_YOUR_ROLE,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      " some text ",
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}
		assertNonErrorReply(t, got, want, err)
	})
	//
	t.Run("SELECT LOWER BOUNDS state - show error message on invalid value", func(t *testing.T) {
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.SelectLowerBoundsState{
			OpponentId: &FROM_2.ID,
			Role:       types.ROLE_EMPLOYEE,
		})
		store.states[USER_ID] = &sm
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_USE_VALID_POSITIVE_NUMBER,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      "qweasd",
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("SELECT LOWER BOUNDS state - proceed further to upper bounds state if number is correct", func(t *testing.T) {
		var lower_bound int64 = 100500
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.SelectLowerBoundsState{
			OpponentId: &USER_ID_2,
			Role:       types.ROLE_EMPLOYEE, // TODO: make it required
		})
		store.states[USER_ID] = &sm

		var nextState types.State_NG = &types.SelectUpperBoundsState{
			OpponentId: &USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
			LowerBound: &lower_bound,
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_SELECT_SALARY_UPPER_BOUND,
						ReplyMarkup: nil,
					},
				},
				NextState: &nextState,
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      fmt.Sprintf("%d", lower_bound),
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		assertNonErrorReply(t, got, want, err)
	})

	t.Run("ANY state with /reset - should set to initial state", func(t *testing.T) {
		store := NewInMemoryStore()
		sm1 := types.StateMachine{}
		sm1.SetState(&types.SelectLowerBoundsState{
			OpponentId: &USER_ID_2,
			Role:       types.ROLE_EMPLOYEE,
		})
		store.states[USER_ID] = &sm1

		sm2 := types.StateMachine{}
		sm2.SetState(&types.SelectLowerBoundsState{
			OpponentId: &USER_ID,
			Role:       types.ROLE_EMPLOYER,
		})
		store.states[USER_ID_2] = &sm2

		var FROM = types.From{
			ID:       int64(USER_ID),
			USERNAME: "hello",
		}
		want := []types.ReplyDTO{
			{
				UserId: USER_ID_2,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
			},
			{
				UserId: USER_ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_START_GUIDE,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      " /reset ",
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}

		assertNonErrorReply(t, got, want, err)

		wantState := &types.StateMachine{}
		wantState.SetState(&types.InitialState{})

		gotUserState := store.GetDialogState(&USER_ID)
		gotUser2State := store.GetDialogState(&USER_ID_2)

		assertState(t, *gotUserState, *wantState)
		assertState(t, *gotUser2State, *wantState)
	})

	t.Run("UNEXPECTED state - should suggest the user to reset the state", func(t *testing.T) {
		store := NewInMemoryStore()
		sm := types.StateMachine{}
		sm.SetState(&types.UnexpectedState{})
		store.states[USER_ID] = &sm

		var FROM = types.From{
			ID:       int64(USER_ID),
			USERNAME: "hello",
		}
		want := []types.ReplyDTO{
			{
				UserId: FROM.ID,
				Messages: []types.ReplyMessage{
					{
						Message:     MESSAGE_UNEXPECTED_STATE,
						ReplyMarkup: nil,
					},
				},
			},
		}

		update := types.TelegramUpdate{
			UpdateID: 1,
			Message: types.Message{
				MessageID: 1,
				Text:      "",
				From:      FROM,
			},
		}

		got, err := createReply(update, store)
		if err != nil {
			t.Errorf("shouldn't have error")
		}

		assertNonErrorReply(t, got, want, err)
	})
}

func assertNonErrorReply(t testing.TB, got, want []types.ReplyDTO, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("shouldn't have error")
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func assertSingleReply(t testing.TB, got, want types.ReplyDTO, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("shouldn't have error")
	}

	if diff := cmp.Diff([]types.ReplyDTO{want}, []types.ReplyDTO{got}); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func assertState(t testing.TB, got, want types.StateMachine) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected state \n%v, \ngot state\n%v", want, got)
	}
}
