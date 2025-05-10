package main

import (
	"liansangyu/config"
	"liansangyu/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(config.Config.AppMode)
	r := router.NewServer()

	if config.Config.Online != "no" {
		err := r.RunTLS("0.0.0.0:8088", "server.crt", "server.key")
		if err != nil {
			panic(err)
		}
	} else {
		err := r.Run("0.0.0.0:8088")
		if err != nil {
			panic(err)
		}
	}
}
