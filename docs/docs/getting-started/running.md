---
sidebar_position: 4
---

# Running

## With Make

```bash
make run
```

This builds the binary and starts the bot with `./laserbeak serve`.

## Directly

```bash
./laserbeak serve
```

## With environment variables

```bash
LASERBEAK_DISCORD_TOKEN=... LASERBEAK_LLM_APIKEY=... ./laserbeak serve
```

## With Docker Compose

```bash
# Create .env with your tokens (see .env.example)
make docker-up
```

Or without Make:

```bash
docker compose up -d
```

See [Docker Deployment](../deployment/docker) for details.

## With a config file

```bash
./laserbeak serve --config /path/to/config.yaml
```
