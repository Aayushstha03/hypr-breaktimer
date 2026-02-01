package popup

import (
	"context"
	"math"
	"time"

	"github.com/Aayushstha03/hypr-breaktimer/internal/config"
	"github.com/Aayushstha03/hypr-breaktimer/internal/state"
	"github.com/Aayushstha03/hypr-breaktimer/internal/xdg"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type Options struct{}

func Run(ctx context.Context, _ Options) (int, error) {
	configPath, _ := xdg.ConfigFile()
	statePath, _ := xdg.StateFile()

	cfg, _ := config.Load(configPath)
	st, _ := state.Load(statePath)

	now := time.Now()
	st.LastPopupShownAt = &now
	_ = state.SaveAtomic(statePath, st)

	m := newModel(cfg)
	// newModel already loaded state; keep this in sync anyway.
	m.lastPopupShownAt = st.LastPopupShownAt
	m.lastBreakStartedAt = st.LastBreakStartedAt
	m.lastBreakCompletedAt = st.LastBreakCompletedAt
	m.snoozedUntil = st.SnoozedUntil
	m.lastAction = st.LastAction
	m.lastActionAt = st.LastActionAt
	// respect config auto-start
	if cfg.Popup.AutoStartBreak {
		m.state = stateBreaking
		m.breakEndsAt = time.Now().Add(m.breakDuration)
		m.writeAction(state.ActionStarted)
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return 1, err
	}
	mm, _ := final.(model)
	if mm.exitCode == 0 {
		return 0, nil
	}
	return mm.exitCode, nil
}

type stateID int

const (
	statePrompt stateID = iota
	stateBreaking
	stateDone
)

type tickMsg time.Time

type model struct {
	state stateID

	width  int
	height int

	breakDuration  time.Duration
	breakEndsAt    time.Time
	snoozeDuration time.Duration
	breakProgress  float64
	progress       progress.Model

	title     string
	message   string
	statePath string

	lastPopupShownAt     *time.Time
	lastBreakStartedAt   *time.Time
	lastBreakCompletedAt *time.Time
	snoozedUntil         *time.Time
	lastAction           state.Action
	lastActionAt         *time.Time

	exitCode int
}

func newModel(cfg config.Config) model {
	statePath := mustStatePath()
	st, _ := state.Load(statePath)

	pm := progress.New(progress.WithoutPercentage())
	pm.FullColor = "6"  // ANSI cyan (theme-defined)
	pm.EmptyColor = "8" // ANSI bright black (theme-defined)

	return model{
		state:                statePrompt,
		breakDuration:        cfg.Schedule.BreakDuration.Duration(),
		snoozeDuration:       cfg.Schedule.SnoozeDuration.Duration(),
		breakProgress:        0,
		progress:             pm,
		title:                cfg.Popup.Title,
		message:              cfg.Popup.Message,
		statePath:            statePath,
		lastPopupShownAt:     st.LastPopupShownAt,
		lastBreakStartedAt:   st.LastBreakStartedAt,
		lastBreakCompletedAt: st.LastBreakCompletedAt,
		snoozedUntil:         st.SnoozedUntil,
		lastAction:           st.LastAction,
		lastActionAt:         st.LastActionAt,
		exitCode:             0,
	}
}

func mustStatePath() string {
	p, err := xdg.StateFile()
	if err != nil {
		return ""
	}
	return p
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		s := msg.String()
		switch s {
		case "ctrl+c", "esc", "q":
			// Treat quitting as break taken (except snooze).
			m.writeAction(state.ActionCompleted)
			m.exitCode = 0
			return m, tea.Quit
		}
		switch m.state {
		case statePrompt:
			switch s {
			case "enter", "b":
				m.state = stateBreaking
				m.breakEndsAt = time.Now().Add(m.breakDuration)
				m.breakProgress = 0
				m.writeAction(state.ActionStarted)
				return m, tickCmd()
			case "s":
				m.writeSnooze(m.snoozeDuration)
				m.exitCode = 0
				return m, tea.Quit
			}
		case stateBreaking:
			switch s {
			case "e":
				m.state = stateDone
				m.writeAction(state.ActionCompleted)
				return m, doneCmd()
			}
		case stateDone:
			return m, tea.Quit
		}
	case progress.FrameMsg:
		if m.state != stateBreaking {
			return m, nil
		}
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case tickMsg:
		if m.state != stateBreaking {
			return m, nil
		}

		m.breakProgress = clamp01(breakProgress(m.breakDuration, m.breakEndsAt, time.Now()))
		progressCmd := m.progress.SetPercent(m.breakProgress)

		if time.Now().After(m.breakEndsAt) {
			m.state = stateDone
			m.writeAction(state.ActionCompleted)
			return m, doneCmd()
		}
		return m, tea.Batch(tickCmd(), progressCmd)
	}
	return m, nil
}

func breakProgress(d time.Duration, endsAt time.Time, now time.Time) float64 {
	if d <= 0 {
		return 1
	}
	startedAt := endsAt.Add(-d)
	elapsed := now.Sub(startedAt)
	return float64(elapsed) / float64(d)
}

func clamp01(v float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func (m *model) writeAction(a state.Action) {
	if m.statePath == "" {
		return
	}
	st, err := state.Load(m.statePath)
	if err != nil {
		return
	}
	now := time.Now()
	st.LastAction = a
	st.LastActionAt = &now
	if a == state.ActionCompleted {
		st.LastBreakCompletedAt = &now
	}
	if a == state.ActionStarted {
		st.LastBreakStartedAt = &now
	}
	if st.LastPopupShownAt == nil {
		st.LastPopupShownAt = &now
	}

	m.lastPopupShownAt = st.LastPopupShownAt
	m.lastBreakStartedAt = st.LastBreakStartedAt
	m.lastBreakCompletedAt = st.LastBreakCompletedAt
	m.snoozedUntil = st.SnoozedUntil
	m.lastAction = st.LastAction
	m.lastActionAt = st.LastActionAt
	_ = state.SaveAtomic(m.statePath, st)
}

func (m *model) writeSnooze(d time.Duration) {
	if m.statePath == "" {
		return
	}
	st, err := state.Load(m.statePath)
	if err != nil {
		return
	}
	now := time.Now()
	until := now.Add(d)
	st.SnoozedUntil = &until
	st.LastAction = state.ActionSnoozed
	st.LastActionAt = &now
	if st.LastPopupShownAt == nil {
		st.LastPopupShownAt = &now
	}

	m.lastPopupShownAt = st.LastPopupShownAt
	m.lastBreakStartedAt = st.LastBreakStartedAt
	m.lastBreakCompletedAt = st.LastBreakCompletedAt
	m.snoozedUntil = st.SnoozedUntil
	m.lastAction = st.LastAction
	m.lastActionAt = st.LastActionAt
	_ = state.SaveAtomic(m.statePath, st)
}

func tickCmd() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func doneCmd() tea.Cmd {
	return tea.Tick(750*time.Millisecond, func(time.Time) tea.Msg { return tea.Quit() })
}
