package controller

import (
	"errors"
	"liansangyu/common"
	"liansangyu/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Monitor struct{}

func (Monitor) Add(c *gin.Context) {
	var info service.UpdateMonitorInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	session := SessionGet(c, "user-session")
	userSession, ok := session.(UserSession)
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err := srv.Monitor.Add(info, userSession.Openid)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusAccepted, ResponseNew(c, nil, "添加成功"))
}

func (Monitor) DeMonitor(c *gin.Context) {
	session := SessionGet(c, "user-session")
	userSession, ok := session.(UserSession)
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err := srv.Monitor.DeMonitor(userSession.Openid)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusAccepted, ResponseNew(c, nil, "解绑成功"))
}
