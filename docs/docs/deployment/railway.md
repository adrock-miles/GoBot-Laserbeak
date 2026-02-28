---
sidebar_position: 2
---

# Railway

Laserbeak can be deployed to [Railway](https://railway.app) with automatic deploys from GitHub.

## railway.toml

The repository includes a `railway.toml` that configures the build and deployment:

```toml
[build]
builder = "DOCKERFILE"
dockerfilePath = "Dockerfile"

[deploy]
startCommand = "laserbeak serve"
restartPolicyType = "ON_FAILURE"
restartPolicyMaxRetries = 5
```

## Setup

1. Create a new project on [Railway](https://railway.app)
2. Connect your GitHub repository
3. Add environment variables in the Railway dashboard:

| Variable | Required | Description |
|----------|----------|-------------|
| `LASERBEAK_DISCORD_TOKEN` | Yes | Discord bot token |
| `LASERBEAK_LLM_APIKEY` | Yes | LLM API key |
| `LASERBEAK_STT_APIKEY` | No | STT API key (for voice) |
| `LASERBEAK_DISCORD_GUILDID` | No | Guild ID for auto-join |
| `LASERBEAK_DISCORD_VOICECHANNELID` | No | Voice channel to auto-join |
| `LASERBEAK_DISCORD_TEXTCHANNELID` | No | Text channel for voice output |

4. Deploy — Railway will build using the Dockerfile and start the bot

## Auto-deploy

Railway automatically redeploys when you push to the connected branch. No additional CI/CD configuration is needed.
