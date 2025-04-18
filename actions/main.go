package actions

import (
	"github.com/gotd/td/tg"
)

const (
	ACTION_SELECT_EMPLOYEE = "employee"
	ACTION_SELECT_EMPLOYER = "employer"
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
)
