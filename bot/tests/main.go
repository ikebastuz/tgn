package bot_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ikebastuz/tgn/types"
)

var TEST_USER_ID int64 = 123
var TEST_USER_ID_2 int64 = 124
var TEST_FROM = types.From{
	ID:       int64(TEST_USER_ID),
	USERNAME: "hello",
}
var TEST_FROM_2 = types.From{
	ID:       int64(TEST_USER_ID_2),
	USERNAME: "hello 2",
}
var TEST_CONNECTION_ID int64 = 100500

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
