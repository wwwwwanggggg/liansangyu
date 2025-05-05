package controller

import (
	"errors"
	"liansangyu/common"

	"github.com/gin-gonic/gin"
)

type Monitor struct{}

func (Monitor) Add(c *gin.Context) {
	var info struct {
		ElderPhone string `json:"elder_phone" binding:"required,len=11,numeric"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err := srv.Monitor.Add(openid, info.ElderPhone)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
}

func (Monitor) Minus(c *gin.Context) {
	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err := srv.Monitor.Minus(openid)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
}
