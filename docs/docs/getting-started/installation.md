---
sidebar_position: 2
---

# Installation

## Clone the repository

```bash
git clone https://github.com/adrock-miles/GoBot-Laserbeak.git
cd GoBot-Laserbeak
```

## Build

Using Make (recommended):

```bash
make build
```

This injects the current git version tag into the binary via `-ldflags`.

Or build directly with Go:

```bash
go build -o laserbeak .
```

## Verify

```bash
./laserbeak version
```

You should see output like:

```
Laserbeak v0.1.0
```
