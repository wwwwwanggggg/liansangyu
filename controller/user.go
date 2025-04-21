package controller

import (
	"errors"
	"liansangyu/common"
	"liansangyu/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct{}

func (User) Login(c *gin.Context) {
	var info struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	openid, err := code2session(info.Code)
	if err != nil {
		c.Error(common.ErrNew(err, common.WxErr))
		return
	}

	level, err := srv.User.GetLevel(openid)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return

	}

	SessionClear(c)
	SessionSet(c, "user-session", UserSession{
		Openid: openid,
		Level:  level,
	})

	c.JSON(http.StatusOK, ResponseNew(c, nil, "登录成功"))
}

func (User) Update(c *gin.Context) {
	var info service.UpdateUserInfo

	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	userType, err := strconv.Atoi(c.Param("type"))
	if err != nil {
		c.Error(common.ErrNew(errors.New("请输入有效的数字"), common.ParamErr))
		return
	}

	session := SessionGet(c, "user-session")
	userSession, ok := session.(UserSession)
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err = srv.User.Update(info, userSession.Openid, uint8(userType))
	if err != nil {
		common.ErrNew(err, common.SysErr)
		return
	}

	SessionClear(c)
	SessionSet(c, "user-session", UserSession{
		Openid: userSession.Openid,
		Level:  uint8(userType),
	})

	c.JSON(http.StatusAccepted, ResponseNew(c, nil, "更新成功"))
}
