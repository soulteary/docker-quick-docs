package server

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	static "github.com/soulteary/gin-static"
)

func InternalServer(host string, port int, dirRoot string, embedFS embed.FS, embedMode bool) {
	r := gin.New()

	if embedMode {
		r.NoRoute(static.ServeEmbed(dirRoot, embedFS))
	} else {
		r.Use(static.Serve("/", static.LocalFile(dirRoot, true)))
	}

	r.GET("/keep-quiet.js", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/javascript; charset=utf-8", []byte("console.log('Hello, world!')"))
	})

	r.Run(fmt.Sprintf("%s:%d", host, port))
}
