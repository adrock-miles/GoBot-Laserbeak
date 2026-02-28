---
sidebar_position: 1
---

# Text Commands

All text commands use the configured prefix (default: `!laser`).

| Command | Description |
|---------|-------------|
| `!laser <message>` | Chat with the LLM |
| `!laser join` | Join your voice channel and start listening |
| `!laser leave` | Leave voice channel |
| `!laser clear` | Clear conversation history for the channel |
| `!laser help` | Show available commands |

## Examples

```
!laser What is the capital of France?
!laser join
!laser clear
```

The bot maintains per-channel conversation history, so follow-up questions work naturally. Use `!laser clear` to reset the conversation context.
