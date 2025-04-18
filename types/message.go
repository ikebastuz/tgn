package types

type From struct {
	USERNAME string `json:"username"`
	IS_BOT   bool   `json:"is_bot"`
	ID       int64  `json:"id"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      From   `json:"from"`
	Text      string `json:"text"`
}

type TelegramUpdate struct {
	UpdateID      int     `json:"update_id"`
	Message       Message `json:"message"`
	CallbackQuery struct {
		ID      string  `json:"id"`
		From    From    `json:"from"`
		Message Message `json:"message"`
		Data    string  `json:"data"`
	} `json:"callback_query"`
}
