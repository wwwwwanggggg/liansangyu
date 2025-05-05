package controller

import (
	"errors"
	"fmt"
	"liansangyu/common"
	"liansangyu/service"
	"net/http"

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

	openid, err := code2openid(info.Code)
	if err != nil {
		c.Error(common.ErrNew(err, common.AuthErr))
		return
	}

	u, err := srv.User.Login(openid)
	if err != nil && err.Error() != "没有相应用户信息" {
		fmt.Println(err.Error())
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	SessionSet(c, "user-session", UserSession{
		Openid: openid,
		Level:  2,
	})

	c.JSON(http.StatusOK, ResponseNew(c, u, "登录成功"))
}

func (User) Register(c *gin.Context) {
	var info service.UpdateUserInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err := srv.User.Register(openid, info)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	c.JSON(http.StatusCreated, ResponseNew(c, nil, "注册成功"))
}

func (User) Update(c *gin.Context) {
	var info service.UpdateUserInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err := srv.User.Update(openid, info)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	c.JSON(http.StatusCreated, ResponseNew(c, nil, "更新成功"))
}

func (User) Get(c *gin.Context) {
	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	info, err := srv.User.Get(openid)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, info, "查询成功"))
}

func (User) Logout(c *gin.Context) {
	SessionDelete(c, "user-session")

	c.JSON(http.StatusOK, ResponseNew(c, nil, "退出登录"))
}
