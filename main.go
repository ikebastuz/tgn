package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"

	"github.com/ikebastuz/tgn/router"
)

type TelegramUpdate struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID int64 `json:"id"`
		} `json:"from"`
		Text string `json:"text"`
	} `json:"message"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("INFO: Initializing Telegram bot client...")

	botPort := os.Getenv("BOT_PORT")
	if botPort == "" {
		log.Fatal("ERROR: BOT_PORT environment variable is required")
	}
	log.Printf("INFO: Server will run on port %s", botPort)

	// Create a new client
	client := telegram.NewClient(telegram.TestAppID, telegram.TestAppHash, telegram.Options{
		// Add middleware for handling flood waits
		Middlewares: []telegram.Middleware{
			floodwait.NewSimpleWaiter(),
		},
	})

	// Create a bot authenticator
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("ERROR: BOT_TOKEN environment variable is required")
	}

	log.Println("INFO: Starting bot with token:", botToken[:10]+"...")

	// Create HTTP server
	http.HandleFunc("/health", router.HandleHealthCheck)

	// Add webhook handler
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
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

		// Get user info to construct proper InputPeer
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
	})

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("INFO: Starting HTTP server on port %s", botPort)
		if err := http.ListenAndServe(":"+botPort, nil); err != nil {
			log.Fatalf("ERROR: HTTP server failed: %v", err)
		}
	}()

	// Run the client
	if err := client.Run(ctx, func(ctx context.Context) error {
		// Authenticate as a bot
		log.Println("INFO: Authenticating bot...")
		if _, err := client.Auth().Bot(ctx, botToken); err != nil {
			log.Printf("ERROR: Authentication failed: %v", err)
			return err
		}

		// Get the current bot info
		log.Println("INFO: Fetching bot self info...")
		me, err := client.Self(ctx)
		if err != nil {
			log.Printf("ERROR: Failed to get bot info: %v", err)
			return err
		}

		log.Printf("INFO: Bot successfully logged in as @%s", me.Username)

		// Start receiving updates
		log.Println("INFO: Starting to receive updates...")
		return telegram.RunUntilCanceled(ctx, client)
	}); err != nil {
		log.Fatalf("ERROR: Failed to run client: %v", err)
	}
}
