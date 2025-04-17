package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/td/telegram"

	"github.com/ikebastuz/tgn/router"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// TODO: extract ENV handling
	appIdString := os.Getenv("BOT_APP_ID")
	if appIdString == "" {
		log.Fatal("ERROR: BOT_APP_ID environment variable is required")
	}
	appId, err := strconv.ParseInt(appIdString, 10, 64)
	if err != nil {
		log.Fatal("ERROR: BOT_APP_ID is not a valid number")
	}

	appHash := os.Getenv("BOT_APP_HASH")
	if appHash == "" {
		log.Fatal("ERROR: BOT_APP_HASH environment variable is required")
	}

	log.Println("INFO: Initializing Telegram bot client...")

	botPort := os.Getenv("BOT_PORT")
	if botPort == "" {
		log.Fatal("ERROR: BOT_PORT environment variable is required")
	}
	log.Printf("INFO: Server will run on port %s", botPort)

	// Create a new client
	client := telegram.NewClient(int(appId), appHash, telegram.Options{
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

	// Add healthcheck handler
	http.HandleFunc("/health", router.HandleHealthCheck)
	// Add webhook handler
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		router.HandleWebhook(ctx, client, w, r)
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
