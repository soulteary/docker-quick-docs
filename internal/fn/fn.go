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

package fn

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	defaultPort    = 8080
	defaultDocsDir = "docs"
	defaultHost    = "0.0.0.0"
)

func GetPort() int {
	portStr := os.Getenv("PORT")

	if portStr == "" {
		log.Printf("未设置环境变量 `PORT`，使用默认端口：%d\n", defaultPort)
		return defaultPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("环境变量 `PORT` 设置不正确，使用默认端口：%d\n", defaultPort)
		return defaultPort
	}

	if port < 1 || port > 65535 {
		log.Printf("环境变量 `PORT` 设置不正确，使用默认端口：%d\n", defaultPort)
		return defaultPort
	}
	return port
}

func GetDocsDir() string {
	docsDir := strings.TrimSpace(os.Getenv("DOCS"))
	if docsDir == "" {
		return defaultDocsDir
	}
	return docsDir
}

func GetHost() string {
	host := strings.TrimSpace(os.Getenv("HOST"))
	if host == "" {
		return defaultHost
	}
	return host
}

func IsEmbedMode() bool {
	return strings.ToLower(strings.TrimSpace(os.Getenv("EMBED"))) == "on"
}

func IsIndexEnabled() bool {
	return strings.ToLower(strings.TrimSpace(os.Getenv("INDEX"))) == "on"
}

func IndexNames() []string {
	if IsIndexEnabled() {
		return []string{"index.html"}
	}
	return []string{".__no_index__"}
}

func FixResType(typed string) string {
	typed = strings.ToLower(typed)
	switch typed {
	case "html":
		return "text/html"
	case "css":
		return "text/css"
	case "js":
		return "application/javascript"
	case "json":
		return "application/json"
	}
	return typed
}
