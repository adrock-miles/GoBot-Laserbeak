package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "laserbeak",
	Short: "Laserbeak — a Discord LLM bot that listens to voice and chats in text",
	Long: `Laserbeak is a Discord bot built in Go that:
  - Responds to text commands with LLM-powered replies
  - Joins voice channels to listen for voice commands
  - Recognizes "laser" wake phrase for stop/play commands
  - Sends command output as text messages in a configurable Discord channel`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "config file path (default: ./config.yaml)")
	rootCmd.PersistentFlags().String("discord-token", "", "Discord bot token")
	rootCmd.PersistentFlags().String("llm-api-key", "", "LLM API key")
	rootCmd.PersistentFlags().String("llm-base-url", "", "LLM API base URL")
	rootCmd.PersistentFlags().String("llm-model", "", "LLM model name")
	rootCmd.PersistentFlags().String("stt-api-key", "", "STT API key")
	rootCmd.PersistentFlags().String("command-prefix", "", "Bot command prefix")
	rootCmd.PersistentFlags().String("guild-id", "", "Discord guild ID for auto-join")
	rootCmd.PersistentFlags().String("voice-channel-id", "", "Voice channel ID to auto-join")
	rootCmd.PersistentFlags().String("text-channel-id", "", "Text channel ID for voice command output")
	rootCmd.PersistentFlags().String("wake-phrase", "", "Wake phrase for voice commands")
	rootCmd.PersistentFlags().String("play-options-url", "", "URL to fetch play options from")
	rootCmd.PersistentFlags().String("play-options-cache-ttl", "", "Cache TTL for play options (e.g. 5m)")

	viper.BindPFlag("discord.token", rootCmd.PersistentFlags().Lookup("discord-token"))
	viper.BindPFlag("llm.apikey", rootCmd.PersistentFlags().Lookup("llm-api-key"))
	viper.BindPFlag("llm.baseurl", rootCmd.PersistentFlags().Lookup("llm-base-url"))
	viper.BindPFlag("llm.model", rootCmd.PersistentFlags().Lookup("llm-model"))
	viper.BindPFlag("stt.apikey", rootCmd.PersistentFlags().Lookup("stt-api-key"))
	viper.BindPFlag("discord.commandprefix", rootCmd.PersistentFlags().Lookup("command-prefix"))
	viper.BindPFlag("discord.guildid", rootCmd.PersistentFlags().Lookup("guild-id"))
	viper.BindPFlag("discord.voicechannelid", rootCmd.PersistentFlags().Lookup("voice-channel-id"))
	viper.BindPFlag("discord.textchannelid", rootCmd.PersistentFlags().Lookup("text-channel-id"))
	viper.BindPFlag("bot.wakephrase", rootCmd.PersistentFlags().Lookup("wake-phrase"))
	viper.BindPFlag("playoptions.apiurl", rootCmd.PersistentFlags().Lookup("play-options-url"))
	viper.BindPFlag("playoptions.cachettl", rootCmd.PersistentFlags().Lookup("play-options-cache-ttl"))
}
