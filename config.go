package main

import (
	"errors"
	"log"
	"os"
	"strconv"
)

type Config struct {
	APP_ID   int64
	APP_HASH string
	PORT     string
	TOKEN    string
}

func getConfig() (*Config, error) {
	appIdString := os.Getenv("BOT_APP_ID")
	if appIdString == "" {
		return nil, errors.New("BOT_APP_ID environment variable is required")
	}
	appId, err := strconv.ParseInt(appIdString, 10, 64)
	if err != nil {
		return nil, errors.New("BOT_APP_ID is not a valid number")
	}

	appHash := os.Getenv("BOT_APP_HASH")
	if appHash == "" {
		return nil, errors.New("BOT_APP_HASH environment variable is required")
	}

	log.Println("INFO: Initializing Telegram bot client...")

	botPort := os.Getenv("BOT_PORT")
	if botPort == "" {
		return nil, errors.New("BOT_PORT environment variable is required")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, errors.New("BOT_TOKEN environment variable is required")
	}

	return &Config{
		APP_ID:   appId,
		APP_HASH: appHash,
		PORT:     botPort,
		TOKEN:    botToken,
	}, nil
}
