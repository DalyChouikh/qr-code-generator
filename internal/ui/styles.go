// Package ui provides the terminal user interface components.
package ui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	primaryColor   = lipgloss.AdaptiveColor{Light: "#7C3AED", Dark: "#A78BFA"}
	secondaryColor = lipgloss.AdaptiveColor{Light: "#059669", Dark: "#34D399"}
	accentColor    = lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#F87171"}
	subtleColor    = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}
	textColor      = lipgloss.AdaptiveColor{Light: "#1F2937", Dark: "#F9FAFB"}
	bgColor        = lipgloss.AdaptiveColor{Light: "#F3F4F6", Dark: "#1F2937"}
)

// Styles defines all the styling for the TUI.
type Styles struct {
	App           lipgloss.Style
	Title         lipgloss.Style
	Subtitle      lipgloss.Style
	Header        lipgloss.Style
	FocusedInput  lipgloss.Style
	BlurredInput  lipgloss.Style
	Cursor        lipgloss.Style
	Label         lipgloss.Style
	LabelFocused  lipgloss.Style
	Option        lipgloss.Style
	OptionActive  lipgloss.Style
	Help          lipgloss.Style
	Error         lipgloss.Style
	Success       lipgloss.Style
	Button        lipgloss.Style
	ButtonActive  lipgloss.Style
	Preview       lipgloss.Style
	PreviewBorder lipgloss.Style
	Divider       lipgloss.Style
	Step          lipgloss.Style
	StepActive    lipgloss.Style
	StepCompleted lipgloss.Style
}

// NewStyles creates and returns the default styles.
func NewStyles() *Styles {
	return &Styles{
		App: lipgloss.NewStyle().
			Padding(1, 2).
			MarginLeft(1),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1).
			Padding(0, 1).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(primaryColor).
			BorderBottom(true),

		Subtitle: lipgloss.NewStyle().
			Foreground(subtleColor).
			Italic(true).
			MarginBottom(2),

		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			MarginBottom(1).
			MarginTop(1),

		FocusedInput: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1).
			Width(50),

		BlurredInput: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(subtleColor).
			Padding(0, 1).
			Width(50),

		Cursor: lipgloss.NewStyle().
			Foreground(primaryColor),

		Label: lipgloss.NewStyle().
			Foreground(subtleColor).
			Bold(false).
			MarginRight(1),

		LabelFocused: lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginRight(1),

		Option: lipgloss.NewStyle().
			Foreground(subtleColor).
			Padding(0, 2),

		OptionActive: lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 2),

		Help: lipgloss.NewStyle().
			Foreground(subtleColor).
			Italic(true).
			MarginTop(2),

		Error: lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Padding(0, 1).
			MarginTop(1),

		Success: lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			Padding(1, 2).
			MarginTop(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor),

		Button: lipgloss.NewStyle().
			Foreground(subtleColor).
			Padding(0, 3).
			Width(13).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(subtleColor),

		ButtonActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Padding(0, 3).
			Width(13).
			Align(lipgloss.Center).
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor),

		Preview: lipgloss.NewStyle().
			Foreground(textColor).
			Padding(1, 2).
			MarginTop(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(subtleColor).
			Width(54),

		PreviewBorder: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1).
			Width(56),

		Divider: lipgloss.NewStyle().
			Foreground(subtleColor).
			SetString("─────────────────────────────────────────────────"),

		Step: lipgloss.NewStyle().
			Foreground(subtleColor).
			Width(3).
			Align(lipgloss.Center),

		StepActive: lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Width(3).
			Align(lipgloss.Center),

		StepCompleted: lipgloss.NewStyle().
			Foreground(secondaryColor).
			Width(3).
			Align(lipgloss.Center),
	}
}

// RenderProgressBar renders a step progress indicator.
func (s *Styles) RenderProgressBar(current, total int) string {
	var steps string
	for i := 1; i <= total; i++ {
		var step string
		switch {
		case i < current:
			step = s.StepCompleted.Render("✓")
		case i == current:
			step = s.StepActive.Render("●")
		default:
			step = s.Step.Render("○")
		}
		steps += step
		if i < total {
			if i < current {
				steps += lipgloss.NewStyle().Foreground(secondaryColor).Render("───")
			} else if i == current {
				steps += lipgloss.NewStyle().Foreground(primaryColor).Render("───")
			} else {
				steps += lipgloss.NewStyle().Foreground(subtleColor).Render("───")
			}
		}
	}
	return lipgloss.NewStyle().MarginBottom(1).Render(steps)
}
