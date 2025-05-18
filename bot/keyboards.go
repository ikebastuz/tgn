package bot

import (
	"github.com/gotd/td/tg"

	"github.com/ikebastuz/tgn/bot/actions"
)

var KEYBOARD_SELECT_YOUR_ROLE = &tg.ReplyInlineMarkup{
	Rows: []tg.KeyboardButtonRow{
		{
			Buttons: []tg.KeyboardButtonClass{
				actions.BUTTON_SELECT_EMPLOYEE,
				actions.BUTTON_SELECT_EMPLOYER,
			},
		},
	},
}
var KEYBOARD_SELECT_YES_NO = &tg.ReplyInlineMarkup{
	Rows: []tg.KeyboardButtonRow{
		{
			Buttons: []tg.KeyboardButtonClass{
				actions.BUTTON_SELECT_YES,
				actions.BUTTON_SELECT_NO,
			},
		},
	},
}
