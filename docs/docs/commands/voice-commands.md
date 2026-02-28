---
sidebar_position: 2
---

# Voice Commands

Voice commands require the bot to be in a voice channel (`!laser join`) and an STT API key to be configured.

## How it works

1. The bot listens to all users in the voice channel
2. Audio is collected per-user as Opus frames
3. Silence detection triggers processing of the buffered audio
4. The audio is transcribed via the STT service (OpenAI Whisper)
5. The transcription is checked for the wake phrase
6. If the wake phrase is detected, the command is parsed and executed
7. Output is sent to the configured text channel (`discord.textchannelid`)

## Wake phrase

The default wake phrase is **"laser"**. The bot also accepts common alternate spellings like "lazer". The wake phrase can be changed via the `bot.wakephrase` config setting.

## Available voice commands

| Voice Command | Output |
|---------------|--------|
| "laser stop" | `!stop` |
| "laser play \<query\>" | `!play \<query\>` |

### Play command matching

When `playoptions.apiurl` is configured, the bot fetches a list of available play options from the API. When a user says "laser play \<something\>", the bot uses the LLM to fuzzy-match the spoken query against the available options and outputs the best match.

If no play options API is configured, a local `play_options.json` file is used as a fallback. If neither is available, the raw query is passed through as-is.

The play options list is cached with a configurable TTL (default: 5 minutes).
