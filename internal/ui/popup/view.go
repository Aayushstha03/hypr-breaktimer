package popup

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	accent := lipgloss.Color("6")
	accentSuccess := lipgloss.Color("10")
	mutedC := lipgloss.Color("8")
	muted := lipgloss.NewStyle().Foreground(mutedC).Render

	motivationStyle := lipgloss.NewStyle().Bold(true).Foreground(accentSuccess)
	nudgeStyle := lipgloss.NewStyle().Bold(true).Foreground(accentSuccess)
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(accent)
	dividerStyle := lipgloss.NewStyle().Foreground(mutedC)
	progressStyle := lipgloss.NewStyle().Foreground(accent)
	boxStyle := lipgloss.NewStyle().
		Padding(1, 3).
		Border(lipgloss.ThickBorder()).
		BorderForeground(accent)

	divider := dividerStyle.Render("\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500")

	var body string
	switch m.state {
	case statePrompt:
		body = motivationStyle.Render(nonEmpty(m.message, "Stand up, look away from the screen, and stretch.")) + "\n\n" +
			titleStyle.Render(nonEmpty(m.title, "Take a break")) + "\n" +
			muted(fmt.Sprintf("for %s", m.breakDuration)) + "\n\n" +
			divider + "\n" +
			muted("b / enter  start break") + "\n" +
			muted(fmt.Sprintf("s          snooze (%s)", m.snoozeDuration)) + "\n" +
			muted("q / esc    quit")
	case stateBreaking:
		remaining := max(time.Until(m.breakEndsAt), 0)
		bar := progressStyle.Render(m.progress.ViewAs(m.breakProgress))
		nudge := ""
		if m.nudge != "" {
			nudge = nudgeStyle.Render(m.nudge) + "\n\n"
		}
		body = titleStyle.Render("Break in progress") + "\n\n" +
			nudge +
			fmt.Sprintf("%s\n", roundSeconds(remaining)) +
			bar + "\n\n" +
			muted("e          end early (counts as a break taken)")
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

func roundSeconds(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	// Rounded down for less jitter.
	secs := int(d.Seconds())
	return (time.Duration(secs) * time.Second).String()
}
