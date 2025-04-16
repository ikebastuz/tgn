package router

import (
	"context"
	"encoding/json"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type TelegramUpdate struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			FIRST_NAME string `json:"first_name"`
			LAST_NAME  string `json:"last_name"`
			USERNAME   string `json:"username"`
			IS_BOT     bool   `json:"is_bot"`
			ID         int64  `json:"id"`
		} `json:"from"`
		Text string `json:"text"`
	} `json:"message"`
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

	log.Printf("INFO: Raw body:", string(body))

	var update TelegramUpdate
	if err := json.Unmarshal(body, &update); err != nil {
		log.Printf("ERROR: Failed to parse update: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("INFO: Received webhook update: %+v", update)

	// Send response back to user
	response := "Your Telegram ID is: " + strconv.FormatInt(update.Message.From.ID, 10)
	log.Printf("INFO: Sending response: %s", response)

	_, err = client.API().MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
		Peer:     &tg.InputPeerUser{UserID: update.Message.From.ID},
		Message:  response,
		RandomID: rand.Int63(),
	})
	if err != nil {
		log.Printf("ERROR: Failed to send message: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
