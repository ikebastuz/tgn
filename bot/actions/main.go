package actions

import (
	"github.com/gotd/td/tg"
)

const (
	ACTION_SELECT_EMPLOYEE = "employee"
	ACTION_SELECT_EMPLOYER = "employer"
	ACTION_SELECT_YES      = "yes"
	ACTION_SELECT_NO       = "no"
)

var (
	BUTTON_SELECT_EMPLOYEE = &tg.KeyboardButtonCallback{
		Text: "Employee",
		Data: []byte(ACTION_SELECT_EMPLOYEE),
	}
	BUTTON_SELECT_EMPLOYER = &tg.KeyboardButtonCallback{
		Text: "Employer",
		Data: []byte(ACTION_SELECT_EMPLOYER),
	}
	BUTTON_SELECT_YES = &tg.KeyboardButtonCallback{
		Text: "YES",
		Data: []byte(ACTION_SELECT_YES),
	}
	BUTTON_SELECT_NO = &tg.KeyboardButtonCallback{
		Text: "NO",
		Data: []byte(ACTION_SELECT_NO),
	}
)
