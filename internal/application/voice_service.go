package application

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/adrock-miles/GoBot-Laserbeak/internal/domain/bot"
)

// VoiceCommand represents a parsed voice command result.
type VoiceCommand struct {
	// Text is the message to send to the output text channel.
	Text string
}

// VoiceService handles voice-to-text-to-command pipeline.
// It transcribes audio, checks for the wake phrase, and parses voice commands.
// For "play" commands, it uses the LLM to match against available options.
type VoiceService struct {
	stt         bot.STTService
	llm         bot.LLMService
	playOptions bot.PlayOptionsService
	wakePhrase  string
}

// NewVoiceService creates a new VoiceService.
// playOptions and llm may be nil — if so, play commands pass through the raw transcription.
func NewVoiceService(stt bot.STTService, wakePhrase string, llm bot.LLMService, playOptions bot.PlayOptionsService) *VoiceService {
	return &VoiceService{
		stt:         stt,
		llm:         llm,
		playOptions: playOptions,
		wakePhrase:  strings.ToLower(wakePhrase),
	}
}

// HandleVoice transcribes audio and parses voice commands.
// Returns the command text to send to chat, or empty string if no valid command.
func (s *VoiceService) HandleVoice(ctx context.Context, channelID, userID string, audioWAV []byte) (string, error) {
	text, err := s.stt.Transcribe(ctx, audioWAV)
	if err != nil {
		return "", fmt.Errorf("transcribe audio: %w", err)
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return "", nil
	}

	log.Printf("voice transcription from user %s: %s", userID, text)

	cmd, ok := s.parseCommand(ctx, text)
	if !ok {
		return "", nil
	}

	log.Printf("voice command from user %s: %s", userID, cmd.Text)
	return cmd.Text, nil
}

// parseCommand checks if the transcription starts with the wake phrase
// and parses the subsequent command.
func (s *VoiceService) parseCommand(ctx context.Context, transcription string) (VoiceCommand, bool) {
	lower := strings.ToLower(transcription)

	// Normalize common alternate spellings (e.g. "lazer" → "laser")
	normalized := strings.NewReplacer("lazer", "laser").Replace(lower)

	// Check for wake phrase
	if !strings.HasPrefix(normalized, s.wakePhrase) {
		return VoiceCommand{}, false
	}

	// Extract everything after the wake phrase
	rest := strings.TrimSpace(normalized[len(s.wakePhrase):])
	restLower := strings.ToLower(rest)

	// Strip punctuation for command matching (STT may transcribe "Stop!" or "stop.")
	stripped := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' {
			return r
		}
		return -1
	}, restLower)
	stripped = strings.TrimSpace(stripped)

	switch {
	case strings.HasPrefix(stripped, "stop"):
		return VoiceCommand{Text: "!stop"}, true

	case strings.HasPrefix(restLower, "play"):
		query := strings.TrimSpace(rest[len("play"):])
		if query == "" {
			return VoiceCommand{}, false
		}
		if strings.Contains(strings.ToLower(query), "random") {
			return VoiceCommand{Text: "!pr"}, true
		}
		matched := s.matchPlayQuery(ctx, query)
		return VoiceCommand{Text: "!play " + matched}, true
	}

	return VoiceCommand{}, false
}

// matchPlayQuery tries to match a spoken query against the available play options
// using the LLM. Falls back to the raw query if matching is unavailable.
func (s *VoiceService) matchPlayQuery(ctx context.Context, query string) string {
	if s.playOptions == nil || s.llm == nil {
		return query
	}

	options, err := s.playOptions.GetOptions(ctx)
	if err != nil {
		log.Printf("failed to get play options for matching: %v", err)
		return query
	}

	if len(options) == 0 {
		return query
	}

	// Build the options list for the LLM prompt
	var optionNames []string
	for _, opt := range options {
		optionNames = append(optionNames, opt.Name)
	}
	optionsList := strings.Join(optionNames, "\n")

	prompt := fmt.Sprintf(
		"The user said: %q\n\n"+
			"Available options:\n%s\n\n"+
			"Which option best matches what the user asked for? "+
			"Reply with ONLY the exact option name, nothing else. "+
			"If nothing matches, reply with the user's original query exactly as given.",
		query, optionsList,
	)

	messages := []bot.LLMMessage{
		{Role: "system", Content: "You are a matching assistant. Given a spoken query and a list of available options, pick the best match. Reply with only the option name, no explanation."},
		{Role: "user", Content: prompt},
	}

	result, err := s.llm.ChatCompletion(ctx, messages)
	if err != nil {
		log.Printf("LLM matching failed, using raw query: %v", err)
		return query
	}

	result = strings.TrimSpace(result)
	if result == "" {
		return query
	}

	log.Printf("LLM matched %q -> %q", query, result)
	return result
}
