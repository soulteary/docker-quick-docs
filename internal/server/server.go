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
	"io/fs"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/soulteary/docker-quick-docs/internal/config"
	"github.com/soulteary/docker-quick-docs/internal/fn"
	"github.com/soulteary/docker-quick-docs/internal/network"
)

func NewApp(embedFS embed.FS) (*fiber.App, error) {
	isEmbedMode := fn.IsEmbedMode()
	docsDir := fn.GetDocsDir()

	app := fiber.New()
	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(network.PostProcess())

	app.Get("/keep-quiet.js", func(c fiber.Ctx) error {
		return c.Type("js").SendString("console.log('Hello, world!')")
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	var docFS fs.FS
	if isEmbedMode {
		sub, err := fs.Sub(embedFS, docsDir)
		if err != nil {
			return nil, fmt.Errorf("embed docs dir %q: %w", docsDir, err)
		}
		docFS = sub
	} else {
		if err := os.MkdirAll(docsDir, os.ModePerm); err != nil {
			return nil, err
		}
		docFS = os.DirFS(docsDir)
	}

	staticCfg := static.Config{
		Browse:     true,
		IndexNames: fn.IndexNames(),
		FS:         docFS,
	}
	staticHandler := static.New("", staticCfg)
	app.Get("/*", staticHandler)
	app.Head("/*", staticHandler)

	return app, nil
}

func Launch(embedFS embed.FS) {
	port := fn.GetPort()
	host := fn.GetHost()
	config.GetConfig()

	app, err := NewApp(embedFS)
	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("监听 http://%s\n", addr)

	go func() {
		if err := app.Listen(addr, fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
			log.Fatal(err)
		}
	}()

	fn.WaitForShutdown(app)
}
