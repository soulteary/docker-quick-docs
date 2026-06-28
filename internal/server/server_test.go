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

package server

import (
	"embed"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/soulteary/docker-quick-docs/internal/config"
	"github.com/stretchr/testify/require"
)

//go:embed testdata
var testEmbedFS embed.FS

func TestNewApp_RootShowsListing(t *testing.T) {
	t.Setenv("EMBED", "off")
	t.Setenv("INDEX", "off")
	docsDir := filepath.Join("..", "..", "docs")
	t.Setenv("DOCS", docsDir)

	app, err := NewApp(testEmbedFS)
	require.NoError(t, err)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	require.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, 200, resp.StatusCode)
	require.Contains(t, string(body), "README.md")
}

func TestNewApp_KeepQuietJS(t *testing.T) {
	t.Setenv("EMBED", "off")
	t.Setenv("DOCS", t.TempDir())

	app, err := NewApp(testEmbedFS)
	require.NoError(t, err)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/keep-quiet.js", nil))
	require.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, 200, resp.StatusCode)
	require.Contains(t, string(body), "console.log")
}

func TestNewApp_Health(t *testing.T) {
	t.Setenv("EMBED", "off")
	t.Setenv("DOCS", t.TempDir())

	app, err := NewApp(testEmbedFS)
	require.NoError(t, err)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/health", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

func TestNewApp_EmbedMode(t *testing.T) {
	t.Setenv("EMBED", "on")
	t.Setenv("DOCS", "testdata/docs")

	app, err := NewApp(testEmbedFS)
	require.NoError(t, err)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/hello.txt", nil))
	require.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, 200, resp.StatusCode)
	require.Contains(t, string(body), "hello from embed")
}

func TestNewApp_IndexEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	siteDir := filepath.Join(tmpDir, "site")
	require.NoError(t, os.MkdirAll(siteDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(siteDir, "index.html"), []byte("<html>home</html>"), 0o644))

	t.Setenv("EMBED", "off")
	t.Setenv("DOCS", tmpDir)
	t.Setenv("INDEX", "on")

	app, err := NewApp(testEmbedFS)
	require.NoError(t, err)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/site/", nil))
	require.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, 200, resp.StatusCode)
	require.Contains(t, string(body), "home")
}

func TestPostProcessIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs", "san")
	require.NoError(t, os.MkdirAll(docsDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(docsDir, "page.html"), []byte(
		`<html><head><script src="https://www.google-analytics.com/analytics.js"></script></head>`+
			`<body>Link: https://ecomfe.github.io/san/docs</body></html>`,
	), 0o644))

	configFile := filepath.Join(tmpDir, "config.json")
	require.NoError(t, os.WriteFile(configFile, []byte(
		`[{"from":"https://ecomfe.github.io/san/","to":"/san/","type":"html","dir":"/san/"}]`,
	), 0o644))

	t.Setenv("CONFIG", configFile)
	t.Setenv("EMBED", "off")
	t.Setenv("DOCS", filepath.Join(tmpDir, "docs"))
	t.Setenv("INDEX", "off")
	config.PostRules = nil
	config.GetConfig()

	app, err := NewApp(testEmbedFS)
	require.NoError(t, err)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/san/page.html", nil))
	require.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, 200, resp.StatusCode)
	require.Contains(t, string(body), "/keep-quiet.js")
	require.Contains(t, string(body), "/san/docs")
	require.NotContains(t, string(body), "google-analytics.com")
	require.NotContains(t, string(body), "ecomfe.github.io/san/")
}

func TestNewApp_HeadRequest(t *testing.T) {
	t.Setenv("EMBED", "off")
	t.Setenv("DOCS", filepath.Join("..", "..", "docs"))

	app, err := NewApp(testEmbedFS)
	require.NoError(t, err)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodHead, "/README.md", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

func TestLaunch_ServesHealth(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("EMBED", "off")
	t.Setenv("DOCS", tmpDir)
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("CONFIG", filepath.Join(tmpDir, "no-config.json"))

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())
	t.Setenv("PORT", strconv.Itoa(port))

	done := make(chan struct{})
	go func() {
		Launch(testEmbedFS)
		close(done)
	}()

	require.Eventually(t, func() bool {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/health", port))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 3*time.Second, 50*time.Millisecond)

	proc, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	require.NoError(t, proc.Signal(syscall.SIGTERM))

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Launch did not shut down in time")
	}
}
