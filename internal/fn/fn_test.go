/**
 * Copyright 2024-2026 Su Yang (soulteary)
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
	os.Unsetenv("PORT")
	port := fn.GetPort()
	if port != 8080 {
		t.Errorf("GetPort() without PORT set = %d; want 8080", port)
	}

	os.Setenv("PORT", "invalid")
	port = fn.GetPort()
	if port != 8080 {
		t.Errorf("GetPort() with invalid PORT = %d; want 8080", port)
	}

	os.Setenv("PORT", "-1")
	port = fn.GetPort()
	if port != 8080 {
		t.Errorf("GetPort() with out-of-range PORT = %d; want 8080", port)
	}

	os.Setenv("PORT", "5000")
	port = fn.GetPort()
	if port != 5000 {
		t.Errorf("GetPort() with valid PORT = %d; want 5000", port)
	}
}

func TestGetDocsDir(t *testing.T) {
	os.Unsetenv("DOCS")
	if got := fn.GetDocsDir(); got != "docs" {
		t.Errorf("GetDocsDir() = %q; want docs", got)
	}

	os.Setenv("DOCS", "  my-docs  ")
	if got := fn.GetDocsDir(); got != "my-docs" {
		t.Errorf("GetDocsDir() = %q; want my-docs", got)
	}
}

func TestGetHost(t *testing.T) {
	os.Unsetenv("HOST")
	if got := fn.GetHost(); got != "0.0.0.0" {
		t.Errorf("GetHost() = %q; want 0.0.0.0", got)
	}

	os.Setenv("HOST", "127.0.0.1")
	if got := fn.GetHost(); got != "127.0.0.1" {
		t.Errorf("GetHost() = %q; want 127.0.0.1", got)
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

func TestIsIndexEnabled(t *testing.T) {
	cases := []struct {
		env   string
		value bool
	}{
		{"on", true},
		{"ON", true},
		{"off", false},
		{"", false},
	}

	for _, c := range cases {
		os.Setenv("INDEX", c.env)
		if got := fn.IsIndexEnabled(); got != c.value {
			t.Errorf("IsIndexEnabled() with INDEX=%q = %v; want %v", c.env, got, c.value)
		}
	}
}

func TestIndexNames(t *testing.T) {
	os.Setenv("INDEX", "off")
	names := fn.IndexNames()
	if len(names) != 1 || names[0] != ".__no_index__" {
		t.Errorf("IndexNames() off = %v; want [.__no_index__]", names)
	}

	os.Setenv("INDEX", "on")
	names = fn.IndexNames()
	if len(names) != 1 || names[0] != "index.html" {
		t.Errorf("IndexNames() on = %v; want [index.html]", names)
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
