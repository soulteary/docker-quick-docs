package main

import (
	"embed"
	"log"

	"github.com/soulteary/docker-quick-docs/internal/server"
)

//go:embed docs
var EmbedFS embed.FS
var Version = "dev"

func main() {
	log.Println("Docker Quick Docs", Version)
	server.Launch(EmbedFS)
}
