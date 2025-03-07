/**
 * Copyright 2023-2025 Su Yang (soulteary)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
