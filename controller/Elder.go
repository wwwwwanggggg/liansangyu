package controller

import (
	"errors"
	"liansangyu/common"
	"liansangyu/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Elder struct{}

func (Elder) Register(c *gin.Context) {
	var info service.UpdateEInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err := srv.Elder.Register(openid, info)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	c.JSON(http.StatusCreated, ResponseNew(c, nil, "注册成功"))
}

func (Elder) Update(c *gin.Context) {
	var info service.UpdateEInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err := srv.Elder.Update(openid, info)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	c.JSON(http.StatusCreated, ResponseNew(c, nil, "更新成功"))
}

func (Elder) Join(c *gin.Context) {
	var info struct {
		OrganizationName string `json:"organization_name" binding:"required"`
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

	err := srv.Elder.Join(openid, info.OrganizationName)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, nil, "加入成功"))
}

func (Elder) Leave(c *gin.Context) {
	var info struct {
		OrganizationName string `json:"organization_name" binding:"required"`
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

	err := srv.Elder.Leave(openid, info.OrganizationName)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, nil, "加入成功"))
}

func (Elder) Decide(c *gin.Context) {
	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err := srv.Elder.Decide(openid)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, nil, "更新成功"))
}
