---
sidebar_position: 6
---

# Discord Setup

This guide walks through creating a Discord bot and inviting it to your server.

## 1. Create an application

1. Go to the [Discord Developer Portal](https://discord.com/developers/applications)
2. Click **New Application**
3. Give it a name (e.g. "Laserbeak") and click **Create**

## 2. Create a bot user

1. In your application, go to the **Bot** tab
2. Click **Add Bot** (if not already created)
3. Under **Token**, click **Copy** — this is your `discord.token`

:::caution
Never commit your bot token to version control. Use environment variables or a config file that is gitignored.
:::

## 3. Enable intents

Under **Privileged Gateway Intents**, enable:

- **Message Content Intent** — required for reading text commands
- **Server Members Intent** — optional, for member-related features

## 4. Set up OAuth2 and invite

1. Go to the **OAuth2** tab
2. Under **Scopes**, select:
   - `bot`
   - `applications.commands`
3. Under **Bot Permissions**, select:
   - Send Messages
   - Read Message History
   - Connect (voice)
   - Speak (voice)
   - Use Voice Activity
4. Copy the generated URL and open it in your browser
5. Select your server and click **Authorize**

## 5. Get channel and guild IDs

To get Discord IDs, enable **Developer Mode** in Discord settings (User Settings → Advanced → Developer Mode). Then right-click on any server, channel, or user to copy their ID.

You'll need:

| ID | Config key | How to get it |
|----|-----------|---------------|
| Guild (server) ID | `discord.guildid` | Right-click server name → Copy Server ID |
| Voice channel ID | `discord.voicechannelid` | Right-click voice channel → Copy Channel ID |
| Text channel ID | `discord.textchannelid` | Right-click text channel → Copy Channel ID |
