# go-laserbeak

A Discord LLM bot built in Go that listens to voice channels and responds with text. Built with Cobra, Viper, and Domain Driven Design.

## Documentation

Full documentation is available at **[adrock-miles.github.io/go-laserbeak](https://adrock-miles.github.io/go-laserbeak/)**.

## Features

- **Text Chat**: Respond to text commands with LLM-powered replies
- **Voice Commands**: Listen in voice channels for wake-phrase-activated commands
- **Wake Phrase**: Say "laser" followed by a command (configurable)
- **Configurable Channels**: Set default voice channel to join and text channel for output
- **Conversation Memory**: Per-channel conversation history with configurable limits
- **OpenAI Compatible**: Works with any OpenAI-compatible API (OpenAI, Ollama, etc.)

## Architecture (DDD)

```
internal/
├── domain/              # Domain layer — entities, value objects, interfaces
│   ├── conversation/    # Conversation aggregate (entity, message, repository)
│   └── bot/             # LLM and STT service port interfaces
├── application/         # Application layer — use cases / orchestration
│   ├── chat_service.go  # Text chat use case
│   └── voice_service.go # Voice command parsing (wake phrase + stop/play)
├── infrastructure/      # Infrastructure layer — external adapters
│   ├── discord/         # Discord bot + voice listener
│   ├── llm/             # OpenAI chat + Whisper STT clients
│   ├── audio/           # Opus decoder + WAV encoder
│   └── persistence/     # In-memory conversation repository
└── config/              # Viper configuration loading
cmd/                     # Interface layer — Cobra CLI commands
```

## Adding the Bot to Your Server

A server admin must authorize the bot by visiting the following link:

https://discord.com/oauth2/authorize?client_id=1477087945273118780&scope=bot&permissions=1328823617

This grants the bot the permissions it needs to read and send messages, connect to voice channels, and speak.

## Quick Start

1. **Copy config**:
   ```bash
   cp config.yaml.example config.yaml
   # Edit config.yaml with your tokens and channel IDs
   ```

2. **Build**:
   ```bash
   make build
   ```
   Or directly with Go: `go build -o laserbeak .`

3. **Run**:
   ```bash
   make run
   ```
   Or directly: `./laserbeak serve`

Or with environment variables:
```bash
LASERBEAK_DISCORD_TOKEN=... LASERBEAK_LLM_APIKEY=... ./laserbeak serve
```

## Development

Run `make help` to see all available targets:

```
  help             Show this help
  build            Build the binary
  clean            Remove build artifacts
  run              Build and run the bot
  docker-build     Build Docker image
  docker-up        Start containers in background
  docker-down      Stop containers
  docs             Build the documentation site
  docs-serve       Start local docs dev server
```

## Prerequisites

- Go 1.24+
- A Discord bot token with Message Content and Voice intents
- An OpenAI API key (or compatible API)
- libopus (for voice decoding): `apt install libopus-dev libopusfile-dev` / `brew install opus opusfile`
