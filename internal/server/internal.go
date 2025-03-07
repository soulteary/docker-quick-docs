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
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	static "github.com/soulteary/gin-static"
)

func InternalServer(host string, port int, dirRoot string, embedFS embed.FS, embedMode bool) {
	r := gin.New()

	if embedMode {
		r.NoRoute(static.ServeEmbed(dirRoot, embedFS))
	} else {
		os.MkdirAll(dirRoot, os.ModePerm)
		r.Use(static.Serve("/", static.LocalFile(dirRoot, true)))
	}

	r.GET("/keep-quiet.js", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/javascript; charset=utf-8", []byte("console.log('Hello, world!')"))
	})

	r.Run(fmt.Sprintf("%s:%d", host, port))
}
