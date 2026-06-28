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

package network_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/soulteary/docker-quick-docs/internal/config"
	"github.com/soulteary/docker-quick-docs/internal/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetPostRules(t *testing.T) {
	t.Helper()
	config.PostRules = nil
}

func TestMuteNoise(t *testing.T) {
	body := []byte(`<script src="https://www.googletagmanager.com/gtag/js"></script>` +
		`<script src="https://www.google-analytics.com/analytics.js"></script>` +
		`<script src="https://hm.baidu.com/hm.js"></script>` +
		`<script src="https://cdn.matomo.cloud/app.js"></script>` +
		`<script src="piwik.js"></script>` +
		`<script src="matomo.js"></script>` +
		`<script src="https://s4.cnzz.com/z_stat.php?id=123"></script>`)

	result := network.MuteNoise(body)
	assert.NotContains(t, string(result), "googletagmanager.com")
	assert.NotContains(t, string(result), "google-analytics.com")
	assert.NotContains(t, string(result), "hm.baidu.com")
	assert.NotContains(t, string(result), "cdn.matomo.cloud")
	assert.NotContains(t, string(result), "piwik.js")
	assert.NotContains(t, string(result), "matomo.js")
	assert.NotContains(t, string(result), "cnzz.com/z_stat.php")
	assert.Contains(t, string(result), "/keep-quiet.js")
}

func TestUpdateBody(t *testing.T) {
	content := []byte("hello world")
	assert.Equal(t, []byte("hello world"), network.UpdateBody(content, []byte("missing"), []byte("x")))

	replaced := network.UpdateBody([]byte("foo bar foo"), []byte("foo"), []byte("baz"))
	assert.Equal(t, []byte("baz bar baz"), replaced)
}

func testPostProcessApp(t *testing.T, handler fiber.Handler) *fiber.App {
	t.Helper()
	app := fiber.New()
	app.Use(network.PostProcess())
	app.Get("/test", handler)
	return app
}

func TestPostProcess_NoConfig_MutesHTML(t *testing.T) {
	resetPostRules(t)

	app := testPostProcessApp(t, func(c fiber.Ctx) error {
		return c.Type("html").SendString(`<script src="https://www.google-analytics.com/analytics.js"></script>`)
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/test", nil))
	require.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "/keep-quiet.js")
	assert.NotContains(t, string(body), "google-analytics.com")
}

func TestPostProcess_WithRules_ReplacesURL(t *testing.T) {
	resetPostRules(t)
	config.PostRules = []config.PostRule{
		{From: "https://example.com/docs/", To: "/local/", Type: "text/html", Dir: "*"},
	}

	app := testPostProcessApp(t, func(c fiber.Ctx) error {
		return c.Type("html").SendString(`<a href="https://example.com/docs/page">link</a>`)
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/test", nil))
	require.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `/local/page`)
	assert.NotContains(t, string(body), "example.com/docs")
}

func TestPostProcess_SkipsBinary(t *testing.T) {
	resetPostRules(t)
	original := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a}

	app := testPostProcessApp(t, func(c fiber.Ctx) error {
		c.Set("Content-Type", "image/png")
		return c.Send(original)
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/test", nil))
	require.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, original, body)
}

func TestPostProcess_DirScope(t *testing.T) {
	resetPostRules(t)
	config.PostRules = []config.PostRule{
		{From: "REMOTE", To: "LOCAL", Type: "text/html", Dir: "/scoped/"},
	}

	app := fiber.New()
	app.Use(network.PostProcess())
	app.Get("/scoped/page", func(c fiber.Ctx) error {
		return c.Type("html").SendString("REMOTE content")
	})
	app.Get("/other/page", func(c fiber.Ctx) error {
		return c.Type("html").SendString("REMOTE content")
	})

	scopedResp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/scoped/page", nil))
	require.NoError(t, err)
	scopedBody, _ := io.ReadAll(scopedResp.Body)
	assert.Contains(t, string(scopedBody), "LOCAL content")

	otherResp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/other/page", nil))
	require.NoError(t, err)
	otherBody, _ := io.ReadAll(otherResp.Body)
	assert.Contains(t, string(otherBody), "REMOTE content")
}
