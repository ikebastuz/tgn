package bot

import (
	"github.com/gotd/td/tg"

	"github.com/ikebastuz/tgn/actions"
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
