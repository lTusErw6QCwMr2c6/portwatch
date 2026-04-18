package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ruleJSON struct {
	Name     string `json:"name"`
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol"`
	Action   string `json:"action"`
}

// Load reads a JSON file containing an array of rules and returns a Policy.
func Load(path string) (*Policy, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("policy: open %s: %w", path, err)
	}
	defer f.Close()

	var raw []ruleJSON
	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return nil, fmt.Errorf("policy: decode: %w", err)
	}

	rules := make([]Rule, 0, len(raw))
	for _, r := range raw {
		a, err := parseAction(r.Action)
		if err != nil {
			return nil, fmt.Errorf("policy: rule %q: %w", r.Name, err)
		}
		rules = append(rules, Rule{
			Name:     r.Name,
			Port:     r.Port,
			Protocol: strings.ToLower(r.Protocol),
			Action:   a,
		})
	}
	return New(rules), nil
}

func parseAction(s string) (Action, error) {
	switch Action(strings.ToLower(s)) {
	case ActionAllow:
		return ActionAllow, nil
	case ActionDeny:
		return ActionDeny, nil
	case ActionWarn:
		return ActionWarn, nil
	default:
		return "", fmt.Errorf("unknown action %q", s)
	}
}
