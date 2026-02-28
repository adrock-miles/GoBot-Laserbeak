---
sidebar_position: 5
---

# Architecture

Laserbeak follows **Domain-Driven Design (DDD)** with a **port/adapter** (hexagonal) pattern. All private packages live under `internal/`.

## Layer overview

```
cmd/                         # Interface layer — Cobra CLI commands
internal/
├── domain/                  # Domain layer — pure business logic
│   ├── bot/                 # Service port interfaces (LLMService, STTService, PlayOptionsService)
│   └── conversation/        # Conversation aggregate + Message value object
├── application/             # Application layer — use-case orchestration
│   ├── chat_service.go      # Text chat use case
│   └── voice_service.go     # Voice command parsing
├── infrastructure/          # Infrastructure layer — adapter implementations
│   ├── discord/             # Discord bot handler + voice listener
│   ├── llm/                 # OpenAI-compatible LLM + Whisper STT clients
│   ├── audio/               # Opus decoder, PCM-to-WAV encoder
│   ├── persistence/         # In-memory conversation repository
│   └── playoptions/         # HTTP client with background TTL cache
└── config/                  # Viper-based configuration loading
```

## Domain layer

The domain layer contains pure business logic with no external dependencies.

- **`bot/`** — defines service port interfaces: `LLMService`, `STTService`, `PlayOptionsService`
- **`conversation/`** — the `Conversation` aggregate manages message history; `Message` is a value object

## Application layer

Orchestrates domain logic and infrastructure.

- **`ChatService`** — handles text conversations with history management, calls `LLMService`
- **`VoiceService`** — processes transcribed audio into commands: wake phrase detection, stop/play parsing, LLM-powered fuzzy matching against play options

## Infrastructure layer

Adapters that implement domain ports.

- **`discord/`** — Discord bot handler routes messages to services; voice listener collects Opus frames per user with silence detection
- **`llm/`** — OpenAI-compatible chat completions client and Whisper-compatible STT client
- **`audio/`** — decodes Opus frames to PCM, encodes PCM to WAV for STT submission
- **`persistence/`** — in-memory conversation repository guarded by `sync.RWMutex`
- **`playoptions/`** — HTTP client that fetches and caches play options with a configurable TTL

## Data flow

### Text command flow

```
Discord message
  → Bot.handler routes by command prefix
    → ChatService manages conversation history
      → LLMService generates response
        → Response sent back to Discord channel
```

### Voice command flow

```
Voice audio in Discord
  → VoiceListener collects Opus frames per user
    → Silence detection triggers processing
      → Audio decoded (Opus → PCM → WAV)
        → STTService transcribes audio
          → VoiceService checks for wake phrase
            → Command parsed (stop / play <query>)
              → Optional: LLM fuzzy-matches query against play options
                → Output sent to configured text channel
```
