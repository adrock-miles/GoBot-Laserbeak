---
slug: /
sidebar_position: 1
---

# Laserbeak

A Discord LLM bot built in Go that listens to voice channels and responds with text.

## Features

- **Text Chat** — respond to text commands with LLM-powered replies
- **Voice Commands** — listen in voice channels for wake-phrase-activated commands
- **Wake Phrase** — say "laser" followed by a command (configurable)
- **Configurable Channels** — set default voice channel to join and text channel for output
- **Conversation Memory** — per-channel conversation history with configurable limits
- **OpenAI Compatible** — works with any OpenAI-compatible API (OpenAI, Ollama, etc.)

## Quick Links

- [Prerequisites](getting-started/prerequisites) — what you need before you start
- [Installation](getting-started/installation) — build from source
- [Configuration](getting-started/configuration) — all settings explained
- [Text Commands](commands/text-commands) — available chat commands
- [Voice Commands](commands/voice-commands) — voice pipeline and wake phrase
- [Architecture](architecture) — DDD layers and data flow
- [Discord Setup](discord-setup) — creating a Discord bot and inviting it
- [Docker Deployment](deployment/docker) — run with Docker and docker-compose
- [Railway Deployment](deployment/railway) — deploy to Railway
