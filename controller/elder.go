package controller

import (
	"errors"
	"liansangyu/common"
	"liansangyu/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Elder struct{}

func (Elder) Update(c *gin.Context) {
	var info service.UpdateElderInfo
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

	err := srv.Elder.Update(info, userSession.Openid)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusAccepted, ResponseNew(c, nil, "更新成功"))

}

func (Elder) GetMonitor(c *gin.Context) {
	session := SessionGet(c, "user-session")
	userSession, ok := session.(UserSession)
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	monitor, err := srv.Elder.GetMonitor(userSession.Openid)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, monitor, "获取被监护人信息成功"))
}

func (Elder) Add(c *gin.Context) {

}
