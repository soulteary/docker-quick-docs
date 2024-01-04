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
