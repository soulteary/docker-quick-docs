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

package server

import (
	"embed"

	"github.com/gin-gonic/gin"
	"github.com/soulteary/docker-quick-docs/internal/config"
	"github.com/soulteary/docker-quick-docs/internal/fn"
	"github.com/soulteary/docker-quick-docs/internal/network"
)

func Launch(embedFS embed.FS) {
	isEmbedMode := fn.IsEmbedMode()
	publicPort := fn.GetPort()
	config.GetConfig()
	internalPort := publicPort - 1
	gin.SetMode(gin.ReleaseMode)
	forwarder := network.Forward(internalPort)
	go InternalServer(config.DOCS_INTERNAL_HOST, internalPort, config.DOCS_DIR_ROOT, embedFS, isEmbedMode)
	PublicServer(config.DOCS_PUBLIC_HOST, publicPort, forwarder)
}
