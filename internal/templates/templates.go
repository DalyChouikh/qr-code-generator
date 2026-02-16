// Package templates provides structured content templates for QR code generation.
//
// Instead of users manually formatting WiFi credentials, vCards, or mailto links,
// these templates provide a guided experience that generates the correct QR code
// encoding format automatically.
package templates

import (
	"fmt"
	"strings"
)

// ContentType represents the type of content to encode.
type ContentType int

const (
	ContentURL ContentType = iota
	ContentWiFi
	ContentVCard
	ContentEmail
	ContentSMS
	ContentText
)

// ContentTypeInfo holds display metadata for a content type.
type ContentTypeInfo struct {
	Type        ContentType
	Name        string
	Icon        string
	Description string
}

// AvailableTypes returns all available content types in display order.
func AvailableTypes() []ContentTypeInfo {
	return []ContentTypeInfo{
		{ContentURL, "URL", "🔗", "Website link"},
		{ContentWiFi, "WiFi", "📶", "WiFi network credentials"},
		{ContentVCard, "Contact", "👤", "Contact card (vCard)"},
		{ContentEmail, "Email", "✉️ ", "Email with subject & body"},
		{ContentSMS, "SMS", "💬", "Text message"},
		{ContentText, "Text", "📝", "Plain text"},
	}
}

// WiFiEncryption represents WiFi security types.
type WiFiEncryption string

const (
	WiFiWPA  WiFiEncryption = "WPA"
	WiFiWEP  WiFiEncryption = "WEP"
	WiFiNone WiFiEncryption = "nopass"
)

// WiFiData holds WiFi network information.
type WiFiData struct {
	SSID       string
	Password   string
	Encryption WiFiEncryption
	Hidden     bool
}

// Encode generates the QR code content string for WiFi.
// Format: WIFI:T:<encryption>;S:<ssid>;P:<password>;H:<hidden>;;
func (w *WiFiData) Encode() string {
	hidden := ""
	if w.Hidden {
		hidden = "H:true;"
	}

	password := ""
	if w.Encryption != WiFiNone {
		password = fmt.Sprintf("P:%s;", escapeWiFiField(w.Password))
	}

	return fmt.Sprintf("WIFI:T:%s;S:%s;%s%s;",
		w.Encryption,
		escapeWiFiField(w.SSID),
		password,
		hidden,
	)
}

// VCardData holds contact information.
type VCardData struct {
	FirstName    string
	LastName     string
	Phone        string
	Email        string
	Organization string
	Title        string
	URL          string
}

// Encode generates the QR code content string for vCard 3.0.
func (v *VCardData) Encode() string {
	var b strings.Builder

	b.WriteString("BEGIN:VCARD\r\n")
	b.WriteString("VERSION:3.0\r\n")

	fullName := strings.TrimSpace(v.FirstName + " " + v.LastName)
	if fullName != "" {
		b.WriteString(fmt.Sprintf("FN:%s\r\n", fullName))
		b.WriteString(fmt.Sprintf("N:%s;%s;;;\r\n", v.LastName, v.FirstName))
	}
	if v.Organization != "" {
		b.WriteString(fmt.Sprintf("ORG:%s\r\n", v.Organization))
	}
	if v.Title != "" {
		b.WriteString(fmt.Sprintf("TITLE:%s\r\n", v.Title))
	}
	if v.Phone != "" {
		b.WriteString(fmt.Sprintf("TEL;TYPE=CELL:%s\r\n", v.Phone))
	}
	if v.Email != "" {
		b.WriteString(fmt.Sprintf("EMAIL:%s\r\n", v.Email))
	}
	if v.URL != "" {
		b.WriteString(fmt.Sprintf("URL:%s\r\n", v.URL))
	}

	b.WriteString("END:VCARD\r\n")

	return b.String()
}

// EmailData holds email composition data.
type EmailData struct {
	Address string
	Subject string
	Body    string
}

// Encode generates a mailto: URI.
func (e *EmailData) Encode() string {
	var params []string

	if e.Subject != "" {
		params = append(params, "subject="+uriEncode(e.Subject))
	}
	if e.Body != "" {
		params = append(params, "body="+uriEncode(e.Body))
	}

	result := "mailto:" + e.Address
	if len(params) > 0 {
		result += "?" + strings.Join(params, "&")
	}

	return result
}

// SMSData holds SMS message data.
type SMSData struct {
	Phone   string
	Message string
}

// Encode generates an SMS URI.
// Format: smsto:<phone>:<message> or sms:<phone>?body=<message>
func (s *SMSData) Encode() string {
	if s.Message != "" {
		return fmt.Sprintf("smsto:%s:%s", s.Phone, s.Message)
	}
	return fmt.Sprintf("smsto:%s", s.Phone)
}

// escapeWiFiField escapes special characters in WiFi configuration fields.
func escapeWiFiField(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `;`, `\;`)
	s = strings.ReplaceAll(s, `:`, `\:`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

// uriEncode performs basic URI encoding for mailto parameters.
func uriEncode(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, " ", "%20")
	s = strings.ReplaceAll(s, "&", "%26")
	s = strings.ReplaceAll(s, "=", "%3D")
	s = strings.ReplaceAll(s, "#", "%23")
	s = strings.ReplaceAll(s, "\n", "%0A")
	return s
}

// WiFiEncryptionTypes returns available WiFi encryption options.
func WiFiEncryptionTypes() []struct {
	Type WiFiEncryption
	Name string
} {
	return []struct {
		Type WiFiEncryption
		Name string
	}{
		{WiFiWPA, "WPA/WPA2/WPA3"},
		{WiFiWEP, "WEP"},
		{WiFiNone, "None (Open)"},
	}
}
