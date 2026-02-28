---
sidebar_position: 1
---

# Docker

Laserbeak includes a multi-stage Dockerfile and a docker-compose file for easy deployment.

## Docker Compose (recommended)

1. Create a `.env` file with your tokens:

```bash
cp .env.example .env
# Edit .env with your values
```

2. Start the bot:

```bash
make docker-up
```

Or without Make:

```bash
docker compose up -d
```

3. View logs:

```bash
docker compose logs -f laserbeak
```

4. Stop the bot:

```bash
make docker-down
```

## Dockerfile

The Dockerfile uses a multi-stage build:

1. **Build stage** (`golang:1.24-bookworm`) — installs libopus, downloads Go dependencies, compiles the binary
2. **Runtime stage** (`debian:bookworm-slim`) — minimal image with only the binary and runtime libraries

```dockerfile
# Build stage
FROM golang:1.24-bookworm AS builder
RUN apt-get update && apt-get install -y libopus-dev libopusfile-dev pkg-config gcc
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /bin/laserbeak .

# Runtime stage
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y libopus0 libopusfile0 ca-certificates
COPY --from=builder /bin/laserbeak /usr/local/bin/laserbeak
CMD ["laserbeak", "serve"]
```

## docker-compose.yml

```yaml
services:
  laserbeak:
    build: .
    env_file: .env
    restart: on-failure
```

Configuration is passed via the `.env` file using `LASERBEAK_` prefixed environment variables. See [Configuration](../getting-started/configuration) for all available settings.
