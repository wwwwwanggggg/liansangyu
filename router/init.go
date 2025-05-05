package router

import (
	"liansangyu/config"
	"liansangyu/controller"

	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	r := gin.Default()
	config.SetCORS(r)
	config.InitSession(r)
	InitRouter(r)
	return r

}

var ctr = controller.New()
