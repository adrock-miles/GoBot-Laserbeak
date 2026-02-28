# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

Requires Go 1.24+ and system libraries for Opus audio decoding:
- macOS: `brew install opus opusfile`
- Linux: `apt install libopus-dev libopusfile-dev`

```bash
go build -o laserbeak .        # Build binary
./laserbeak serve              # Start the bot
./laserbeak version            # Print version
```

No tests, linter config, Makefile, or CI/CD exist yet.

## Configuration

Three-tier precedence (highest to lowest): CLI flags → environment variables (`LASERBEAK_` prefix, e.g. `LASERBEAK_LLM_APIKEY`) → config file (`config.yaml` in `.`, `~/.gobot-laserbeak/`, or `/etc/gobot-laserbeak/`). See `config.yaml.example` and `.env.example` for templates.

Required: `discord.token`, `llm.apikey`. Voice features require `stt.apikey`.

## Architecture

Domain-Driven Design with port/adapter pattern. All private packages live under `internal/`.

**Layers:**
- **`cmd/`** — Cobra CLI commands. `serve.go` wires up all dependencies and starts the bot.
- **`internal/domain/`** — Pure domain: interfaces (ports) in `bot/` (`LLMService`, `STTService`, `PlayOptionsService`), conversation aggregate + message value object in `conversation/`.
- **`internal/application/`** — Use-case orchestration. `ChatService` handles text conversations with history. `VoiceService` processes transcribed audio into commands (wake phrase detection, stop/play parsing, LLM-powered option matching).
- **`internal/infrastructure/`** — Adapters implementing domain ports:
  - `discord/` — Discord bot handler + voice listener (Opus frame collection, per-user audio buffering, silence detection)
  - `llm/` — OpenAI-compatible LLM client and Whisper-compatible STT client
  - `audio/` — Opus decoding, PCM-to-WAV encoding
  - `persistence/` — In-memory conversation repo (sync.RWMutex guarded)
  - `playoptions/` — HTTP client with background TTL cache for external play options API
- **`internal/config/`** — Viper-based config loading with CLI flag binding.

**Key data flow:** Discord message → `Bot.handler` routes by command → `ChatService` manages conversation history and calls `LLMService` → response sent back to Discord channel. Voice path: `VoiceListener` collects Opus frames per user → silence detection triggers processing → `VoiceService` transcribes via `STTService` → wake phrase check → command parsing → optional LLM matching against play options → output sent to configured text channel.
