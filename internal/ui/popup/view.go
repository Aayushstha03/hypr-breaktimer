package popup

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230"))
	boxStyle := lipgloss.NewStyle().
		Padding(1, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("241"))

	var body string
	switch m.state {
	case statePrompt:
		body = titleStyle.Render(nonEmpty(m.title, "Take a break")) + "\n\n" +
			nonEmpty(m.message, "Stand up, look away from the screen, and stretch.") + "\n\n" +
			metaBlock(m) + "\n\n" +
			"b / enter  start break\n" +
			"s          snooze\n" +
			"d          dismiss\n" +
			"q / esc    quit"
	case stateBreaking:
		remaining := max(time.Until(m.breakEndsAt), 0)
		body = titleStyle.Render("Break in progress") + "\n\n" +
			fmt.Sprintf("Time left: %s\n\n", roundSeconds(remaining)) +
			metaBlock(m) + "\n\n" +
			"e          end early\n" +
			"q / esc    quit"
	case stateDone:
		body = titleStyle.Render("Nice.") + "\n\n" + "Break complete."
	default:
		body = ""
	}

	content := boxStyle.Render(body)
	if m.width <= 0 || m.height <= 0 {
		return content
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func nonEmpty(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func metaBlock(m model) string {
	meta := lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render

	line1 := fmt.Sprintf("break %s | snooze %s", m.breakDuration, m.snoozeDuration)

	cfg := m.configPath
	if cfg == "" {
		cfg = "(unknown config path)"
	}
	st := m.statePath
	if st == "" {
		st = "(unknown state path)"
	}

	line2 := fmt.Sprintf("cfg: %s", cfg)
	line3 := fmt.Sprintf("state: %s", st)

	line4 := ""
	if m.lastPopupShownAt != nil {
		line4 = fmt.Sprintf("last popup: %s", m.lastPopupShownAt.Format(time.RFC3339))
	}
	line5 := ""
	if m.snoozedUntil != nil {
		line5 = fmt.Sprintf("snoozed until: %s", m.snoozedUntil.Format(time.RFC3339))
	}
	line6 := ""
	if m.lastAction != "" {
		at := ""
		if m.lastActionAt != nil {
			at = m.lastActionAt.Format(time.RFC3339)
		}
		line6 = fmt.Sprintf("last action: %s %s", m.lastAction, at)
	}

	out := meta(line1) + "\n" + muted(line2) + "\n" + muted(line3)
	if line4 != "" {
		out += "\n" + meta(line4)
	}
	if line5 != "" {
		out += "\n" + meta(line5)
	}
	if line6 != "" {
		out += "\n" + meta(line6)
	}
	return out
}

func roundSeconds(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	// Rounded down for less jitter.
	secs := int(d.Seconds())
	return (time.Duration(secs) * time.Second).String()
}
