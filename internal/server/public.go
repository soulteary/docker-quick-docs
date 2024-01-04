package server

import (
	"fmt"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func PublicServer(host string, port int, forwarder func(c *gin.Context)) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Any("/*path", forwarder)
	r.Run(fmt.Sprintf("%s:%d", host, port))
}
