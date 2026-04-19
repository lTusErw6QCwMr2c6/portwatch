package masking

import "strings"

// Level controls how aggressively port/process info is masked.
type Level int

const (
	LevelNone  Level = iota // no masking
	LevelPartial            // mask middle characters
	LevelFull               // replace with placeholder
)

// Config holds masking configuration.
type Config struct {
	Level       Level
	Placeholder string
	Fields      []string // field names to mask
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Level:       LevelPartial,
		Placeholder: "[MASKED]",
		Fields:      []string{"process", "user", "cmdline"},
	}
}

// Masker applies masking rules to string values.
type Masker struct {
	cfg    Config
	fields map[string]bool
}

// New creates a Masker from the given Config.
func New(cfg Config) *Masker {
	f := make(map[string]bool, len(cfg.Fields))
	for _, field := range cfg.Fields {
		f[strings.ToLower(field)] = true
	}
	return &Masker{cfg: cfg, fields: f}
}

// ShouldMask reports whether the named field should be masked.
func (m *Masker) ShouldMask(field string) bool {
	return m.fields[strings.ToLower(field)]
}

// Apply masks the given value according to the configured level.
func (m *Masker) Apply(value string) string {
	if value == "" {
		return value
	}
	switch m.cfg.Level {
	case LevelFull:
		return m.cfg.Placeholder
	case LevelPartial:
		return partialMask(value)
	default:
		return value
	}
}

// ApplyField masks value only if the field is in the mask list.
func (m *Masker) ApplyField(field, value string) string {
	if !m.ShouldMask(field) {
		return value
	}
	return m.Apply(value)
}

func partialMask(s string) string {
	if len(s) <= 2 {
		return strings.Repeat("*", len(s))
	}
	return string(s[0]) + strings.Repeat("*", len(s)-2) + string(s[len(s)-1])
}
