package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrock-miles/GoBot-Laserbeak/internal/application"
	"github.com/adrock-miles/GoBot-Laserbeak/internal/config"
	"github.com/adrock-miles/GoBot-Laserbeak/internal/domain/bot"
	"github.com/adrock-miles/GoBot-Laserbeak/internal/infrastructure/discord"
	"github.com/adrock-miles/GoBot-Laserbeak/internal/infrastructure/llm"
	"github.com/adrock-miles/GoBot-Laserbeak/internal/infrastructure/persistence"
	"github.com/adrock-miles/GoBot-Laserbeak/internal/infrastructure/playoptions"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Laserbeak Discord bot",
	RunE:  runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Infrastructure
	convRepo := persistence.NewInMemoryConversationRepo()
	llmClient := llm.NewOpenAIClient(cfg.LLM.APIKey, cfg.LLM.BaseURL, cfg.LLM.Model)

	// Application services
	chatService := application.NewChatService(
		convRepo,
		llmClient,
		cfg.Bot.SystemPrompt,
		cfg.Bot.MaxHistory,
	)

	// Discord bot
	botCfg := discord.BotConfig{
		Token:          cfg.Discord.Token,
		CommandPrefix:  cfg.Discord.CommandPrefix,
		GuildID:        cfg.Discord.GuildID,
		VoiceChannelID: cfg.Discord.VoiceChannelID,
		TextChannelID:  cfg.Discord.TextChannelID,
	}

	discordBot, err := discord.NewBot(botCfg)
	if err != nil {
		return fmt.Errorf("create bot: %w", err)
	}

	discordBot.SetChatHandler(chatService.HandleMessage)

	// Set up voice service if STT API key is provided
	if cfg.STT.APIKey != "" {
		sttClient := llm.NewSTTClient(cfg.STT.APIKey, cfg.STT.BaseURL, cfg.STT.Model)

		// Build play options sources (local file + optional API)
		var playOptsSources []bot.PlayOptionsService
		playOptsSources = append(playOptsSources, playoptions.NewFileSource("play_options.json"))

		var playOptsClient *playoptions.Client
		if cfg.PlayOptions.APIURL != "" {
			playOptsClient = playoptions.NewClient(cfg.PlayOptions.APIURL, cfg.PlayOptions.CacheTTL)
			playOptsClient.Start()
			playOptsSources = append(playOptsSources, playOptsClient)
			log.Printf("Play options matching enabled (API: %s, cache TTL: %s)",
				cfg.PlayOptions.APIURL, cfg.PlayOptions.CacheTTL)
		}

		playOpts := playoptions.NewComposite(playOptsSources...)
		voiceService := application.NewVoiceService(sttClient, cfg.Bot.WakePhrase, llmClient, playOpts)
		discordBot.SetVoiceHandler(voiceService.HandleVoice)
		log.Printf("Voice commands enabled (wake phrase: %q)", cfg.Bot.WakePhrase)

		if playOptsClient != nil {
			defer playOptsClient.Stop()
		}
	} else {
		log.Println("Voice commands disabled (no STT API key configured)")
	}

	if err := discordBot.Start(); err != nil {
		return fmt.Errorf("start bot: %w", err)
	}
	defer discordBot.Stop()

	log.Printf("Laserbeak is running. Command prefix: %s", cfg.Discord.CommandPrefix)
	if cfg.Discord.TextChannelID != "" {
		log.Printf("Voice command output channel: %s", cfg.Discord.TextChannelID)
	}
	log.Println("Press Ctrl+C to exit")

	// Wait for shutdown signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	log.Println("Shutting down...")
	return nil
}
