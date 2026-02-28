package playoptions

import (
	"context"
	"log"

	"github.com/adrock-miles/GoBot-Laserbeak/internal/domain/bot"
)

// Composite merges play options from multiple sources into a single PlayOptionsService.
type Composite struct {
	sources []bot.PlayOptionsService
}

// NewComposite creates a Composite from one or more PlayOptionsService sources.
func NewComposite(sources ...bot.PlayOptionsService) *Composite {
	return &Composite{sources: sources}
}

// GetOptions returns the merged options from all sources.
// Errors from individual sources are logged but don't fail the whole call.
func (c *Composite) GetOptions(ctx context.Context) ([]bot.PlayOption, error) {
	var all []bot.PlayOption
	for _, src := range c.sources {
		opts, err := src.GetOptions(ctx)
		if err != nil {
			log.Printf("play options source error: %v", err)
			continue
		}
		all = append(all, opts...)
	}
	return all, nil
}
