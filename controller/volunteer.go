package controller

import (
	"errors"
	"liansangyu/common"
	"liansangyu/model"
	"liansangyu/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Volunteer struct{}

func (Volunteer) Update(c *gin.Context) {
	var info service.UpdateVInfo
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

	err := srv.Volunteer.Update(info, userSession.Openid)
	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusAccepted, ResponseNew(c, nil, "更新成功"))
}

func (Volunteer) SignUp(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(common.ErrNew(errors.New("请输入合法的id"), common.ParamErr))
		return
	}

	session := SessionGet(c, "user-session")
	userSession, ok := session.(UserSession)
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err = srv.Volunteer.SignUp(id, userSession.Openid)

	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusAccepted, ResponseNew(c, nil, "报名成功"))
}

func (Volunteer) SignOut(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(common.ErrNew(errors.New("请输入合法的id"), common.ParamErr))
		return
	}

	session := SessionGet(c, "user-session")
	userSession, ok := session.(UserSession)
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	err = srv.Volunteer.SignOut(id, userSession.Openid)

	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusAccepted, ResponseNew(c, nil, "退报成功"))
}

func (Volunteer) Checkin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(common.ErrNew(errors.New("请输入合法的id"), common.ParamErr))
		return
	}

	session := SessionGet(c, "user-session")
	userSession, ok := session.(UserSession)
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}
	t := time.Now()
	err = srv.Volunteer.Checkin(&t, userSession.Openid, id)

	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusAccepted, ResponseNew(c, nil, "签到成功"))
}

func (Volunteer) Checkout(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(common.ErrNew(errors.New("请输入合法的id"), common.ParamErr))
		return
	}

	session := SessionGet(c, "user-session")
	userSession, ok := session.(UserSession)
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}
	t := time.Now()
	err = srv.Volunteer.Checkout(&t, userSession.Openid, id)

	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusAccepted, ResponseNew(c, nil, "签退成功"))
}

func (Volunteer) GetTasks(c *gin.Context) {
	var loc service.Location

	if err := c.ShouldBindJSON(&loc); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}

	if !(loc.Latitude < 90 && loc.Latitude > -90 && loc.Longtitude < 180 && loc.Longtitude > -180) {
		c.Error(common.ErrNew(errors.New("请输入合法的经纬度"), common.ParamErr))
		return
	}

	session := SessionGet(c, "user-session")
	userSession, ok := session.(UserSession)
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	tasks, err := srv.Volunteer.GetTasks(userSession.Openid, loc)

	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusAccepted, ResponseNew(c, tasks, "获取任务成功"))
}

func (Volunteer) GetInfo(c *gin.Context) {
	session := SessionGet(c, "user-session")
	userSession, ok := session.(UserSession)
	if !ok {
		c.Error(common.ErrNew(errors.New("登录状态有问题"), common.AuthErr))
		return
	}

	u, t, err := srv.Volunteer.GetInfo(userSession.Openid)

	resp := struct {
		VolunteerInfo model.Volunteer `json:"volunteer_info"`
		Tasks         []model.Task    `json:"tasks"`
	}{
		VolunteerInfo: u,
		Tasks:         t,
	}

	if err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}

	c.JSON(http.StatusOK, ResponseNew(c, resp, "获取志愿者信息成功"))
}
