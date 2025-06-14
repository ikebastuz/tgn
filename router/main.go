package router

import (
	"context"
	"encoding/json"
	"io"

	log "github.com/sirupsen/logrus"

	"net/http"

	"github.com/gotd/td/telegram"

	"github.com/ikebastuz/tgn/bot"
	"github.com/ikebastuz/tgn/types"
)

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Service is healthy"))
}

func HandleWebhook(ctx context.Context, client *telegram.Client, store types.Store, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Warnf("❌ Invalid request method: %s", r.Method)
		http.Error(w, "⚠️ Method not allowed - only POST requests are accepted", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("❌ Request error: Failed to read request body: %v", err)
		http.Error(w, "⚠️ Invalid request format - could not read request body", http.StatusBadRequest)
		return
	}

	// log.Print("INFO: Raw body:\n", string(body))
	// body_debug, _ := json.MarshalIndent(body, "", " ")
	// fmt.Println(string(body_debug))

	// keyboard := &tg.ReplyInlineMarkup{
	// 	Rows: []tg.KeyboardButtonRow{
	// 		{
	// 			Buttons: []tg.KeyboardButtonClass{
	// 				actions.BUTTON_SELECT_EMPLOYEE,
	// 				actions.BUTTON_SELECT_EMPLOYER,
	// 			},
	// 		},
	// 	},
	// }
	var update types.TelegramUpdate
	if err := json.Unmarshal(body, &update); err != nil {
		log.Errorf("❌ Parse error: Failed to parse update: %v", err)
		http.Error(w, "⚠️ Invalid request format - could not parse JSON body", http.StatusBadRequest)
		return
	}

	// log.Printf("INFO: Received webhook update:\n")
	// update_debug, _ := json.MarshalIndent(update, "", " ")
	// fmt.Println(string(update_debug))

	err = bot.HandleMessage(ctx, client, update, store)

	if err != nil {
		log.Errorf("❌ Server error: Failed to handle message: %v", err)
		http.Error(w, "⚠️ Internal server error - please try again later", http.StatusInternalServerError)
		return
	}

	// if update.CallbackQuery.ID != "" {
	// 	userId := update.CallbackQuery.From.ID
	// 	response := fmt.Sprintf("Received %s\n/connect_123\n/connect#123", update.CallbackQuery.Data)
	//
	// 	// Example of sending notification (may be good for an error)
	// 	// queryId, err := strconv.ParseInt(update.CallbackQuery.ID, 10, 64)
	// 	// if err != nil {
	// 	// 	log.Printf("ERROR: Failed to parse query ID: %v", err)
	// 	// 	http.Error(w, "Bad request", http.StatusBadRequest)
	// 	// 	return
	// 	// }
	// 	// _, err = client.API().MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
	// 	// 	QueryID: queryId,
	// 	// 	Message: response,
	// 	// })
	// 	// if err != nil {
	// 	// 	log.Printf("ERROR: Failed to send message: %v", err)
	// 	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	// 	return
	// 	// }
	//
	// 	btn_a_value := rand.Intn(100)
	// 	btn_b_value := rand.Intn(100)
	// 	keyboard := &tg.ReplyInlineMarkup{
	// 		Rows: []tg.KeyboardButtonRow{
	// 			{
	// 				Buttons: []tg.KeyboardButtonClass{
	// 					&tg.KeyboardButtonCallback{
	// 						Text: fmt.Sprintf("Button %d", btn_a_value),
	// 						Data: fmt.Appendf(nil, "%d", btn_a_value),
	// 					},
	// 					&tg.KeyboardButtonCallback{
	// 						Text: fmt.Sprintf("Button %d", btn_b_value),
	// 						Data: fmt.Appendf(nil, "%d", btn_b_value),
	// 					},
	// 				},
	// 			},
	// 		},
	// 	}
	// 	_, err = client.API().MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
	// 		ID:          update.CallbackQuery.Message.MessageID,
	// 		Peer:        &tg.InputPeerUser{UserID: userId},
	// 		Message:     response,
	// 		ReplyMarkup: keyboard,
	// 	})
	// 	if err != nil {
	// 		log.Printf("ERROR: Failed to send message: %v", err)
	// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 		return
	// 	}
	// } else {
	// 	userId := update.Message.From.ID
	// 	response := "Your Telegram ID is: " + strconv.FormatInt(userId, 10)
	//
	// 	// Send response back to user
	// 	log.Printf("INFO: Sending response: %s", response)
	//
	// 	_, err = client.API().MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
	// 		Peer:        &tg.InputPeerUser{UserID: update.Message.From.ID},
	// 		Message:     response,
	// 		RandomID:    rand.Int63(),
	// 		ReplyMarkup: keyboard,
	// 	})
	// 	if err != nil {
	// 		log.Printf("ERROR: Failed to send message: %v", err)
	// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Message processed successfully"))
}
