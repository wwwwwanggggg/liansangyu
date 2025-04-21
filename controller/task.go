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
	var info service.NewTaskInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	openid, ok := GetOpenid(c, "user-session")

	if !ok {
		c.Error(common.ErrNew(errors.New("登录装填有问题"), common.AuthErr))
		return
	}

	err := srv.Task.New(info, openid)

	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusAccepted, ResponseNew(c, nil, "创建成功"))
}

func (Task) Update(c *gin.Context) {
	var info service.UpdateTaskInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(common.ErrNew(errors.New("请输入合法的id"), common.ParamErr))
		return
	}

	openid, ok := GetOpenid(c, "user-session")

	if !ok {
		c.Error(common.ErrNew(errors.New("登录装填有问题"), common.AuthErr))
		return
	}

	err = srv.Task.Update(info, id, openid)

	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusAccepted, ResponseNew(c, nil, "更新成功"))

}

func (Task) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(common.ErrNew(errors.New("请输入合法的id"), common.ParamErr))
		return
	}

	openid, ok := GetOpenid(c, "user-session")

	if !ok {
		c.Error(common.ErrNew(errors.New("登录装填有问题"), common.AuthErr))
		return
	}

	err = srv.Task.Delete(openid, id)

	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusAccepted, ResponseNew(c, nil, "删除成功"))
}
