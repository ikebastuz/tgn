package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/td/telegram"

	"github.com/ikebastuz/tgn/router"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	config, err := getConfig()

	if err != nil {
		log.Fatalf("ERROR: %s", err)
		return
	}

	// Create a new client
	client := telegram.NewClient(int(config.APP_ID), config.APP_HASH, telegram.Options{
		// Add middleware for handling flood waits
		Middlewares: []telegram.Middleware{
			floodwait.NewSimpleWaiter(),
		},
	})

	log.Println("INFO: Starting bot with token:", config.TOKEN[:10]+"...")

	// Add healthcheck handler
	http.HandleFunc("/health", router.HandleHealthCheck)
	// Add webhook handler
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		router.HandleWebhook(ctx, client, w, r)
	})

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("INFO: Starting HTTP server on port %s", config.PORT)
		if err := http.ListenAndServe(":"+config.PORT, nil); err != nil {
			log.Fatalf("ERROR: HTTP server failed: %v", err)
		}
	}()

	// Run the client
	if err := client.Run(ctx, func(ctx context.Context) error {
		// Authenticate as a bot
		log.Println("INFO: Authenticating bot...")
		if _, err := client.Auth().Bot(ctx, config.TOKEN); err != nil {
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
