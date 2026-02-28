---
sidebar_position: 1
---

# Prerequisites

## Required Software

- **Go 1.24+** — [download](https://go.dev/dl/)
- **libopus** — required for voice channel audio decoding

### macOS

```bash
brew install opus opusfile
```

### Linux (Debian/Ubuntu)

```bash
apt install libopus-dev libopusfile-dev
```

## Required Accounts & Tokens

| Token | Purpose | Where to get it |
|-------|---------|-----------------|
| Discord bot token | Connect to Discord | [Discord Developer Portal](https://discord.com/developers/applications) |
| LLM API key | Chat completions | [OpenAI](https://platform.openai.com/api-keys) or any compatible provider |
| STT API key | Voice transcription (optional) | Same as LLM key if using OpenAI |

The LLM API key is required. The STT API key is only needed if you want voice commands.

See [Discord Setup](../discord-setup) for a walkthrough on creating a Discord bot.
