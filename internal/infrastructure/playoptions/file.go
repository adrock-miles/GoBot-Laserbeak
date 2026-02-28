package playoptions

import (
	"context"
	"encoding/json"
	"os"

	"github.com/adrock-miles/GoBot-Laserbeak/internal/domain/bot"
)

// FileSource reads play options from a local JSON file on each call.
// Supports both string arrays (["a","b"]) and object arrays ([{"name":"a"}]).
type FileSource struct {
	path string
}

// NewFileSource creates a FileSource that reads from the given path.
func NewFileSource(path string) *FileSource {
	return &FileSource{path: path}
}

// GetOptions reads and parses the JSON file, returning the options.
// Returns an empty list (not an error) if the file doesn't exist.
func (f *FileSource) GetOptions(ctx context.Context) ([]bot.PlayOption, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var options []bot.PlayOption
	if err := json.Unmarshal(data, &options); err != nil {
		// Try as array of strings
		var names []string
		if err := json.Unmarshal(data, &names); err != nil {
			return nil, err
		}
		for _, name := range names {
			options = append(options, bot.PlayOption{Name: name})
		}
	}
	return options, nil
}
