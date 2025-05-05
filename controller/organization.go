package controller

import (
	"errors"
	"liansangyu/common"
	"liansangyu/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Organization struct {
}

func (Organization) Register(c *gin.Context) {
	var info service.UpdateOInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err := srv.Organization.Register(openid, info)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, nil, "成功注册"))
}

func (Organization) Update(c *gin.Context) {
	var info service.UpdateOInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err := srv.Organization.Update(openid, info)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, nil, "成功更新"))
}

func (Organization) Decide(c *gin.Context) {
	var info struct {
		Users []string `json:"users" binding:"required,dive"`
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

	err := srv.Organization.Decide(openid, info.Users)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, nil, "成功修改"))
}

func (Organization) Get(c *gin.Context) {
	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	o, err := srv.Organization.Get(openid)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, o, "查询成功"))
}

func (Organization) GetList(c *gin.Context) {
	data, err := srv.Organization.GetList()

	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, data, "查询成功"))
}
