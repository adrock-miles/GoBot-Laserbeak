package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Discord     DiscordConfig
	LLM         LLMConfig
	STT         STTConfig
	Bot         BotConfig
	PlayOptions PlayOptionsConfig
}

// PlayOptionsConfig holds settings for the play options API.
type PlayOptionsConfig struct {
	APIURL   string        // URL to fetch play options from (e.g. http://localhost:8080/options)
	CacheTTL time.Duration // how long to cache the options list
}

// DiscordConfig holds Discord-specific settings.
type DiscordConfig struct {
	Token          string
	CommandPrefix  string
	GuildID        string // guild to operate in (required for auto-join)
	VoiceChannelID string // voice channel to auto-join on startup
	TextChannelID  string // text channel for voice command output
}

// LLMConfig holds LLM API settings.
type LLMConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

// STTConfig holds speech-to-text API settings.
type STTConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

// BotConfig holds general bot behavior settings.
type BotConfig struct {
	SystemPrompt string
	MaxHistory   int
	WakePhrase   string // wake phrase for voice commands (e.g. "hey m'bot")
}

// Load reads configuration from environment variables, config files, and flags.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.gobot-laserbeak")
	viper.AddConfigPath("/etc/gobot-laserbeak")

	// Environment variable mapping
	viper.SetEnvPrefix("LASERBEAK")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Bind each config key to both LASERBEAK_-prefixed and un-prefixed env vars.
	// This allows platforms like Railway to use simpler variable names
	// (e.g. DISCORD_TOKEN instead of LASERBEAK_DISCORD_TOKEN).
	// The LASERBEAK_-prefixed version takes precedence when both are set.
	envBindings := map[string][2]string{
		"discord.token":          {"LASERBEAK_DISCORD_TOKEN", "DISCORD_TOKEN"},
		"discord.commandprefix":  {"LASERBEAK_DISCORD_COMMANDPREFIX", "DISCORD_COMMANDPREFIX"},
		"discord.guildid":        {"LASERBEAK_DISCORD_GUILDID", "DISCORD_GUILDID"},
		"discord.voicechannelid": {"LASERBEAK_DISCORD_VOICECHANNELID", "DISCORD_VOICECHANNELID"},
		"discord.textchannelid":  {"LASERBEAK_DISCORD_TEXTCHANNELID", "DISCORD_TEXTCHANNELID"},
		"llm.apikey":             {"LASERBEAK_LLM_APIKEY", "LLM_APIKEY"},
		"llm.baseurl":            {"LASERBEAK_LLM_BASEURL", "LLM_BASEURL"},
		"llm.model":              {"LASERBEAK_LLM_MODEL", "LLM_MODEL"},
		"stt.apikey":             {"LASERBEAK_STT_APIKEY", "STT_APIKEY"},
		"stt.baseurl":            {"LASERBEAK_STT_BASEURL", "STT_BASEURL"},
		"stt.model":              {"LASERBEAK_STT_MODEL", "STT_MODEL"},
		"bot.systemprompt":       {"LASERBEAK_BOT_SYSTEMPROMPT", "BOT_SYSTEMPROMPT"},
		"bot.maxhistory":         {"LASERBEAK_BOT_MAXHISTORY", "BOT_MAXHISTORY"},
		"bot.wakephrase":         {"LASERBEAK_BOT_WAKEPHRASE", "BOT_WAKEPHRASE"},
		"playoptions.apiurl":     {"LASERBEAK_PLAYOPTIONS_APIURL", "PLAYOPTIONS_APIURL"},
		"playoptions.cachettl":   {"LASERBEAK_PLAYOPTIONS_CACHETTL", "PLAYOPTIONS_CACHETTL"},
	}
	for key, envVars := range envBindings {
		viper.BindEnv(key, envVars[0], envVars[1])
	}

	// Defaults
	viper.SetDefault("discord.commandprefix", "!laser")
	viper.SetDefault("llm.baseurl", "https://api.openai.com/v1")
	viper.SetDefault("llm.model", "gpt-4")
	viper.SetDefault("stt.baseurl", "https://api.openai.com/v1")
	viper.SetDefault("stt.model", "whisper-1")
	viper.SetDefault("bot.systemprompt", "You are Laserbeak, a helpful Discord assistant. Respond concisely and helpfully.")
	viper.SetDefault("bot.maxhistory", 50)
	viper.SetDefault("bot.wakephrase", "hey m'bot")
	viper.SetDefault("playoptions.cachettl", "5m")

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
	}

	cfg := &Config{
		Discord: DiscordConfig{
			Token:          viper.GetString("discord.token"),
			CommandPrefix:  viper.GetString("discord.commandprefix"),
			GuildID:        viper.GetString("discord.guildid"),
			VoiceChannelID: viper.GetString("discord.voicechannelid"),
			TextChannelID:  viper.GetString("discord.textchannelid"),
		},
		LLM: LLMConfig{
			APIKey:  viper.GetString("llm.apikey"),
			BaseURL: viper.GetString("llm.baseurl"),
			Model:   viper.GetString("llm.model"),
		},
		STT: STTConfig{
			APIKey:  viper.GetString("stt.apikey"),
			BaseURL: viper.GetString("stt.baseurl"),
			Model:   viper.GetString("stt.model"),
		},
		Bot: BotConfig{
			SystemPrompt: viper.GetString("bot.systemprompt"),
			MaxHistory:   viper.GetInt("bot.maxhistory"),
			WakePhrase:   viper.GetString("bot.wakephrase"),
		},
	}

	cacheTTL, err := time.ParseDuration(viper.GetString("playoptions.cachettl"))
	if err != nil {
		cacheTTL = 5 * time.Minute
	}
	cfg.PlayOptions = PlayOptionsConfig{
		APIURL:   viper.GetString("playoptions.apiurl"),
		CacheTTL: cacheTTL,
	}

	if cfg.Discord.Token == "" {
		return nil, fmt.Errorf("discord.token is required (set DISCORD_TOKEN or LASERBEAK_DISCORD_TOKEN)")
	}
	if cfg.LLM.APIKey == "" {
		return nil, fmt.Errorf("llm.apikey is required (set LLM_APIKEY or LASERBEAK_LLM_APIKEY)")
	}

	return cfg, nil
}
