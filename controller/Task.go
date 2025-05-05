package controller

import (
	"errors"
	"liansangyu/common"
	"liansangyu/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Task struct{}

func (Task) New(c *gin.Context) {
	pt := c.Param("type")
	var info service.UpdateTaskInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err := srv.Task.New(openid, info, pt)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, ResponseNew(c, nil, "创建成功"))
}

func (Task) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(common.ErrNew(errors.New("请输入有效的id"), common.ParamErr))
		return
	}
	var info service.UpdateTaskInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err = srv.Task.Update(openid, id, info)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, nil, "更新成功"))
}

func (Task) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(common.ErrNew(errors.New("请输入有效的id"), common.ParamErr))
		return
	}

	openid, ok := GetOpenid(c, "user-session")
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err = srv.Task.Delete(openid, id)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, nil, "删除成功"))
}
