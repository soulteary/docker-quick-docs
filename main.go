package main

import (
	"embed"
	"log"

	"github.com/soulteary/docker-quick-docs/internal/server"
)

//go:embed docs
var EmbedFS embed.FS

func main() {
	log.Println("Docker Quick Docs")
	server.Launch(EmbedFS)
}
