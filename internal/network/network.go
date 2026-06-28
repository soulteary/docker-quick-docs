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

package network

import (
	"bytes"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/soulteary/docker-quick-docs/internal/config"
)

var noiseReplacements = [][2]string{
	{"https://www.googletagmanager.com/gtag/js", "/keep-quiet.js"},
	{"https://www.google-analytics.com/analytics.js", "/keep-quiet.js"},
	{"https://hm.baidu.com/hm.js", "/keep-quiet.js"},
	{"https://cdn.matomo.cloud/", "/keep-quiet.js"},
	{"piwik.js", "/keep-quiet.js"},
	{"matomo.js", "/keep-quiet.js"},
	{"cnzz.com/z_stat.php", "/keep-quiet.js"},
}

func MuteNoise(body []byte) []byte {
	content := body
	for _, pair := range noiseReplacements {
		content = UpdateBody(content, []byte(pair[0]), []byte(pair[1]))
	}
	return content
}

func UpdateBody(content []byte, from []byte, to []byte) []byte {
	return bytes.ReplaceAll(content, from, to)
}

type UpdateJob struct {
	From string
	To   string
}

func PostProcess() fiber.Handler {
	return func(c fiber.Ctx) error {
		if err := c.Next(); err != nil {
			return err
		}

		body := c.Response().Body()
		if len(body) == 0 {
			return nil
		}

		mimeType := strings.ToLower(string(c.Response().Header.ContentType()))
		needUpdate := false

		var jobs []UpdateJob
		for _, rule := range config.PostRules {
			if strings.HasPrefix(mimeType, rule.Type) {
				if rule.Dir == "*" || strings.HasPrefix(c.Path(), rule.Dir) {
					needUpdate = true
					jobs = append(jobs, UpdateJob{From: rule.From, To: rule.To})
				}
			}
		}

		if !strings.HasPrefix(mimeType, "text/html") && !needUpdate {
			return nil
		}

		bodyUpdated := body
		if strings.HasPrefix(mimeType, "text/html") {
			bodyUpdated = MuteNoise(bodyUpdated)
		}
		for _, job := range jobs {
			bodyUpdated = UpdateBody(bodyUpdated, []byte(job.From), []byte(job.To))
		}

		if !bytes.Equal(body, bodyUpdated) {
			c.Response().SetBody(bodyUpdated)
		}
		return nil
	}
}
