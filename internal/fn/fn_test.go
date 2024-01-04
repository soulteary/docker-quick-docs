package fn_test

import (
	"os"
	"testing"

	"github.com/soulteary/docker-quick-docs/internal/fn"
)

func TestGetPort(t *testing.T) {
	// Test with no environment variable set
	os.Unsetenv("PORT")
	port := fn.GetPort()
	if port != 8080 {
		t.Errorf("GetPort() without PORT set = %d; want 8080", port)
	}

	// Test with invalid integer value
	os.Setenv("PORT", "invalid")
	port = fn.GetPort()
	if port != 8080 {
		t.Errorf("GetPort() with invalid PORT = %d; want 8080", port)
	}

	// Test with valid integer but out of range value
	os.Setenv("PORT", "-1")
	port = fn.GetPort()
	if port != 8080 {
		t.Errorf("GetPort() with out-of-range PORT = %d; want 8080", port)
	}

	// Test with valid port number
	os.Setenv("PORT", "5000")
	port = fn.GetPort()
	if port != 5000 {
		t.Errorf("GetPort() with valid PORT = %d; want 5000", port)
	}
}

func TestIsEmbedMode(t *testing.T) {
	cases := []struct {
		env   string
		value bool
	}{
		{"on", true},
		{"ON", true},
		{"oN", true},
		{"off", false},
		{"OFF", false},
		{"invalid", false},
		{"", false},
		{"  on  ", true},
		{"on\n", true},
		{"\non\t", true},
	}

	for _, c := range cases {
		os.Setenv("EMBED", c.env)
		result := fn.IsEmbedMode()
		if result != c.value {
			t.Errorf("IsEmbedMode() with EMBED=%q = %v; want %v", c.env, result, c.value)
		}
	}
}

func TestFixResType(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType string
	}{
		{"HTML lower case", "html", "text/html"},
		{"HTML upper case", "HTML", "text/html"},
		{"CSS lower case", "css", "text/css"},
		{"CSS mixed case", "CsS", "text/css"},
		{"JS lower case", "js", "application/javascript"},
		{"JS upper case", "JS", "application/javascript"},
		{"JSON lower case", "json", "application/json"},
		{"JSON mixed case", "Json", "application/json"},
		{"Plain text", "txt", "txt"},
		{"Image", "jpg", "jpg"},
		{"Random string", "randomstring", "randomstring"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fn.FixResType(tt.input)
			if result != tt.expectedType {
				t.Errorf("FixResType(%v): expected %v, got %v", tt.input, tt.expectedType, result)
			}
		})
	}
}
