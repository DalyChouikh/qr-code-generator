// Template wizard UI component for structured content input.
//
// This component provides multi-field forms for WiFi, vCard, Email, and SMS
// content types. Each template guides users through filling in the relevant
// fields, then generates the properly formatted QR code content string.
package ui

import (
	"strings"

	"github.com/DalyChouikh/internal/templates"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TemplateWizard manages the multi-field input forms for content templates.
type TemplateWizard struct {
	contentType templates.ContentType

	// WiFi fields
	wifiSSID     textinput.Model
	wifiPassword textinput.Model
	wifiEncIndex int  // 0=WPA, 1=WEP, 2=None
	wifiHidden   bool // Hidden network toggle

	// vCard fields
	vcardFirstName textinput.Model
	vcardLastName  textinput.Model
	vcardPhone     textinput.Model
	vcardEmail     textinput.Model
	vcardOrg       textinput.Model
	vcardTitle     textinput.Model
	vcardURL       textinput.Model

	// Email fields
	emailAddress textinput.Model
	emailSubject textinput.Model
	emailBody    textinput.Model

	// SMS fields
	smsPhone   textinput.Model
	smsMessage textinput.Model

	// Field navigation
	focusIndex int    // Which field is focused
	confirmed  bool   // User confirmed the form
	result     string // Encoded content string
}

// NewTemplateWizard creates a new wizard for the given content type.
func NewTemplateWizard(ct templates.ContentType) TemplateWizard {
	tw := TemplateWizard{contentType: ct}

	// WiFi
	tw.wifiSSID = newInput("MyNetwork", 64)
	tw.wifiPassword = newInput("password123", 128)

	// vCard
	tw.vcardFirstName = newInput("John", 64)
	tw.vcardLastName = newInput("Doe", 64)
	tw.vcardPhone = newInput("+1234567890", 20)
	tw.vcardEmail = newInput("john@example.com", 128)
	tw.vcardOrg = newInput("Acme Inc.", 128)
	tw.vcardTitle = newInput("Software Engineer", 128)
	tw.vcardURL = newInput("https://example.com", 256)

	// Email
	tw.emailAddress = newInput("user@example.com", 128)
	tw.emailSubject = newInput("Hello!", 256)
	tw.emailBody = newInput("I wanted to reach out...", 512)

	// SMS
	tw.smsPhone = newInput("+1234567890", 20)
	tw.smsMessage = newInput("Hello!", 256)

	// Focus the first field
	tw.focusFirst()

	return tw
}

func newInput(placeholder string, charLimit int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = charLimit
	ti.Width = 44
	return ti
}

// focusFirst focuses the first input field.
func (tw *TemplateWizard) focusFirst() {
	tw.focusIndex = 0
	tw.blurAll()

	switch tw.contentType {
	case templates.ContentWiFi:
		tw.wifiSSID.Focus()
	case templates.ContentVCard:
		tw.vcardFirstName.Focus()
	case templates.ContentEmail:
		tw.emailAddress.Focus()
	case templates.ContentSMS:
		tw.smsPhone.Focus()
	}
}

// blurAll blurs all input fields.
func (tw *TemplateWizard) blurAll() {
	tw.wifiSSID.Blur()
	tw.wifiPassword.Blur()
	tw.vcardFirstName.Blur()
	tw.vcardLastName.Blur()
	tw.vcardPhone.Blur()
	tw.vcardEmail.Blur()
	tw.vcardOrg.Blur()
	tw.vcardTitle.Blur()
	tw.vcardURL.Blur()
	tw.emailAddress.Blur()
	tw.emailSubject.Blur()
	tw.emailBody.Blur()
	tw.smsPhone.Blur()
	tw.smsMessage.Blur()
}

// fieldCount returns the number of fields for the current content type.
func (tw *TemplateWizard) fieldCount() int {
	switch tw.contentType {
	case templates.ContentWiFi:
		return 4 // SSID, Password, Encryption, Hidden
	case templates.ContentVCard:
		return 7 // First, Last, Phone, Email, Org, Title, URL
	case templates.ContentEmail:
		return 3 // Address, Subject, Body
	case templates.ContentSMS:
		return 2 // Phone, Message
	}
	return 0
}

// isToggleField returns true if the current field is a toggle (not a text input).
func (tw *TemplateWizard) isToggleField() bool {
	switch tw.contentType {
	case templates.ContentWiFi:
		return tw.focusIndex == 2 || tw.focusIndex == 3 // Encryption selector, Hidden toggle
	}
	return false
}

// focusField focuses the field at the current focusIndex.
func (tw *TemplateWizard) focusField() tea.Cmd {
	tw.blurAll()

	switch tw.contentType {
	case templates.ContentWiFi:
		switch tw.focusIndex {
		case 0:
			return tw.wifiSSID.Focus()
		case 1:
			return tw.wifiPassword.Focus()
			// 2 = encryption selector (no text input)
			// 3 = hidden toggle (no text input)
		}
	case templates.ContentVCard:
		switch tw.focusIndex {
		case 0:
			return tw.vcardFirstName.Focus()
		case 1:
			return tw.vcardLastName.Focus()
		case 2:
			return tw.vcardPhone.Focus()
		case 3:
			return tw.vcardEmail.Focus()
		case 4:
			return tw.vcardOrg.Focus()
		case 5:
			return tw.vcardTitle.Focus()
		case 6:
			return tw.vcardURL.Focus()
		}
	case templates.ContentEmail:
		switch tw.focusIndex {
		case 0:
			return tw.emailAddress.Focus()
		case 1:
			return tw.emailSubject.Focus()
		case 2:
			return tw.emailBody.Focus()
		}
	case templates.ContentSMS:
		switch tw.focusIndex {
		case 0:
			return tw.smsPhone.Focus()
		case 1:
			return tw.smsMessage.Focus()
		}
	}
	return nil
}

// Update handles key messages for the template wizard.
func (tw *TemplateWizard) Update(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "tab", "down":
		if tw.focusIndex < tw.fieldCount()-1 {
			tw.focusIndex++
			return tw.focusField()
		}
	case "shift+tab", "up":
		if tw.focusIndex > 0 {
			tw.focusIndex--
			return tw.focusField()
		}
	case "enter":
		// If on last field or any field, try to confirm
		if tw.tryEncode() {
			tw.confirmed = true
			return nil
		}
	case "left", "right", " ":
		// Handle toggle fields
		if tw.isToggleField() {
			tw.handleToggle(msg.String())
			return nil
		}
	}

	// Forward to current text input (for non-toggle fields)
	if !tw.isToggleField() {
		return tw.updateCurrentInput(msg)
	}
	return nil
}

// handleToggle handles left/right/space for toggle and selector fields.
func (tw *TemplateWizard) handleToggle(key string) {
	if tw.contentType == templates.ContentWiFi {
		if tw.focusIndex == 2 { // Encryption
			switch key {
			case "left":
				if tw.wifiEncIndex > 0 {
					tw.wifiEncIndex--
				}
			case "right":
				if tw.wifiEncIndex < 2 {
					tw.wifiEncIndex++
				}
			case " ":
				tw.wifiEncIndex = (tw.wifiEncIndex + 1) % 3
			}
		} else if tw.focusIndex == 3 { // Hidden
			tw.wifiHidden = !tw.wifiHidden
		}
	}
}

// updateCurrentInput forwards the key message to the currently focused text input.
func (tw *TemplateWizard) updateCurrentInput(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch tw.contentType {
	case templates.ContentWiFi:
		switch tw.focusIndex {
		case 0:
			tw.wifiSSID, cmd = tw.wifiSSID.Update(msg)
		case 1:
			tw.wifiPassword, cmd = tw.wifiPassword.Update(msg)
		}
	case templates.ContentVCard:
		switch tw.focusIndex {
		case 0:
			tw.vcardFirstName, cmd = tw.vcardFirstName.Update(msg)
		case 1:
			tw.vcardLastName, cmd = tw.vcardLastName.Update(msg)
		case 2:
			tw.vcardPhone, cmd = tw.vcardPhone.Update(msg)
		case 3:
			tw.vcardEmail, cmd = tw.vcardEmail.Update(msg)
		case 4:
			tw.vcardOrg, cmd = tw.vcardOrg.Update(msg)
		case 5:
			tw.vcardTitle, cmd = tw.vcardTitle.Update(msg)
		case 6:
			tw.vcardURL, cmd = tw.vcardURL.Update(msg)
		}
	case templates.ContentEmail:
		switch tw.focusIndex {
		case 0:
			tw.emailAddress, cmd = tw.emailAddress.Update(msg)
		case 1:
			tw.emailSubject, cmd = tw.emailSubject.Update(msg)
		case 2:
			tw.emailBody, cmd = tw.emailBody.Update(msg)
		}
	case templates.ContentSMS:
		switch tw.focusIndex {
		case 0:
			tw.smsPhone, cmd = tw.smsPhone.Update(msg)
		case 1:
			tw.smsMessage, cmd = tw.smsMessage.Update(msg)
		}
	}

	return cmd
}

// UpdateBlink handles non-key messages (cursor blink) for the active input.
func (tw *TemplateWizard) UpdateBlink(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch tw.contentType {
	case templates.ContentWiFi:
		switch tw.focusIndex {
		case 0:
			tw.wifiSSID, cmd = tw.wifiSSID.Update(msg)
		case 1:
			tw.wifiPassword, cmd = tw.wifiPassword.Update(msg)
		}
	case templates.ContentVCard:
		switch tw.focusIndex {
		case 0:
			tw.vcardFirstName, cmd = tw.vcardFirstName.Update(msg)
		case 1:
			tw.vcardLastName, cmd = tw.vcardLastName.Update(msg)
		case 2:
			tw.vcardPhone, cmd = tw.vcardPhone.Update(msg)
		case 3:
			tw.vcardEmail, cmd = tw.vcardEmail.Update(msg)
		case 4:
			tw.vcardOrg, cmd = tw.vcardOrg.Update(msg)
		case 5:
			tw.vcardTitle, cmd = tw.vcardTitle.Update(msg)
		case 6:
			tw.vcardURL, cmd = tw.vcardURL.Update(msg)
		}
	case templates.ContentEmail:
		switch tw.focusIndex {
		case 0:
			tw.emailAddress, cmd = tw.emailAddress.Update(msg)
		case 1:
			tw.emailSubject, cmd = tw.emailSubject.Update(msg)
		case 2:
			tw.emailBody, cmd = tw.emailBody.Update(msg)
		}
	case templates.ContentSMS:
		switch tw.focusIndex {
		case 0:
			tw.smsPhone, cmd = tw.smsPhone.Update(msg)
		case 1:
			tw.smsMessage, cmd = tw.smsMessage.Update(msg)
		}
	}

	return cmd
}

// tryEncode validates and encodes the form data. Returns true if valid.
func (tw *TemplateWizard) tryEncode() bool {
	switch tw.contentType {
	case templates.ContentWiFi:
		ssid := strings.TrimSpace(tw.wifiSSID.Value())
		if ssid == "" {
			return false
		}
		encTypes := templates.WiFiEncryptionTypes()
		data := &templates.WiFiData{
			SSID:       ssid,
			Password:   tw.wifiPassword.Value(),
			Encryption: encTypes[tw.wifiEncIndex].Type,
			Hidden:     tw.wifiHidden,
		}
		tw.result = data.Encode()
		return true

	case templates.ContentVCard:
		firstName := strings.TrimSpace(tw.vcardFirstName.Value())
		lastName := strings.TrimSpace(tw.vcardLastName.Value())
		if firstName == "" && lastName == "" {
			return false
		}
		data := &templates.VCardData{
			FirstName:    firstName,
			LastName:     lastName,
			Phone:        strings.TrimSpace(tw.vcardPhone.Value()),
			Email:        strings.TrimSpace(tw.vcardEmail.Value()),
			Organization: strings.TrimSpace(tw.vcardOrg.Value()),
			Title:        strings.TrimSpace(tw.vcardTitle.Value()),
			URL:          strings.TrimSpace(tw.vcardURL.Value()),
		}
		tw.result = data.Encode()
		return true

	case templates.ContentEmail:
		addr := strings.TrimSpace(tw.emailAddress.Value())
		if addr == "" {
			return false
		}
		data := &templates.EmailData{
			Address: addr,
			Subject: strings.TrimSpace(tw.emailSubject.Value()),
			Body:    strings.TrimSpace(tw.emailBody.Value()),
		}
		tw.result = data.Encode()
		return true

	case templates.ContentSMS:
		phone := strings.TrimSpace(tw.smsPhone.Value())
		if phone == "" {
			return false
		}
		data := &templates.SMSData{
			Phone:   phone,
			Message: strings.TrimSpace(tw.smsMessage.Value()),
		}
		tw.result = data.Encode()
		return true
	}

	return false
}

// IsConfirmed returns true if the form was submitted.
func (tw *TemplateWizard) IsConfirmed() bool {
	return tw.confirmed
}

// Result returns the encoded content string.
func (tw *TemplateWizard) Result() string {
	return tw.result
}

// View renders the template wizard form.
func (tw *TemplateWizard) View(styles *Styles) string {
	switch tw.contentType {
	case templates.ContentWiFi:
		return tw.viewWiFi(styles)
	case templates.ContentVCard:
		return tw.viewVCard(styles)
	case templates.ContentEmail:
		return tw.viewEmail(styles)
	case templates.ContentSMS:
		return tw.viewSMS(styles)
	}
	return ""
}

func (tw *TemplateWizard) viewWiFi(styles *Styles) string {
	var s strings.Builder

	s.WriteString(renderField(styles, "Network Name (SSID):", &tw.wifiSSID, tw.focusIndex == 0, true))
	s.WriteString(renderField(styles, "Password:", &tw.wifiPassword, tw.focusIndex == 1, false))

	// Encryption selector
	s.WriteString("\n")
	label := styles.Label
	if tw.focusIndex == 2 {
		label = styles.LabelFocused
	}
	s.WriteString(label.Render("Encryption:"))
	s.WriteString("\n")

	encTypes := templates.WiFiEncryptionTypes()
	var encBtns []string
	for i, enc := range encTypes {
		style := styles.Button
		if i == tw.wifiEncIndex {
			style = styles.ButtonActive
		}
		encBtns = append(encBtns, style.Render(enc.Name))
	}
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, encBtns...))

	// Hidden toggle
	s.WriteString("\n\n")
	label = styles.Label
	if tw.focusIndex == 3 {
		label = styles.LabelFocused
	}
	toggleStr := "○ No"
	if tw.wifiHidden {
		toggleStr = "● Yes"
	}
	s.WriteString(label.Render("Hidden Network: ") + label.Render(toggleStr))

	return s.String()
}

func (tw *TemplateWizard) viewVCard(styles *Styles) string {
	var s strings.Builder

	s.WriteString(renderField(styles, "First Name:", &tw.vcardFirstName, tw.focusIndex == 0, true))
	s.WriteString(renderField(styles, "Last Name:", &tw.vcardLastName, tw.focusIndex == 1, false))
	s.WriteString(renderField(styles, "Phone:", &tw.vcardPhone, tw.focusIndex == 2, false))
	s.WriteString(renderField(styles, "Email:", &tw.vcardEmail, tw.focusIndex == 3, false))
	s.WriteString(renderField(styles, "Organization:", &tw.vcardOrg, tw.focusIndex == 4, false))
	s.WriteString(renderField(styles, "Job Title:", &tw.vcardTitle, tw.focusIndex == 5, false))
	s.WriteString(renderField(styles, "Website:", &tw.vcardURL, tw.focusIndex == 6, false))

	return s.String()
}

func (tw *TemplateWizard) viewEmail(styles *Styles) string {
	var s strings.Builder

	s.WriteString(renderField(styles, "Email Address:", &tw.emailAddress, tw.focusIndex == 0, true))
	s.WriteString(renderField(styles, "Subject:", &tw.emailSubject, tw.focusIndex == 1, false))
	s.WriteString(renderField(styles, "Body:", &tw.emailBody, tw.focusIndex == 2, false))

	return s.String()
}

func (tw *TemplateWizard) viewSMS(styles *Styles) string {
	var s strings.Builder

	s.WriteString(renderField(styles, "Phone Number:", &tw.smsPhone, tw.focusIndex == 0, true))
	s.WriteString(renderField(styles, "Message:", &tw.smsMessage, tw.focusIndex == 1, false))

	return s.String()
}

// renderField renders a labeled text input field.
func renderField(styles *Styles, labelText string, input *textinput.Model, focused bool, first bool) string {
	var s strings.Builder

	if !first {
		s.WriteString("\n")
	}

	label := styles.Label
	inputStyle := styles.BlurredInput
	if focused {
		label = styles.LabelFocused
		inputStyle = styles.FocusedInput
	}

	s.WriteString(label.Render(labelText))
	s.WriteString("\n")
	s.WriteString(inputStyle.Render(input.View()))

	return s.String()
}
