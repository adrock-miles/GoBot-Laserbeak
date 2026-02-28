package playoptions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/adrock-miles/GoBot-Laserbeak/internal/domain/bot"
)

// Client fetches play options from an external bot's HTTP API and caches them.
type Client struct {
	apiURL   string
	cacheTTL time.Duration
	client   *http.Client

	mu        sync.RWMutex
	cache     []bot.PlayOption
	cacheTime time.Time
	stopCh    chan struct{}
}

// NewClient creates a new play options client with caching.
func NewClient(apiURL string, cacheTTL time.Duration) *Client {
	return &Client{
		apiURL:   apiURL,
		cacheTTL: cacheTTL,
		client:   &http.Client{Timeout: 10 * time.Second},
		stopCh:   make(chan struct{}),
	}
}

// Start begins background cache refresh. Call Stop to clean up.
func (c *Client) Start() {
	// Initial fetch
	if err := c.refresh(); err != nil {
		log.Printf("initial play options fetch failed (url=%s): %v", c.apiURL, err)
	} else {
		log.Printf("loaded %d play options from %s", len(c.cache), c.apiURL)
	}

	go c.refreshLoop()
}

// Stop ends the background refresh loop.
func (c *Client) Stop() {
	close(c.stopCh)
}

// GetOptions returns the cached list of play options, refreshing if stale.
func (c *Client) GetOptions(ctx context.Context) ([]bot.PlayOption, error) {
	c.mu.RLock()
	options := c.cache
	cacheTime := c.cacheTime
	c.mu.RUnlock()

	if time.Since(cacheTime) > c.cacheTTL {
		if err := c.refresh(); err != nil {
			// Return stale cache if available
			if options != nil {
				log.Printf("play options refresh failed (url=%s), using stale cache (%d items, age %s): %v",
					c.apiURL, len(options), time.Since(cacheTime).Truncate(time.Second), err)
				return options, nil
			}
			return nil, fmt.Errorf("play options fetch failed (url=%s) and no cached data available: %w", c.apiURL, err)
		}
		c.mu.RLock()
		options = c.cache
		c.mu.RUnlock()
	}

	return options, nil
}

func (c *Client) refreshLoop() {
	ticker := time.NewTicker(c.cacheTTL)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			if err := c.refresh(); err != nil {
				log.Printf("play options refresh failed (url=%s): %v", c.apiURL, err)
			} else {
				c.mu.RLock()
				count := len(c.cache)
				c.mu.RUnlock()
				log.Printf("refreshed play options: %d items", count)
			}
		}
	}
}

func (c *Client) refresh() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request to %s failed: %w", c.apiURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("play options API returned status %d from %s: %s", resp.StatusCode, c.apiURL, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	// Try parsing as array of strings first, then as array of objects with "name" field
	var options []bot.PlayOption

	var names []string
	if err := json.Unmarshal(body, &names); err == nil {
		for _, name := range names {
			options = append(options, bot.PlayOption{Name: name})
		}
	} else {
		if err := json.Unmarshal(body, &options); err != nil {
			return fmt.Errorf("parse play options: %w", err)
		}
	}

	c.mu.Lock()
	c.cache = options
	c.cacheTime = time.Now()
	c.mu.Unlock()

	return nil
}
