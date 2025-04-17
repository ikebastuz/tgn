package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gotd/td/telegram"
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

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandleWebhook(ctx context.Context, client *telegram.Client, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("INFO: Raw body:\n%v", string(body))

	keyboard := &tg.ReplyInlineMarkup{
		Rows: []tg.KeyboardButtonRow{
			{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonCallback{
						Text: "Button 1",
						Data: []byte("button1"),
					},
					&tg.KeyboardButtonCallback{
						Text: "Button 2",
						Data: []byte("button2"),
					},
				},
			},
		},
	}
	var update TelegramUpdate
	if err := json.Unmarshal(body, &update); err != nil {
		log.Printf("ERROR: Failed to parse update: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("INFO: Received webhook update:\n%+v", update)

	if update.CallbackQuery.ID != "" {
		userId := update.CallbackQuery.From.ID
		response := fmt.Sprintf("Received %s", update.CallbackQuery.Data)

		// Example of sending notification (may be good for an error)
		// queryId, err := strconv.ParseInt(update.CallbackQuery.ID, 10, 64)
		// if err != nil {
		// 	log.Printf("ERROR: Failed to parse query ID: %v", err)
		// 	http.Error(w, "Bad request", http.StatusBadRequest)
		// 	return
		// }
		// _, err = client.API().MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
		// 	QueryID: queryId,
		// 	Message: response,
		// })
		// if err != nil {
		// 	log.Printf("ERROR: Failed to send message: %v", err)
		// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
		// 	return
		// }

		btn_a_value := rand.Intn(100)
		btn_b_value := rand.Intn(100)
		keyboard := &tg.ReplyInlineMarkup{
			Rows: []tg.KeyboardButtonRow{
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonCallback{
							Text: fmt.Sprintf("Button %d", btn_a_value),
							Data: []byte(fmt.Sprintf("%d", btn_a_value)),
						},
						&tg.KeyboardButtonCallback{
							Text: fmt.Sprintf("Button %d", btn_b_value),
							Data: []byte(fmt.Sprintf("%d", btn_b_value)),
						},
					},
				},
			},
		}
		_, err = client.API().MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
			ID:          update.CallbackQuery.Message.MessageID,
			Peer:        &tg.InputPeerUser{UserID: userId},
			Message:     response,
			ReplyMarkup: keyboard,
		})
		if err != nil {
			log.Printf("ERROR: Failed to send message: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		userId := update.Message.From.ID
		response := "Your Telegram ID is: " + strconv.FormatInt(userId, 10)

		// Send response back to user
		log.Printf("INFO: Sending response: %s", response)

		_, err = client.API().MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
			Peer:        &tg.InputPeerUser{UserID: update.Message.From.ID},
			Message:     response,
			RandomID:    rand.Int63(),
			ReplyMarkup: keyboard,
		})
		if err != nil {
			log.Printf("ERROR: Failed to send message: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
