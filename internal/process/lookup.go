package process

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Info holds basic information about a process owning a port.
type Info struct {
	PID  int
	Name string
}

// String returns a human-readable representation of the process info.
func (p Info) String() string {
	if p.PID == 0 {
		return "unknown"
	}
	return fmt.Sprintf("%s(%d)", p.Name, p.PID)
}

// Lookup attempts to find the process that owns the given port using lsof.
// Returns an empty Info and no error when no process is found.
func Lookup(port int) (Info, error) {
	cmd := exec.Command("lsof", "-iTCP:"+strconv.Itoa(port), "-sTCP:LISTEN", "-n", "-P", "-F", "pc")
	out, err := cmd.Output()
	if err != nil {
		// lsof exits non-zero when no process is found; treat as empty result.
		return Info{}, nil
	}

	return parseOutput(string(out))
}

// parseOutput parses lsof -F output to extract PID and process name.
func parseOutput(output string) (Info, error) {
	var info Info
	for _, line := range strings.Split(output, "\n") {
		if len(line) < 2 {
			continue
		}
		switch line[0] {
		case 'p':
			pid, err := strconv.Atoi(strings.TrimSpace(line[1:]))
			if err != nil {
				return Info{}, fmt.Errorf("process: parse pid %q: %w", line[1:], err)
			}
			info.PID = pid
		case 'c':
			info.Name = strings.TrimSpace(line[1:])
		}
	}
	return info, nil
}
