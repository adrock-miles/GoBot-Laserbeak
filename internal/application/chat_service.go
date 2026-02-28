package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/adrock-miles/go-laserbeak/internal/domain/bot"
	"github.com/adrock-miles/go-laserbeak/internal/domain/conversation"
)

// ChatService is the application service that orchestrates text-based conversations.
type ChatService struct {
	repo         conversation.Repository
	llm          bot.LLMService
	systemPrompt string
	maxHistory   int
}

// NewChatService creates a new ChatService.
func NewChatService(
	repo conversation.Repository,
	llm bot.LLMService,
	systemPrompt string,
	maxHistory int,
) *ChatService {
	return &ChatService{
		repo:         repo,
		llm:          llm,
		systemPrompt: systemPrompt,
		maxHistory:   maxHistory,
	}
}

// HandleMessage processes a user message and returns the LLM response.
func (s *ChatService) HandleMessage(ctx context.Context, channelID, userID, content string) (string, error) {
	// Handle special commands
	if strings.TrimSpace(content) == "/clear" {
		s.repo.Delete(channelID)
		return "", nil
	}

	conv := s.getOrCreateConversation(channelID)

	conv.AddMessage(conversation.NewMessage(conversation.RoleUser, content))

	llmMessages := toLLMMessages(conv.AllMessages())

	reply, err := s.llm.ChatCompletion(ctx, llmMessages)
	if err != nil {
		return "", fmt.Errorf("LLM completion: %w", err)
	}

	conv.AddMessage(conversation.NewMessage(conversation.RoleAssistant, reply))
	s.repo.Save(conv)

	return reply, nil
}

func (s *ChatService) getOrCreateConversation(channelID string) *conversation.Conversation {
	conv, found := s.repo.FindByChannel(channelID)
	if !found {
		conv = conversation.NewConversation(channelID, s.systemPrompt, s.maxHistory)
		s.repo.Save(conv)
	}
	return conv
}

func toLLMMessages(msgs []conversation.Message) []bot.LLMMessage {
	result := make([]bot.LLMMessage, len(msgs))
	for i, m := range msgs {
		result[i] = bot.LLMMessage{
			Role:    string(m.Role),
			Content: m.Content,
		}
	}
	return result
}
