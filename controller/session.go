package controller

import (
	"encoding/gob"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserSession struct {
	ID     int
	Openid string
	Level  uint8
}

func SessionGet(c *gin.Context, name string) any {
	session := sessions.Default(c)
	return session.Get(name)
}

func SessionSet(c *gin.Context, name string, body any) {
	session := sessions.Default(c)
	if body == nil {
		return
	}
	gob.Register(body)
	session.Set(name, body)

}

func SessionUpdate(c *gin.Context, name string, body any) {
	SessionSet(c, name, body)
}

func SessionClear(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
}

func SessionDelete(c *gin.Context, name string) {
	session := sessions.Default(c)
	session.Delete(name)
}

func GetOpenid(c *gin.Context, name string) (string, bool) {
	session := SessionGet(c, "user-session")
	userSession, ok := session.(UserSession)
	fmt.Println("getOpenid", userSession, ok)

	return userSession.Openid, ok
}
