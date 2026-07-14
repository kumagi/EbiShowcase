package uilab

import "strings"

type Align int

const (
	Left Align = iota
	Center
	Right
)

func X(anchor, width float64, align Align) float64 {
	if align == Center {
		return anchor - width/2
	}
	if align == Right {
		return anchor - width
	}
	return anchor
}

// Wrap is deliberately text-only: call text/v2 Measure for pixel placement
// after it chooses safe break opportunities. Japanese can break per rune;
// English retains words where possible.
func Wrap(s string, limit int, japanese bool) []string {
	if limit < 1 {
		return []string{s}
	}
	var out []string
	line := ""
	for _, part := range strings.Fields(s) {
		if japanese {
			for _, r := range part {
				if len([]rune(line)) >= limit {
					out = append(out, line)
					line = ""
				}
				line += string(r)
			}
			continue
		}
		next := part
		if line != "" {
			next = " " + part
		}
		if len([]rune(line))+len([]rune(next)) > limit && line != "" {
			out = append(out, line)
			line = part
		} else {
			line += next
		}
	}
	if line != "" {
		out = append(out, line)
	}
	return out
}
