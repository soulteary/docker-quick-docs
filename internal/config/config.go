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

package config

import (
	"bytes"
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

func unmarshalRules(buf []byte) ([]PostRule, error) {
	buf = bytes.TrimSpace(buf)
	var rules []PostRule
	if err := json.Unmarshal(buf, &rules); err == nil {
		return rules, nil
	}
	var single PostRule
	if err := json.Unmarshal(buf, &single); err != nil {
		return nil, err
	}
	log.Println("配置文件为单个对象，已自动包装为数组")
	return []PostRule{single}, nil
}

func GetConfig() {
	buf, err := ReadConfigFile(DOCS_DEFAULT_CONFIG)
	if err != nil {
		return
	}
	rules, err := unmarshalRules(buf)
	if err != nil {
		log.Println("解析配置文件失败:", err)
		PostRules = nil
		return
	}

	if len(rules) == 0 {
		PostRules = nil
		return
	}

	log.Println("解析配置文件成功，规则数量:", len(rules))
	var normalized []PostRule
	for _, rule := range rules {
		if rule.From == "" {
			log.Println("跳过无效规则: from 为空")
			continue
		}
		if rule.To == "" {
			log.Println("跳过无效规则: to 为空")
			continue
		}
		if rule.Dir == "" {
			rule.Dir = "*"
		}
		if rule.Type == "" {
			rule.Type = "html"
		}
		rule.Type = fn.FixResType(rule.Type)
		normalized = append(normalized, rule)
	}
	PostRules = normalized
}
