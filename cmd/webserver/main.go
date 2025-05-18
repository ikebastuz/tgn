package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/td/telegram"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ikebastuz/tgn"
	"github.com/ikebastuz/tgn/bot"
	"github.com/ikebastuz/tgn/metrics"
	"github.com/ikebastuz/tgn/router"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	metrics.InitMetrics()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	config, err := tgn.GetConfig()

	if err != nil {
		log.Fatalf("❌ Configuration error: %s", err)
		return
	}

	// Create a new client
	client := telegram.NewClient(int(config.APP_ID), config.APP_HASH, telegram.Options{
		// Add middleware for handling flood waits
		Middlewares: []telegram.Middleware{
			floodwait.NewSimpleWaiter(),
		},
	})
	store := bot.NewInMemoryStore()

	log.Info("🔑 Starting bot with token:", config.TOKEN[:10]+"...")

	// Add Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	// Add healthcheck handler
	http.HandleFunc("/health", router.HandleHealthCheck)
	// Add webhook handler
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		router.HandleWebhook(ctx, client, store, w, r)
	})

	// Start HTTP server in a goroutine
	go func() {
		log.Infof("🌐 Starting HTTP server on port %s", config.PORT)
		if err := http.ListenAndServe(":"+config.PORT, nil); err != nil {
			log.Fatalf("❌ HTTP server error: %v", err)
		}
	}()

	// Run the client
	if err := client.Run(ctx, func(ctx context.Context) error {
		// Authenticate as a bot
		log.Info("🔐 Authenticating bot...")
		if _, err := client.Auth().Bot(ctx, config.TOKEN); err != nil {
			log.Errorf("❌ Authentication failed: %v", err)
			return err
		}

		// Get the current bot info
		log.Info("ℹ️ Fetching bot information...")
		me, err := client.Self(ctx)
		if err != nil {
			log.Errorf("❌ Failed to get bot info: %v", err)
			return err
		}

		log.Infof("✅ Bot successfully logged in as @%s", me.Username)

		// Start receiving updates
		log.Info("📡 Starting to receive updates...")
		return telegram.RunUntilCanceled(ctx, client)
	}); err != nil {
		log.Fatalf("❌ Client error: %v", err)
	}
}
