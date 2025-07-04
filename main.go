package main

import (
	"liansangyu/config"
	"liansangyu/router"

	"github.com/gin-gonic/gin"
)

var Router = router.NewServer()

func main() {
	gin.SetMode(config.Config.AppMode)

	if config.Config.Online != "no" {
		err := Router.RunTLS("0.0.0.0:8088", "server.crt", "server.key")
		if err != nil {
			panic(err)
		}
	} else {
		err := Router.Run("0.0.0.0:8088")
		if err != nil {
			panic(err)
		}
	}
}
