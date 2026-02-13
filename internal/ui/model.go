package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/DalyChouikh/internal/config"
	"github.com/DalyChouikh/internal/generator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Step represents the current step in the wizard.
type Step int

const (
	StepURL Step = iota
	StepFormat
	StepColor
	StepSize
	StepOutput
	StepConfirm
	StepComplete
)

const totalSteps = 6

// Model represents the application state.
type Model struct {
	styles *Styles
	config *config.QRConfig

	// Current step
	step Step

	// Input fields
	urlInput      textinput.Model
	sizeInput     textinput.Model
	outputInput   textinput.Model
	colorInput    textinput.Model

	// Selection states
	formatIndex int // 0 = PNG, 1 = SVG
	colorIndex  int // Index in predefined colors, -1 for custom
	colorNames  []string

	// UI state
	err         error
	successPath string
	quitting    bool

	// Window size
	width  int
	height int
}

// New creates a new Model with default values.
func New() Model {
	styles := NewStyles()
	cfg := config.DefaultConfig()

	// URL input
	urlInput := textinput.New()
	urlInput.Placeholder = "https://example.com"
	urlInput.CharLimit = 2048
	urlInput.Width = 46
	urlInput.Focus()

	// Size input
	sizeInput := textinput.New()
	sizeInput.Placeholder = "256"
	sizeInput.CharLimit = 4
	sizeInput.Width = 46

	// Output input
	outputInput := textinput.New()
	outputInput.Placeholder = "qrcode"
	outputInput.CharLimit = 256
	outputInput.Width = 46

	// Color input (for custom hex)
	colorInput := textinput.New()
	colorInput.Placeholder = "#000000"
	colorInput.CharLimit = 7
	colorInput.Width = 46

	// Get home directory for default output path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	cfg.OutputPath = filepath.Join(homeDir, "qrcode.png")

	return Model{
		styles:      styles,
		config:      cfg,
		step:        StepURL,
		urlInput:    urlInput,
		sizeInput:   sizeInput,
		outputInput: outputInput,
		colorInput:  colorInput,
		formatIndex: 0,
		colorIndex:  0, // Default to black
		colorNames:  config.GetPredefinedColorNames(),
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Global key handlers
		switch msg.String() {
		case "ctrl+c", "q":
			if m.step == StepComplete {
				m.quitting = true
				return m, tea.Quit
			}
			if msg.String() == "ctrl+c" {
				m.quitting = true
				return m, tea.Quit
			}
		case "esc":
			if m.step > StepURL && m.step < StepComplete {
				m.step--
				m.err = nil
				return m, m.focusStep()
			}
		}

		// Step-specific handlers
		switch m.step {
		case StepURL:
			return m.handleURLStep(msg)
		case StepFormat:
			return m.handleFormatStep(msg)
		case StepColor:
			return m.handleColorStep(msg)
		case StepSize:
			return m.handleSizeStep(msg)
		case StepOutput:
			return m.handleOutputStep(msg)
		case StepConfirm:
			return m.handleConfirmStep(msg)
		case StepComplete:
			return m.handleCompleteStep(msg)
		}
	}

	// Update current text input (for non-key messages like cursor blink)
	var cmd tea.Cmd
	switch m.step {
	case StepURL:
		m.urlInput, cmd = m.urlInput.Update(msg)
	case StepSize:
		m.sizeInput, cmd = m.sizeInput.Update(msg)
	case StepOutput:
		m.outputInput, cmd = m.outputInput.Update(msg)
	case StepColor:
		if m.colorIndex == -1 {
			m.colorInput, cmd = m.colorInput.Update(msg)
		}
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleURLStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		url := strings.TrimSpace(m.urlInput.Value())
		if url == "" {
			m.err = fmt.Errorf("please enter a URL or text to encode")
			return m, nil
		}
		m.config.Content = url
		m.err = nil
		m.step = StepFormat
		m.urlInput.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.urlInput, cmd = m.urlInput.Update(msg)
	return m, cmd
}

func (m Model) handleFormatStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "h":
		if m.formatIndex > 0 {
			m.formatIndex--
		}
	case "right", "l":
		if m.formatIndex < 1 {
			m.formatIndex++
		}
	case "enter", " ":
		if m.formatIndex == 0 {
			m.config.Format = config.FormatPNG
		} else {
			m.config.Format = config.FormatSVG
		}
		m.err = nil
		m.step = StepColor
	case "1":
		m.formatIndex = 0
		m.config.Format = config.FormatPNG
	case "2":
		m.formatIndex = 1
		m.config.Format = config.FormatSVG
	}
	return m, nil
}

func (m Model) handleColorStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle custom color input if selected
	if m.colorIndex == -1 {
		switch msg.String() {
		case "enter":
			hexColor := strings.TrimSpace(m.colorInput.Value())
			if hexColor == "" {
				m.colorIndex = 0 // Reset to first predefined color
				return m, nil
			}
			color, err := config.ParseHexColor(hexColor)
			if err != nil {
				m.err = fmt.Errorf("invalid hex color: %s", hexColor)
				return m, nil
			}
			m.config.Foreground = color
			m.err = nil
			m.step = StepSize
			m.colorInput.Blur()
			return m, m.sizeInput.Focus()
		case "esc":
			m.colorIndex = 0
			m.colorInput.Blur()
			return m, nil
		}

		var cmd tea.Cmd
		m.colorInput, cmd = m.colorInput.Update(msg)
		return m, cmd
	}

	switch msg.String() {
	case "up", "k":
		if m.colorIndex > 0 {
			m.colorIndex--
		}
	case "down", "j":
		if m.colorIndex < len(m.colorNames)-1 {
			m.colorIndex++
		}
	case "c":
		// Switch to custom color input
		m.colorIndex = -1
		m.colorInput.Focus()
		return m, textinput.Blink
	case "enter", " ":
		colorName := m.colorNames[m.colorIndex]
		m.config.Foreground = config.PredefinedColors[colorName]
		m.err = nil
		m.step = StepSize
		return m, m.sizeInput.Focus()
	}
	return m, nil
}

func (m Model) handleSizeStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		sizeStr := strings.TrimSpace(m.sizeInput.Value())
		if sizeStr == "" {
			m.config.Size = 256 // Default size
		} else {
			size, err := strconv.Atoi(sizeStr)
			if err != nil || size < 64 || size > 4096 {
				m.err = fmt.Errorf("size must be a number between 64 and 4096")
				return m, nil
			}
			m.config.Size = size
		}
		m.err = nil
		m.step = StepOutput
		m.sizeInput.Blur()
		return m, m.outputInput.Focus()
	}

	var cmd tea.Cmd
	m.sizeInput, cmd = m.sizeInput.Update(msg)
	return m, cmd
}

func (m Model) handleOutputStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		output := strings.TrimSpace(m.outputInput.Value())
		if output == "" {
			output = "qrcode"
		}

		// Expand ~ to home directory
		if strings.HasPrefix(output, "~") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				output = filepath.Join(homeDir, output[1:])
			}
		}

		// Make absolute if relative
		if !filepath.IsAbs(output) {
			cwd, err := os.Getwd()
			if err == nil {
				output = filepath.Join(cwd, output)
			}
		}

		m.config.SetOutputPath(output)
		m.err = nil
		m.step = StepConfirm
		m.outputInput.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.outputInput, cmd = m.outputInput.Update(msg)
	return m, cmd
}

func (m Model) handleConfirmStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "y", "Y":
		// Generate QR code
		gen := generator.New(m.config)
		if err := gen.Generate(); err != nil {
			m.err = err
			return m, nil
		}
		m.successPath = m.config.OutputPath
		m.step = StepComplete
		return m, nil
	case "n", "N", "esc":
		// Go back to URL step
		m.step = StepURL
		m.urlInput.Focus()
		return m, textinput.Blink
	}
	return m, nil
}

func (m Model) handleCompleteStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "q", " ":
		m.quitting = true
		return m, tea.Quit
	case "r":
		// Reset and start over
		return New(), textinput.Blink
	}
	return m, nil
}

// focusStep focuses the text input for the current step.
func (m *Model) focusStep() tea.Cmd {
	switch m.step {
	case StepURL:
		return m.urlInput.Focus()
	case StepSize:
		return m.sizeInput.Focus()
	case StepOutput:
		return m.outputInput.Focus()
	case StepColor:
		if m.colorIndex == -1 {
			return m.colorInput.Focus()
		}
	default:
		return nil
	}
	return textinput.Blink
}

// View renders the UI.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	// Title
	s.WriteString(m.styles.Title.Render("  QR Code Generator  "))
	s.WriteString("\n")
	s.WriteString(m.styles.Subtitle.Render("Create beautiful QR codes from your terminal"))
	s.WriteString("\n\n")

	// Progress bar
	if m.step < StepComplete {
		s.WriteString(m.styles.RenderProgressBar(int(m.step)+1, totalSteps))
		s.WriteString("\n\n")
	}

	// Current step content
	switch m.step {
	case StepURL:
		s.WriteString(m.renderURLStep())
	case StepFormat:
		s.WriteString(m.renderFormatStep())
	case StepColor:
		s.WriteString(m.renderColorStep())
	case StepSize:
		s.WriteString(m.renderSizeStep())
	case StepOutput:
		s.WriteString(m.renderOutputStep())
	case StepConfirm:
		s.WriteString(m.renderConfirmStep())
	case StepComplete:
		s.WriteString(m.renderCompleteStep())
	}

	// Error display
	if m.err != nil {
		s.WriteString("\n")
		s.WriteString(m.styles.Error.Render("⚠ " + m.err.Error()))
	}

	// Help text
	s.WriteString("\n\n")
	s.WriteString(m.renderHelp())

	return m.styles.App.Render(s.String())
}

func (m Model) renderURLStep() string {
	var s strings.Builder

	s.WriteString(m.styles.Header.Render("Step 1: Enter URL or Text"))
	s.WriteString("\n\n")

	label := m.styles.LabelFocused.Render("Content:")
	s.WriteString(label + "\n")
	s.WriteString(m.styles.FocusedInput.Render(m.urlInput.View()))

	return s.String()
}

func (m Model) renderFormatStep() string {
	var s strings.Builder

	s.WriteString(m.styles.Header.Render("Step 2: Choose Output Format"))
	s.WriteString("\n\n")

	pngStyle := m.styles.Button
	svgStyle := m.styles.Button

	if m.formatIndex == 0 {
		pngStyle = m.styles.ButtonActive
	} else {
		svgStyle = m.styles.ButtonActive
	}

	s.WriteString("  ")
	s.WriteString(pngStyle.Render("  PNG  "))
	s.WriteString("    ")
	s.WriteString(svgStyle.Render("  SVG  "))

	s.WriteString("\n\n")
	s.WriteString(m.styles.Label.Render("PNG: Raster image, best for most uses"))
	s.WriteString("\n")
	s.WriteString(m.styles.Label.Render("SVG: Vector format, scales infinitely"))

	return s.String()
}

func (m Model) renderColorStep() string {
	var s strings.Builder

	s.WriteString(m.styles.Header.Render("Step 3: Choose QR Code Color"))
	s.WriteString("\n\n")

	if m.colorIndex == -1 {
		// Custom color input mode
		s.WriteString(m.styles.LabelFocused.Render("Enter hex color:"))
		s.WriteString("\n")
		s.WriteString(m.styles.FocusedInput.Render(m.colorInput.View()))
		s.WriteString("\n\n")
		s.WriteString(m.styles.Label.Render("Press ESC to go back to color selection"))
	} else {
		// Color selection mode
		for i, name := range m.colorNames {
			var line string
			color := config.PredefinedColors[name]
			colorBox := lipgloss.NewStyle().
				Background(lipgloss.Color(config.ColorToHex(color))).
				Render("  ")

			if i == m.colorIndex {
				cursor := m.styles.OptionActive.Render("▸")
				colorName := m.styles.OptionActive.Render(name)
				line = fmt.Sprintf("%s %s %s", cursor, colorBox, colorName)
			} else {
				cursor := m.styles.Option.Render(" ")
				colorName := m.styles.Option.Render(name)
				line = fmt.Sprintf("%s %s %s", cursor, colorBox, colorName)
			}
			s.WriteString(line + "\n")
		}

		s.WriteString("\n")
		s.WriteString(m.styles.Label.Render("Press 'c' for custom hex color"))
	}

	return s.String()
}

func (m Model) renderSizeStep() string {
	var s strings.Builder

	s.WriteString(m.styles.Header.Render("Step 4: Set Dimensions"))
	s.WriteString("\n\n")

	label := m.styles.LabelFocused.Render("Size (64-4096 pixels):")
	s.WriteString(label + "\n")
	s.WriteString(m.styles.FocusedInput.Render(m.sizeInput.View()))
	s.WriteString("\n\n")
	s.WriteString(m.styles.Label.Render("Leave empty for default (256px)"))

	return s.String()
}

func (m Model) renderOutputStep() string {
	var s strings.Builder

	s.WriteString(m.styles.Header.Render("Step 5: Output Location"))
	s.WriteString("\n\n")

	label := m.styles.LabelFocused.Render("Filename or path:")
	s.WriteString(label + "\n")
	s.WriteString(m.styles.FocusedInput.Render(m.outputInput.View()))
	s.WriteString("\n\n")
	s.WriteString(m.styles.Label.Render("Leave empty for 'qrcode' in current directory"))
	s.WriteString("\n")
	s.WriteString(m.styles.Label.Render("Use ~ for home directory, e.g., ~/Downloads/myqr"))

	return s.String()
}

func (m Model) renderConfirmStep() string {
	var s strings.Builder

	s.WriteString(m.styles.Header.Render("Step 6: Review & Generate"))
	s.WriteString("\n\n")

	// Preview box
	preview := m.styles.Preview.Render(m.renderPreview())
	s.WriteString(preview)
	s.WriteString("\n\n")

	confirmText := lipgloss.NewStyle().Bold(true).Foreground(primaryColor).Render("Generate QR code? ")
	s.WriteString(confirmText)
	s.WriteString(m.styles.Label.Render("[Y/Enter] Yes  [N/Esc] Go Back"))

	return s.String()
}

func (m Model) renderPreview() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("📝 Content:  %s", truncateString(m.config.Content, 40)))
	lines = append(lines, fmt.Sprintf("📄 Format:   %s", strings.ToUpper(string(m.config.Format))))

	colorName := "Custom"
	for name, c := range config.PredefinedColors {
		if c == m.config.Foreground {
			colorName = name
			break
		}
	}
	lines = append(lines, fmt.Sprintf("🎨 Color:    %s (%s)", colorName, config.ColorToHex(m.config.Foreground)))
	lines = append(lines, fmt.Sprintf("📐 Size:     %dx%d pixels", m.config.Size, m.config.Size))
	lines = append(lines, fmt.Sprintf("💾 Output:   %s", truncateString(m.config.OutputPath, 40)))

	return strings.Join(lines, "\n")
}

func (m Model) renderCompleteStep() string {
	var s strings.Builder

	successBox := m.styles.Success.Render(fmt.Sprintf("✓ QR code generated successfully!\n\nSaved to:\n%s", m.successPath))
	s.WriteString(successBox)

	s.WriteString("\n\n")
	s.WriteString(m.styles.Label.Render("Press [R] to create another, [Q/Enter] to exit"))

	return s.String()
}

func (m Model) renderHelp() string {
	var help string

	switch m.step {
	case StepURL, StepSize, StepOutput:
		help = "Enter: Confirm • Esc: Back • Ctrl+C: Quit"
	case StepFormat:
		help = "←/→: Select • Enter/Space: Confirm • Esc: Back • Ctrl+C: Quit"
	case StepColor:
		if m.colorIndex == -1 {
			help = "Enter: Confirm • Esc: Cancel custom color • Ctrl+C: Quit"
		} else {
			help = "↑/↓: Select • Enter/Space: Confirm • C: Custom color • Esc: Back • Ctrl+C: Quit"
		}
	case StepConfirm:
		help = "Y/Enter: Generate • N/Esc: Go Back • Ctrl+C: Quit"
	case StepComplete:
		help = "R: Create another • Q/Enter: Exit"
	}

	return m.styles.Help.Render(help)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
