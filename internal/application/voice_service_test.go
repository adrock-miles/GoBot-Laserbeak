package application

import (
	"context"
	"testing"

	"github.com/adrock-miles/go-laserbeak/internal/domain/bot"
)

// --- Mocks ---

type mockSTT struct {
	text string
	err  error
}

func (m *mockSTT) Transcribe(_ context.Context, _ []byte) (string, error) {
	return m.text, m.err
}

type mockLLM struct {
	reply string
	err   error
}

func (m *mockLLM) ChatCompletion(_ context.Context, _ []bot.LLMMessage) (string, error) {
	return m.reply, m.err
}

type mockPlayOptions struct {
	options []bot.PlayOption
	err     error
}

func (m *mockPlayOptions) GetOptions(_ context.Context) ([]bot.PlayOption, error) {
	return m.options, m.err
}

// helper to create a VoiceService with "laser" wake phrase and no LLM/play options.
func newTestService() *VoiceService {
	return NewVoiceService(&mockSTT{}, "laser", nil, nil)
}

// helper that parses a transcription and returns the command text.
// Returns empty string if no command was matched.
func parse(t *testing.T, svc *VoiceService, transcription string) string {
	t.Helper()
	cmd, ok := svc.parseCommand(context.Background(), transcription)
	if !ok {
		return ""
	}
	return cmd.Text
}

// --- Wake phrase detection ---

func TestWakePhrase_ExactMatch(t *testing.T) {
	svc := newTestService()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercase", "laser stop", "!stop"},
		{"capitalized", "Laser stop", "!stop"},
		{"all caps", "LASER stop", "!stop"},
		{"mixed case", "lAsEr stop", "!stop"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(t, svc, tt.input)
			if got != tt.want {
				t.Errorf("parse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestWakePhrase_AlternateSpellings(t *testing.T) {
	svc := newTestService()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lazer", "lazer stop", "!stop"},
		{"Lazer", "Lazer stop", "!stop"},
		{"LAZER", "LAZER stop", "!stop"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(t, svc, tt.input)
			if got != tt.want {
				t.Errorf("parse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestWakePhrase_FillerWordsBefore(t *testing.T) {
	svc := newTestService()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"hey laser", "hey laser stop", "!stop"},
		{"yo laser", "yo laser stop", "!stop"},
		{"ok laser", "ok laser stop", "!stop"},
		{"hey lazer", "hey lazer stop", "!stop"},
		{"Hey Laser", "Hey Laser stop", "!stop"},
		{"OK LASER", "OK LASER stop", "!stop"},
		{"oh hey laser", "oh hey laser stop", "!stop"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(t, svc, tt.input)
			if got != tt.want {
				t.Errorf("parse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestWakePhrase_NoMatch(t *testing.T) {
	svc := newTestService()

	tests := []struct {
		name  string
		input string
	}{
		{"no wake phrase", "stop the music"},
		{"play without wake", "play something"},
		{"blazer not laser", "blazer stop"},
		{"empty", ""},
		{"just filler", "hey yo"},
		{"wake phrase buried too deep", "i was just saying hey laser stop"},
		{"partial word", "lasers stop"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(t, svc, tt.input)
			if got != "" {
				t.Errorf("parse(%q) = %q, want no match", tt.input, got)
			}
		})
	}
}

func TestWakePhrase_OnlyWakePhrase(t *testing.T) {
	svc := newTestService()

	tests := []struct {
		name  string
		input string
	}{
		{"just laser", "laser"},
		{"hey laser", "hey laser"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(t, svc, tt.input)
			if got != "" {
				t.Errorf("parse(%q) = %q, want no match (no command after wake phrase)", tt.input, got)
			}
		})
	}
}

// --- Stop command ---

func TestStopCommand(t *testing.T) {
	svc := newTestService()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic", "laser stop", "!stop"},
		{"with punctuation", "laser stop!", "!stop"},
		{"with period", "laser stop.", "!stop"},
		{"trailing words", "laser stop please", "!stop"},
		{"trailing words with punctuation", "laser stop it now!", "!stop"},
		{"caps", "laser Stop", "!stop"},
		{"all caps", "laser STOP", "!stop"},
		{"with filler prefix", "hey laser stop", "!stop"},
		{"alternate spelling stop", "lazer stop", "!stop"},
		{"filler plus alternate", "yo lazer stop!", "!stop"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(t, svc, tt.input)
			if got != tt.want {
				t.Errorf("parse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// --- Play random ---

func TestPlayRandom(t *testing.T) {
	svc := newTestService()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic", "laser play random", "!pr"},
		{"capitalized", "laser play Random", "!pr"},
		{"all caps", "laser play RANDOM", "!pr"},
		{"something random", "laser play something random", "!pr"},
		{"random song", "laser play a random song", "!pr"},
		{"with filler", "hey laser play random", "!pr"},
		{"alternate spelling", "lazer play random", "!pr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(t, svc, tt.input)
			if got != tt.want {
				t.Errorf("parse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// --- Play command (no LLM, passthrough) ---

func TestPlayCommand_Passthrough(t *testing.T) {
	svc := newTestService()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic query", "laser play never gonna give you up", "!play never gonna give you up"},
		{"with filler", "hey laser play some song", "!play some song"},
		{"alternate spelling", "lazer play test", "!play test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(t, svc, tt.input)
			if got != tt.want {
				t.Errorf("parse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPlayCommand_EmptyQuery(t *testing.T) {
	svc := newTestService()

	got := parse(t, svc, "laser play")
	if got != "" {
		t.Errorf("parse(%q) = %q, want no match (empty play query)", "laser play", got)
	}
}

// --- Play command with LLM matching ---

func TestPlayCommand_LLMMatching(t *testing.T) {
	llm := &mockLLM{reply: "itsworking"}
	opts := &mockPlayOptions{options: []bot.PlayOption{
		{Name: "itsworking"},
		{Name: "miragewish"},
	}}
	svc := NewVoiceService(&mockSTT{}, "laser", llm, opts)

	got := parse(t, svc, "laser play its working")
	if got != "!play itsworking" {
		t.Errorf("parse with LLM = %q, want %q", got, "!play itsworking")
	}
}

func TestPlayCommand_LLMFallback_NoOptions(t *testing.T) {
	llm := &mockLLM{}
	opts := &mockPlayOptions{options: []bot.PlayOption{}}
	svc := NewVoiceService(&mockSTT{}, "laser", llm, opts)

	got := parse(t, svc, "laser play some song")
	if got != "!play some song" {
		t.Errorf("parse with empty options = %q, want %q", got, "!play some song")
	}
}

// --- HandleVoice integration ---

func TestHandleVoice_TranscribesAndParses(t *testing.T) {
	stt := &mockSTT{text: "hey laser stop"}
	svc := NewVoiceService(stt, "laser", nil, nil)

	got, err := svc.HandleVoice(context.Background(), "ch1", "u1", []byte("fake-audio"))
	if err != nil {
		t.Fatalf("HandleVoice error: %v", err)
	}
	if got != "!stop" {
		t.Errorf("HandleVoice = %q, want %q", got, "!stop")
	}
}

func TestHandleVoice_EmptyTranscription(t *testing.T) {
	stt := &mockSTT{text: ""}
	svc := NewVoiceService(stt, "laser", nil, nil)

	got, err := svc.HandleVoice(context.Background(), "ch1", "u1", []byte("fake-audio"))
	if err != nil {
		t.Fatalf("HandleVoice error: %v", err)
	}
	if got != "" {
		t.Errorf("HandleVoice = %q, want empty", got)
	}
}

func TestHandleVoice_NoWakePhrase(t *testing.T) {
	stt := &mockSTT{text: "hello there"}
	svc := NewVoiceService(stt, "laser", nil, nil)

	got, err := svc.HandleVoice(context.Background(), "ch1", "u1", []byte("fake-audio"))
	if err != nil {
		t.Fatalf("HandleVoice error: %v", err)
	}
	if got != "" {
		t.Errorf("HandleVoice = %q, want empty", got)
	}
}

// --- Custom wake phrase ---

func TestCustomWakePhrase(t *testing.T) {
	svc := NewVoiceService(&mockSTT{}, "jarvis", nil, nil)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"exact", "jarvis stop", "!stop"},
		{"with filler", "hey jarvis stop", "!stop"},
		{"caps", "JARVIS stop", "!stop"},
		{"wrong phrase", "laser stop", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(t, svc, tt.input)
			if got != tt.want {
				t.Errorf("parse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
