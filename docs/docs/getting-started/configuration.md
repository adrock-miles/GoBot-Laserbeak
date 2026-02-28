---
sidebar_position: 3
---

# Configuration

Laserbeak uses a three-tier configuration system. Settings are resolved in this order (highest precedence first):

1. **CLI flags** (e.g. `--discord-token`)
2. **Environment variables** (prefixed with `LASERBEAK_`, e.g. `LASERBEAK_DISCORD_TOKEN`)
3. **Config file** (`config.yaml`)

The config file is loaded from the first path found:
- `./config.yaml`
- `~/.gobot-laserbeak/config.yaml`
- `/etc/gobot-laserbeak/config.yaml`

## Getting started

```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your tokens and channel IDs
```

## All settings

| Setting | CLI Flag | Env Var | Default | Description |
|---------|----------|---------|---------|-------------|
| `discord.token` | `--discord-token` | `LASERBEAK_DISCORD_TOKEN` | — | Discord bot token **(required)** |
| `discord.commandprefix` | `--command-prefix` | `LASERBEAK_DISCORD_COMMANDPREFIX` | `!laser ` | Bot command prefix |
| `discord.guildid` | `--guild-id` | `LASERBEAK_DISCORD_GUILDID` | — | Guild ID for auto-join |
| `discord.voicechannelid` | `--voice-channel-id` | `LASERBEAK_DISCORD_VOICECHANNELID` | — | Voice channel to auto-join |
| `discord.textchannelid` | `--text-channel-id` | `LASERBEAK_DISCORD_TEXTCHANNELID` | — | Text channel for voice command output |
| `llm.apikey` | `--llm-api-key` | `LASERBEAK_LLM_APIKEY` | — | LLM API key **(required)** |
| `llm.baseurl` | `--llm-base-url` | `LASERBEAK_LLM_BASEURL` | `https://api.openai.com/v1` | LLM API base URL |
| `llm.model` | `--llm-model` | `LASERBEAK_LLM_MODEL` | `gpt-4` | LLM model name |
| `stt.apikey` | `--stt-api-key` | `LASERBEAK_STT_APIKEY` | — | STT API key (enables voice) |
| `stt.baseurl` | — | `LASERBEAK_STT_BASEURL` | `https://api.openai.com/v1` | STT API base URL |
| `stt.model` | — | `LASERBEAK_STT_MODEL` | `whisper-1` | STT model name |
| `bot.systemprompt` | — | `LASERBEAK_BOT_SYSTEMPROMPT` | *(built-in)* | System prompt for LLM |
| `bot.maxhistory` | — | `LASERBEAK_BOT_MAXHISTORY` | `50` | Max conversation history per channel |
| `bot.wakephrase` | `--wake-phrase` | `LASERBEAK_BOT_WAKEPHRASE` | `laser` | Wake phrase for voice commands |
| `playoptions.apiurl` | `--play-options-url` | `LASERBEAK_PLAYOPTIONS_APIURL` | — | URL to fetch play options |
| `playoptions.cachettl` | `--play-options-cache-ttl` | `LASERBEAK_PLAYOPTIONS_CACHETTL` | `5m` | Cache TTL for play options |

## Example config file

```yaml
discord:
  token: "YOUR_DISCORD_BOT_TOKEN"
  commandprefix: "!laser "
  guildid: ""
  voicechannelid: ""
  textchannelid: ""

llm:
  apikey: "YOUR_OPENAI_API_KEY"
  baseurl: "https://api.openai.com/v1"
  model: "gpt-4"

stt:
  apikey: "YOUR_OPENAI_API_KEY"
  baseurl: "https://api.openai.com/v1"
  model: "whisper-1"

bot:
  systemprompt: "You are Laserbeak, a helpful Discord assistant."
  maxhistory: 50
  wakephrase: "laser"

playoptions:
  apiurl: ""
  cachettl: "5m"
```

## Example `.env` file

```bash
LASERBEAK_DISCORD_TOKEN=your-discord-token
LASERBEAK_LLM_APIKEY=your-openai-key
LASERBEAK_STT_APIKEY=your-openai-key
LASERBEAK_DISCORD_GUILDID=123456789
LASERBEAK_DISCORD_VOICECHANNELID=123456789
LASERBEAK_DISCORD_TEXTCHANNELID=123456789
```
