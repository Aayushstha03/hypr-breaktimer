package popup

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	accent := lipgloss.Color("2")
	accentSuccess := lipgloss.Color("3")
	mutedC := lipgloss.Color("7")
	muted := lipgloss.NewStyle().Foreground(mutedC).Render

	compact := m.width > 0 && m.height > 0 && (m.width < 60 || m.height < 18)

	motivationStyle := lipgloss.NewStyle().Bold(true).Foreground(accentSuccess)
	nudgeStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1"))
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(accent)
	dividerStyle := lipgloss.NewStyle().Foreground(mutedC)
	progressStyle := lipgloss.NewStyle().Foreground(accent)
	boxStyle := lipgloss.NewStyle().Padding(1, 3)
	if compact {
		boxStyle = lipgloss.NewStyle().Padding(0, 1)
	} else {
		boxStyle = boxStyle.Border(lipgloss.ThickBorder()).BorderForeground(accent)
	}

	// Fix the popup width based on terminal size so the box doesn't "breathe"
	// when content length changes. Render content to an inner width and let the
	// box style (padding/border) add its frame around it.
	frameW, _ := boxStyle.GetFrameSize()
	outerW := 86
	if m.width > 0 {
		outerW = min(outerW, m.width)
	}
	outerW = max(24, outerW)
	innerW := max(16, outerW-frameW)

	wrap := lipgloss.NewStyle().Width(innerW)
	divider := dividerStyle.Render(strings.Repeat("\u2500", innerW))

	var body string
	switch m.state {
	case statePrompt:
		dividerLine := ""
		if !compact {
			dividerLine = wrap.Render(divider) + "\n"
		}
		body = motivationStyle.Width(innerW).Render(nonEmpty(m.message, "Stand up, look away from the screen, and stretch.")) + "\n\n" +
			titleStyle.Width(innerW).Render(nonEmpty(m.title, "Take a break")) + "\n" +
			muted(fmt.Sprintf("for %s", m.breakDuration)) + "\n\n" +
			dividerLine +
			muted("b / enter  start break") + "\n" +
			muted(fmt.Sprintf("s          snooze (%s)", m.snoozeDuration)) + "\n"
	case stateBreaking:
		remaining := max(time.Until(m.breakEndsAt), 0)
		pm := m.progress
		pm.Width = innerW
		bar := progressStyle.Width(innerW).Render(pm.ViewAs(m.breakProgress))

		reservedNudgeH := 1
		for _, msg := range exitAttemptMessages {
			h := lipgloss.Height(nudgeStyle.Width(innerW).Render(msg))
			if h > reservedNudgeH {
				reservedNudgeH = h
			}
		}
		nudgeBlock := padToHeight(nudgeStyle.Width(innerW).Render(m.nudge), reservedNudgeH)
		if m.nudge == "" {
			nudgeBlock = padToHeight("", reservedNudgeH)
		}

		body = titleStyle.Width(innerW).Render("Break in progress") + "\n\n" +
			nudgeBlock + "\n\n" +
			wrap.Render(roundSeconds(remaining)) + "\n" +
			bar + "\n\n" +
			wrap.Render(muted("e          end early (counts as a break taken)"))
	case stateDone:
		body = titleStyle.Width(innerW).Render("Nice.") + "\n\n" + wrap.Render("Break complete.")
	default:
		body = ""
	}

	content := boxStyle.Render(body)
	if m.width <= 0 || m.height <= 0 {
		return content
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func padToHeight(s string, height int) string {
	if height <= 0 {
		return ""
	}
	if s == "" {
		lines := make([]string, height)
		return strings.Join(lines, "\n")
	}
	lines := strings.Split(s, "\n")
	if len(lines) >= height {
		return s
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
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
