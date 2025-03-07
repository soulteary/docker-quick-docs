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

package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/soulteary/docker-quick-docs/internal/fn"
)

func ReadConfigFile(defaultFile string) ([]byte, error) {
	configFile := strings.TrimSpace(os.Getenv("CONFIG"))
	if configFile == "" {
		configFile = defaultFile
	}
	buf, err := os.ReadFile(configFile)
	if err != nil {
		if configFile != defaultFile {
			log.Println("读取配置文件失败，确认文件路径是否正确:", configFile)
		}
		return []byte(""), err
	}
	return buf, nil
}

type PostRule struct {
	From string `json:"from"`
	To   string `json:"to"`
	Dir  string `json:"dir,omitempty"`
	Type string `json:"type,omitempty"`
}

var PostRules []PostRule

func GetConfig() {
	buf, err := ReadConfigFile(DOCS_DEFAULT_CONFIG)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, &PostRules)
	if err != nil {
		return
	}

	if len(PostRules) > 0 {
		log.Println("解析配置文件成功，规则数量:", len(PostRules))
		var rules []PostRule
		for _, rule := range PostRules {
			// skip empty from rules
			if rule.From == "" {
				continue
			}
			// fill default values
			if rule.Dir == "" {
				rule.Dir = "*"
			}
			// fill default values
			if rule.Type == "" {
				rule.Type = "html"
			}
			rule.Type = fn.FixResType(rule.Type)
			rules = append(rules, rule)
		}
		PostRules = rules
	}
}
