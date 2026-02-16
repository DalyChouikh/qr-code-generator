package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/DalyChouikh/internal/config"
	"github.com/DalyChouikh/internal/generator"
	"github.com/DalyChouikh/internal/history"
	"github.com/DalyChouikh/internal/templates"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Step represents the current step in the wizard.
type Step int

const (
	StepContentType Step = iota
	StepURL
	StepTemplate
	StepFormat
	StepColor
	StepBgColor
	StepSize
	StepOutput
	StepConfirm
	StepComplete
)

const totalVisibleSteps = 9

// Model represents the application state.
type Model struct {
	styles *Styles
	config *config.QRConfig

	// Current step
	step Step

	// Input fields
	urlInput    textinput.Model
	sizeInput   textinput.Model
	outputInput textinput.Model
	colorInput  textinput.Model

	// Selection states
	formatIndex int // 0 = PNG, 1 = SVG
	colorIndex  int // Index in predefined colors, -1 for custom
	colorNames  []string

	// Content type selection
	contentTypes   []templates.ContentTypeInfo
	contentTypeIdx int
	templateWizard *TemplateWizard

	// Background color
	bgColorIndex int
	bgColorInput textinput.Model

	// File browser
	fileBrowserActive bool
	filePicker        FilePicker

	// QR terminal preview
	qrPreview string

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

	// Background color input (for custom hex)
	bgColorInput := textinput.New()
	bgColorInput.Placeholder = "#FFFFFF"
	bgColorInput.CharLimit = 7
	bgColorInput.Width = 46

	// Get home directory for default output path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	cfg.OutputPath = filepath.Join(homeDir, "qrcode.png")

	return Model{
		styles:         styles,
		config:         cfg,
		step:           StepContentType,
		urlInput:       urlInput,
		sizeInput:      sizeInput,
		outputInput:    outputInput,
		colorInput:     colorInput,
		formatIndex:    0,
		colorIndex:     0,
		colorNames:     config.GetPredefinedColorNames(),
		contentTypes:   templates.AvailableTypes(),
		contentTypeIdx: 0,
		bgColorIndex:   0,
		bgColorInput:   bgColorInput,
		filePicker:     NewFilePicker(),
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
			if m.step > StepContentType && m.step < StepComplete {
				// Let step handlers manage esc in sub-modes
				if m.step == StepColor && m.colorIndex == -1 {
					break
				}
				if m.step == StepBgColor && m.bgColorIndex == -1 {
					break
				}
				if m.step == StepOutput && m.fileBrowserActive {
					break
				}
				m.step = m.previousStep()
				m.err = nil
				return m, m.focusStep()
			}
		}

		// Step-specific handlers
		switch m.step {
		case StepContentType:
			return m.handleContentTypeStep(msg)
		case StepURL:
			return m.handleURLStep(msg)
		case StepTemplate:
			return m.handleTemplateStep(msg)
		case StepFormat:
			return m.handleFormatStep(msg)
		case StepColor:
			return m.handleColorStep(msg)
		case StepBgColor:
			return m.handleBgColorStep(msg)
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
	case StepTemplate:
		if m.templateWizard != nil {
			cmd = m.templateWizard.UpdateBlink(msg)
		}
	case StepSize:
		m.sizeInput, cmd = m.sizeInput.Update(msg)
	case StepOutput:
		if m.fileBrowserActive {
			cmd = m.filePicker.UpdateTextInput(msg)
		} else {
			m.outputInput, cmd = m.outputInput.Update(msg)
		}
	case StepColor:
		if m.colorIndex == -1 {
			m.colorInput, cmd = m.colorInput.Update(msg)
		}
	case StepBgColor:
		if m.bgColorIndex == -1 {
			m.bgColorInput, cmd = m.bgColorInput.Update(msg)
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
			m.step = StepBgColor
			m.colorInput.Blur()
			return m, nil
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
		m.step = StepBgColor
		return m, nil
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
	if m.fileBrowserActive {
		return m.handleOutputFileBrowser(msg)
	}

	switch msg.String() {
	case "tab":
		// Switch to file browser mode
		m.fileBrowserActive = true
		m.filePicker.Refresh()
		m.outputInput.Blur()
		return m, nil
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

func (m Model) handleOutputFileBrowser(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		// Switch back to manual input mode
		m.fileBrowserActive = false
		m.filePicker.Reset()
		return m, m.outputInput.Focus()
	case "esc":
		if m.filePicker.nameMode {
			// Exit filename input, back to browsing
			cmd := m.filePicker.Update(msg)
			return m, cmd
		}
		// Exit file browser, back to manual input
		m.fileBrowserActive = false
		m.filePicker.Reset()
		return m, m.outputInput.Focus()
	}

	cmd := m.filePicker.Update(msg)

	// Check if the file picker has confirmed a selection
	if m.filePicker.IsConfirmed() {
		path := m.filePicker.SelectedPath()
		m.config.SetOutputPath(path)
		m.err = nil
		m.step = StepConfirm
		m.fileBrowserActive = false
		m.filePicker.Reset()
		return m, nil
	}

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

		// Save to history
		if store, err := history.NewStore(); err == nil {
			_ = store.Add(history.Entry{
				Content:    m.config.Content,
				Format:     string(m.config.Format),
				Size:       m.config.Size,
				FgColor:    config.ColorToHex(m.config.Foreground),
				BgColor:    config.ColorToHex(m.config.Background),
				OutputPath: m.config.OutputPath,
			})
		}

		// Generate terminal preview for scanning
		if preview, err := generator.GenerateTerminalPreview(m.config.Content); err == nil {
			m.qrPreview = preview
		}

		m.step = StepComplete
		return m, nil
	case "n", "N":
		// Go back to content type step
		m.step = StepContentType
		m.fileBrowserActive = false
		return m, nil
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

func (m Model) handleContentTypeStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.contentTypeIdx > 0 {
			m.contentTypeIdx--
		}
	case "down", "j":
		if m.contentTypeIdx < len(m.contentTypes)-1 {
			m.contentTypeIdx++
		}
	case "enter", " ":
		ct := m.contentTypes[m.contentTypeIdx].Type
		m.err = nil
		switch ct {
		case templates.ContentURL, templates.ContentText:
			m.step = StepURL
			return m, m.urlInput.Focus()
		default:
			// WiFi, vCard, Email, SMS → template wizard
			tw := NewTemplateWizard(ct)
			m.templateWizard = &tw
			m.step = StepTemplate
			return m, textinput.Blink
		}
	}
	return m, nil
}

func (m Model) handleTemplateStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.templateWizard == nil {
		return m, nil
	}

	cmd := m.templateWizard.Update(msg)

	if m.templateWizard.IsConfirmed() {
		m.config.Content = m.templateWizard.Result()
		m.step = StepFormat
		m.err = nil
		return m, nil
	}

	return m, cmd
}

func (m Model) handleBgColorStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle custom color input if selected
	if m.bgColorIndex == -1 {
		switch msg.String() {
		case "enter":
			hexColor := strings.TrimSpace(m.bgColorInput.Value())
			if hexColor == "" {
				m.bgColorIndex = 0 // Reset to first predefined color
				return m, nil
			}
			bgColor, err := config.ParseHexColor(hexColor)
			if err != nil {
				m.err = fmt.Errorf("invalid hex color: %s", hexColor)
				return m, nil
			}
			m.config.Background = bgColor
			m.err = nil
			m.step = StepSize
			m.bgColorInput.Blur()
			return m, m.sizeInput.Focus()
		case "esc":
			m.bgColorIndex = 0
			m.bgColorInput.Blur()
			return m, nil
		}

		var cmd tea.Cmd
		m.bgColorInput, cmd = m.bgColorInput.Update(msg)
		return m, cmd
	}

	switch msg.String() {
	case "up", "k":
		if m.bgColorIndex > 0 {
			m.bgColorIndex--
		}
	case "down", "j":
		if m.bgColorIndex < len(m.colorNames)-1 {
			m.bgColorIndex++
		}
	case "c":
		m.bgColorIndex = -1
		m.bgColorInput.Focus()
		return m, textinput.Blink
	case "enter", " ":
		colorName := m.colorNames[m.bgColorIndex]
		m.config.Background = config.PredefinedColors[colorName]
		m.err = nil
		m.step = StepSize
		return m, m.sizeInput.Focus()
	}
	return m, nil
}

// previousStep returns the step to go back to.
func (m Model) previousStep() Step {
	switch m.step {
	case StepURL, StepTemplate:
		return StepContentType
	case StepFormat:
		ct := m.contentTypes[m.contentTypeIdx].Type
		if ct == templates.ContentURL || ct == templates.ContentText {
			return StepURL
		}
		return StepTemplate
	case StepColor:
		return StepFormat
	case StepBgColor:
		return StepColor
	case StepSize:
		return StepBgColor
	case StepOutput:
		return StepSize
	case StepConfirm:
		return StepOutput
	default:
		return m.step
	}
}

// stepDisplayNumber returns the visual step number for the progress bar.
func (m Model) stepDisplayNumber() int {
	switch m.step {
	case StepContentType:
		return 1
	case StepURL, StepTemplate:
		return 2
	case StepFormat:
		return 3
	case StepColor:
		return 4
	case StepBgColor:
		return 5
	case StepSize:
		return 6
	case StepOutput:
		return 7
	case StepConfirm:
		return 8
	case StepComplete:
		return 9
	default:
		return 0
	}
}

// focusStep focuses the text input for the current step.
func (m *Model) focusStep() tea.Cmd {
	switch m.step {
	case StepURL:
		return m.urlInput.Focus()
	case StepTemplate:
		if m.templateWizard != nil {
			m.templateWizard.focusFirst()
			return textinput.Blink
		}
	case StepSize:
		return m.sizeInput.Focus()
	case StepOutput:
		if m.fileBrowserActive {
			if m.filePicker.nameMode {
				return m.filePicker.nameInput.Focus()
			}
			return nil
		}
		return m.outputInput.Focus()
	case StepColor:
		if m.colorIndex == -1 {
			return m.colorInput.Focus()
		}
	case StepBgColor:
		if m.bgColorIndex == -1 {
			return m.bgColorInput.Focus()
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
		s.WriteString(m.styles.RenderProgressBar(m.stepDisplayNumber(), totalVisibleSteps))
		s.WriteString("\n\n")
	}

	// Current step content
	switch m.step {
	case StepContentType:
		s.WriteString(m.renderContentTypeStep())
	case StepURL:
		s.WriteString(m.renderURLStep())
	case StepTemplate:
		s.WriteString(m.renderTemplateStep())
	case StepFormat:
		s.WriteString(m.renderFormatStep())
	case StepColor:
		s.WriteString(m.renderColorStep())
	case StepBgColor:
		s.WriteString(m.renderBgColorStep())
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

func (m Model) renderContentTypeStep() string {
	var s strings.Builder

	s.WriteString(m.styles.Header.Render("Step 1: What do you want to encode?"))
	s.WriteString("\n\n")

	for i, ct := range m.contentTypes {
		var line string
		if i == m.contentTypeIdx {
			cursor := m.styles.OptionActive.Render("▸")
			name := m.styles.OptionActive.Render(ct.Icon + " " + ct.Name)
			desc := lipgloss.NewStyle().Foreground(primaryColor).Italic(true).Render(" — " + ct.Description)
			line = fmt.Sprintf("%s %s%s", cursor, name, desc)
		} else {
			cursor := m.styles.Option.Render(" ")
			name := m.styles.Option.Render(ct.Icon + " " + ct.Name)
			desc := lipgloss.NewStyle().Foreground(subtleColor).Italic(true).Render(" — " + ct.Description)
			line = fmt.Sprintf("%s %s%s", cursor, name, desc)
		}
		s.WriteString(line + "\n")
	}

	return s.String()
}

func (m Model) renderTemplateStep() string {
	var s strings.Builder

	ct := m.contentTypes[m.contentTypeIdx]
	s.WriteString(m.styles.Header.Render(fmt.Sprintf("Step 2: %s %s Details", ct.Icon, ct.Name)))
	s.WriteString("\n\n")

	if m.templateWizard != nil {
		s.WriteString(m.templateWizard.View(m.styles))
	}

	s.WriteString("\n\n")
	s.WriteString(m.styles.Label.Render("Tab/↓: Next field • Shift+Tab/↑: Previous • Enter: Confirm"))

	return s.String()
}

func (m Model) renderBgColorStep() string {
	var s strings.Builder

	s.WriteString(m.styles.Header.Render("Step 5: Background Color"))
	s.WriteString("\n\n")

	if m.bgColorIndex == -1 {
		// Custom color input mode
		s.WriteString(m.styles.LabelFocused.Render("Enter hex color:"))
		s.WriteString("\n")
		s.WriteString(m.styles.FocusedInput.Render(m.bgColorInput.View()))
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

			if i == m.bgColorIndex {
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

func (m Model) renderURLStep() string {
	var s strings.Builder

	s.WriteString(m.styles.Header.Render("Step 2: Enter URL or Text"))
	s.WriteString("\n\n")

	label := m.styles.LabelFocused.Render("Content:")
	s.WriteString(label + "\n")
	s.WriteString(m.styles.FocusedInput.Render(m.urlInput.View()))

	return s.String()
}

func (m Model) renderFormatStep() string {
	var s strings.Builder

	s.WriteString(m.styles.Header.Render("Step 3: Choose Output Format"))
	s.WriteString("\n\n")

	pngStyle := m.styles.Button
	svgStyle := m.styles.Button

	if m.formatIndex == 0 {
		pngStyle = m.styles.ButtonActive
	} else {
		svgStyle = m.styles.ButtonActive
	}

	pngBtn := pngStyle.Render("PNG")
	svgBtn := svgStyle.Render("SVG")
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, "  ", pngBtn, "    ", svgBtn))

	s.WriteString("\n\n")
	s.WriteString(m.styles.Label.Render("PNG: Raster image, best for most uses"))
	s.WriteString("\n")
	s.WriteString(m.styles.Label.Render("SVG: Vector format, scales infinitely"))

	return s.String()
}

func (m Model) renderColorStep() string {
	var s strings.Builder

	s.WriteString(m.styles.Header.Render("Step 4: Foreground Color"))
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

	s.WriteString(m.styles.Header.Render("Step 6: Set Dimensions"))
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

	s.WriteString(m.styles.Header.Render("Step 7: Output Location"))
	s.WriteString("\n\n")

	if m.fileBrowserActive {
		s.WriteString(m.filePicker.View(m.styles))
	} else {
		label := m.styles.LabelFocused.Render("Filename or path:")
		s.WriteString(label + "\n")
		s.WriteString(m.styles.FocusedInput.Render(m.outputInput.View()))
		s.WriteString("\n\n")
		s.WriteString(m.styles.Label.Render("Leave empty for 'qrcode' in current directory"))
		s.WriteString("\n")
		s.WriteString(m.styles.Label.Render("Use ~ for home directory, e.g., ~/Downloads/myqr"))
		s.WriteString("\n\n")
		s.WriteString(m.styles.Label.Render("Press Tab to browse files"))
	}

	return s.String()
}

func (m Model) renderConfirmStep() string {
	var s strings.Builder

	s.WriteString(m.styles.Header.Render("Step 8: Review & Generate"))
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

	// Show content type
	ct := m.contentTypes[m.contentTypeIdx]
	lines = append(lines, fmt.Sprintf("📋 Type:     %s %s", ct.Icon, ct.Name))
	lines = append(lines, fmt.Sprintf("📝 Content:  %s", truncateString(m.config.Content, 40)))
	lines = append(lines, fmt.Sprintf("📄 Format:   %s", strings.ToUpper(string(m.config.Format))))

	fgName := "Custom"
	for name, c := range config.PredefinedColors {
		if c == m.config.Foreground {
			fgName = name
			break
		}
	}
	lines = append(lines, fmt.Sprintf("🎨 FG Color: %s (%s)", fgName, config.ColorToHex(m.config.Foreground)))

	bgName := "Custom"
	for name, c := range config.PredefinedColors {
		if c == m.config.Background {
			bgName = name
			break
		}
	}
	lines = append(lines, fmt.Sprintf("🖼️  BG Color: %s (%s)", bgName, config.ColorToHex(m.config.Background)))

	lines = append(lines, fmt.Sprintf("📐 Size:     %dx%d pixels", m.config.Size, m.config.Size))
	lines = append(lines, fmt.Sprintf("💾 Output:   %s", truncateString(m.config.OutputPath, 40)))

	return strings.Join(lines, "\n")
}

func (m Model) renderCompleteStep() string {
	var s strings.Builder

	successBox := m.styles.Success.Render(fmt.Sprintf("✓ QR code generated successfully!\n\nSaved to:\n%s", m.successPath))
	s.WriteString(successBox)

	if m.qrPreview != "" {
		s.WriteString("\n\n")
		s.WriteString(m.styles.Header.Render("Scan with your phone:"))
		s.WriteString("\n\n")
		s.WriteString(m.qrPreview)
	}

	s.WriteString("\n")
	s.WriteString(m.styles.Label.Render("Press [R] to create another, [Q/Enter] to exit"))

	return s.String()
}

func (m Model) renderHelp() string {
	var help string

	switch m.step {
	case StepContentType:
		help = "↑/↓: Select • Enter/Space: Confirm • Ctrl+C: Quit"
	case StepURL, StepSize:
		help = "Enter: Confirm • Esc: Back • Ctrl+C: Quit"
	case StepTemplate:
		help = "Tab/↓: Next field • Shift+Tab/↑: Prev • Enter: Confirm • Esc: Back • Ctrl+C: Quit"
	case StepOutput:
		if m.fileBrowserActive {
			if m.filePicker.nameMode {
				help = "Enter: Confirm • Esc: Back to browsing • Ctrl+C: Quit"
			} else {
				help = "↑/↓: Navigate • Enter: Open/Select • N: New filename • ~: Home • Tab: Manual input • Ctrl+C: Quit"
			}
		} else {
			help = "Enter: Confirm • Tab: Browse files • Esc: Back • Ctrl+C: Quit"
		}
	case StepFormat:
		help = "←/→: Select • Enter/Space: Confirm • Esc: Back • Ctrl+C: Quit"
	case StepColor:
		if m.colorIndex == -1 {
			help = "Enter: Confirm • Esc: Cancel custom color • Ctrl+C: Quit"
		} else {
			help = "↑/↓: Select • Enter/Space: Confirm • C: Custom color • Esc: Back • Ctrl+C: Quit"
		}
	case StepBgColor:
		if m.bgColorIndex == -1 {
			help = "Enter: Confirm • Esc: Cancel custom color • Ctrl+C: Quit"
		} else {
			help = "↑/↓: Select • Enter/Space: Confirm • C: Custom color • Esc: Back • Ctrl+C: Quit"
		}
	case StepConfirm:
		help = "Y/Enter: Generate • N: Go Back • Esc: Previous step • Ctrl+C: Quit"
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
