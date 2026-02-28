# ---- Build stage ----
FROM golang:1.24-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    libopus-dev libopusfile-dev pkg-config gcc \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o /bin/laserbeak .

# ---- Runtime stage ----
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libopus0 libopusfile0 ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /bin/laserbeak /usr/local/bin/laserbeak

CMD ["laserbeak", "serve"]
