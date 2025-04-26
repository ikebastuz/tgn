package types

import (
	"github.com/gotd/td/tg"
)

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

type CallbackQuery struct {
	ID      string  `json:"id"`
	From    From    `json:"from"`
	Message Message `json:"message"`
	Data    string  `json:"data"`
}

type TelegramUpdate struct {
	UpdateID      int           `json:"update_id"`
	Message       Message       `json:"message"`
	CallbackQuery CallbackQuery `json:"callback_query"`
}

type ReplyMessage struct {
	MessageID   int
	Message     string
	ReplyMarkup tg.ReplyMarkupClass
}

type ReplyDTO struct {
	UserId    int64
	Messages  []ReplyMessage
	NextState *DialogState
}
