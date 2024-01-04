package main

import (
	"embed"
	"log"

	"github.com/soulteary/docker-quick-docs/internal/server"
	"github.com/soulteary/docker-quick-docs/internal/version"
)

//go:embed docs
var EmbedFS embed.FS

func main() {
	log.Println("Quick Docs", version.Version)
	server.Launch(EmbedFS)
}
