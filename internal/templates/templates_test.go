package templates

import (
	"errors"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		data Template
		want error
	}{
		{"wifi ok", &WiFiData{SSID: "Home", Encryption: WiFiWPA}, nil},
		{"wifi open network", &WiFiData{SSID: "Cafe", Encryption: WiFiNone}, nil},
		{"wifi missing ssid", &WiFiData{Encryption: WiFiWPA}, ErrSSIDRequired},
		{"wifi unknown encryption", &WiFiData{SSID: "Home", Encryption: "WPA9"}, ErrUnknownEncryption},
		{"contact first name only", &VCardData{FirstName: "Ada"}, nil},
		{"contact last name only", &VCardData{LastName: "Lovelace"}, nil},
		{"contact without name", &VCardData{Phone: "+1 555 0100"}, ErrContactNameRequired},
		{"email ok", &EmailData{Address: "a@example.com"}, nil},
		{"email missing address", &EmailData{Subject: "Hi"}, ErrEmailAddressRequired},
		{"sms ok", &SMSData{Phone: "+1 555 0100"}, nil},
		{"sms missing phone", &SMSData{Message: "Hi"}, ErrSMSPhoneRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.Validate(); !errors.Is(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name string
		data Template
		want string
	}{
		{
			"wifi escapes special characters",
			&WiFiData{SSID: `My;Net`, Password: `p:ss\`, Encryption: WiFiWPA},
			`WIFI:T:WPA;S:My\;Net;P:p\:ss\\;;`,
		},
		{
			"wifi open hidden network omits password",
			&WiFiData{SSID: "Cafe", Password: "ignored", Encryption: WiFiNone, Hidden: true},
			"WIFI:T:nopass;S:Cafe;H:true;;",
		},
		{
			"vcard",
			&VCardData{FirstName: "Ada", LastName: "Lovelace", Email: "ada@example.com"},
			"BEGIN:VCARD\r\nVERSION:3.0\r\nFN:Ada Lovelace\r\nN:Lovelace;Ada;;;\r\nEMAIL:ada@example.com\r\nEND:VCARD\r\n",
		},
		{
			"email encodes subject and body",
			&EmailData{Address: "a@example.com", Subject: "Hi there", Body: "a&b=c"},
			"mailto:a@example.com?subject=Hi%20there&body=a%26b%3Dc",
		},
		{"sms with message", &SMSData{Phone: "+15550100", Message: "hello"}, "smsto:+15550100:hello"},
		{"sms without message", &SMSData{Phone: "+15550100"}, "smsto:+15550100"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.Encode(); got != tt.want {
				t.Errorf("Encode() = %q, want %q", got, tt.want)
			}
		})
	}
}
