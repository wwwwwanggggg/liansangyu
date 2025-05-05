package main

import (
	"liansangyu/config"
	"liansangyu/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(config.Config.AppMode)
	r := router.NewServer()

	err := r.RunTLS("0.0.0.0:8088", "server.crt", "server.key")
	if err != nil {
		panic(err)
	}
}
