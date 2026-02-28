# GoBot-Laserbeak

A Discord LLM bot built in Go that listens to voice channels and responds with text. Built with Cobra, Viper, and Domain Driven Design.

## Features

- **Text Chat**: Respond to text commands with LLM-powered replies
- **Voice Commands**: Listen in voice channels for wake-phrase-activated commands
- **Wake Phrase**: Say "hey m'bot" followed by a command (configurable)
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
   go build -o laserbeak .
   ```

3. **Run**:
   ```bash
   ./laserbeak serve
   ```

Or with environment variables:
```bash
LASERBEAK_DISCORD_TOKEN=... LASERBEAK_LLM_APIKEY=... ./laserbeak serve
```

## Text Commands

| Command | Description |
|---------|-------------|
| `!laser <message>` | Chat with the LLM |
| `!laser join` | Join your voice channel and start listening |
| `!laser leave` | Leave voice channel |
| `!laser clear` | Clear conversation history |
| `!laser help` | Show available commands |

## Voice Commands

Voice commands require the wake phrase (default: "hey m'bot") to be spoken first. The bot transcribes speech via OpenAI Whisper and parses the command.

| Voice Command | Output to Text Chat |
|---------------|---------------------|
| "hey m'bot stop" | `!stop` |
| "hey m'bot play Never Gonna Give You Up" | `!play Never Gonna Give You Up` |

Voice command output is sent to the configured text channel (`discord.textchannelid`).

## Configuration

Configuration is loaded from (in order of precedence):
1. CLI flags (`--discord-token`, `--guild-id`, `--text-channel-id`, etc.)
2. Environment variables (`LASERBEAK_DISCORD_TOKEN`, etc.)
3. Config file (`config.yaml`)

### Key Settings

| Setting | Env Var | Description |
|---------|---------|-------------|
| `discord.token` | `LASERBEAK_DISCORD_TOKEN` | Discord bot token (required) |
| `discord.guildid` | `LASERBEAK_DISCORD_GUILDID` | Guild ID for auto-join |
| `discord.voicechannelid` | `LASERBEAK_DISCORD_VOICECHANNELID` | Voice channel to auto-join |
| `discord.textchannelid` | `LASERBEAK_DISCORD_TEXTCHANNELID` | Text channel for voice command output |
| `bot.wakephrase` | `LASERBEAK_BOT_WAKEPHRASE` | Wake phrase (default: "hey m'bot") |
| `llm.apikey` | `LASERBEAK_LLM_APIKEY` | LLM API key (required) |
| `stt.apikey` | `LASERBEAK_STT_APIKEY` | STT API key (enables voice) |

## Prerequisites

- Go 1.21+
- A Discord bot token with Message Content and Voice intents
- An OpenAI API key (or compatible API)
- libopus (for voice decoding): `apt install libopus-dev libopusfile-dev` / `brew install opus opusfile`
