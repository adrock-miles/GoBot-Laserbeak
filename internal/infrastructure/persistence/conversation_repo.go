package persistence

import (
	"sync"

	"github.com/adrock-miles/go-laserbeak/internal/domain/conversation"
)

// InMemoryConversationRepo implements conversation.Repository with an in-memory store.
type InMemoryConversationRepo struct {
	mu    sync.RWMutex
	store map[string]*conversation.Conversation
}

// NewInMemoryConversationRepo creates a new in-memory conversation repository.
func NewInMemoryConversationRepo() *InMemoryConversationRepo {
	return &InMemoryConversationRepo{
		store: make(map[string]*conversation.Conversation),
	}
}

func (r *InMemoryConversationRepo) FindByChannel(channelID string) (*conversation.Conversation, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	conv, ok := r.store[channelID]
	return conv, ok
}

func (r *InMemoryConversationRepo) Save(conv *conversation.Conversation) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[conv.ChannelID] = conv
}

func (r *InMemoryConversationRepo) Delete(channelID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, channelID)
}
