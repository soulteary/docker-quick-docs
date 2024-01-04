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
		log.Println("请确认当前目录下存在 config.json 文件")
		log.Println(err)
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
