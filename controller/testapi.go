package controller

import (
	"errors"
	"liansangyu/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestAPI struct{}

func (TestAPI) DirectLogin(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Error(common.ErrNew(errors.New("没有code"), common.ParamErr))
		return
	}

	SessionSet(c, "user-session", UserSession{
		Openid: code,
		Level:  2,
	})

	c.JSON(http.StatusOK, ResponseNew(c, nil, "登录成功"))
}
