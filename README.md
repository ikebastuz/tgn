# TGN - Telegram Negotiation Bot

A Telegram bot that helps two people negotiate salary anonymously. Each person sets their salary range, and the bot finds a fair number both can agree on.

## How it works

1. Person A starts negotiation, gets connection ID
2. Person B joins using the connection ID
3. Both select role (employer/employee)
4. Both enter their salary range (min/max)
5. Bot calculates overlap and suggests fair salary
6. Both get the same result

## How to start

```bash
# Clone and setup
git clone https://github.com/ikebastuz/tgn.git
cd tgn
go mod download

# Set environment variables
export BOT_APP_ID=your_telegram_app_id
export BOT_APP_HASH=your_telegram_app_hash
export BOT_TOKEN=your_bot_token
export BOT_PORT=8080

# Run
go run cmd/webserver/main.go
```

## How to use

1. Message the bot: `/start`
2. Share the connection ID with the other person
3. Other person messages bot: `/connect <ID>`
4. Both select employer/employee role
5. Both enter salary ranges
6. Get negotiated result

## Go concepts used

- **Interfaces** - State machine with different state types
- **Type switches** - Handling different conversation states
- **Structs** - Configuration, message types, user data
- **Methods** - State machine operations, solver algorithm
- **Goroutines** - HTTP server runs concurrently with Telegram client
- **Channels** - Context cancellation for graceful shutdown
- **Error handling** - Custom error types and validation
- **Environment variables** - Configuration management
- **HTTP handlers** - Webhook, health check, metrics endpoints
- **Testing** - Unit tests for solver algorithm and utilities
- **Modules** - Go mod for dependency management
