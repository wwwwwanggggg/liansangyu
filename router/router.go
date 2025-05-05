package router

import (
	"liansangyu/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.Error)
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	apiRouter := r.Group("/api")
	{
		apiRouter.POST("/", ctr.User.Login)
		apiRouter.DELETE("/", ctr.User.Logout)

		ur := apiRouter.Group("user")
		{
			ur.POST("/", ctr.User.Register)
			ur.PUT("/", ctr.User.Update)
			ur.GET("/", ctr.User.Get)
		}

		tr := apiRouter.Group("task")
		tr.Use(middleware.CheckRole(1))
		{
			tr.POST("/:type", ctr.Task.New)
			tr.PUT("/:id", ctr.Task.Update)
			tr.DELETE("/:id", ctr.Task.Delete)
		}

		vr := apiRouter.Group("volunteer")
		vr.POST("/", ctr.Volunteer.Register)
		vr.Use(middleware.CheckRole(1))
		{
			vr.PUT("/", ctr.Volunteer.Update)
			vr.POST("/signin/:id", ctr.Volunteer.Signin)
			vr.DELETE("/signout/:id", ctr.Volunteer.Signout)
			vr.POST("/join", ctr.Volunteer.Join)
			vr.DELETE("/leave", ctr.Volunteer.Leave)
			vr.POST("/checkin/:id", ctr.Volunteer.Checkin)
			vr.DELETE("/checkout/:id", ctr.Volunteer.Checkout)
			vr.GET("/tasks", ctr.Volunteer.GetTaskList)
		}

		er := apiRouter.Group("elder")
		{
			er.POST("/", ctr.Elder.Register)
			er.PUT("/", ctr.Elder.Update)
			er.POST("/join", ctr.Elder.Join)
			er.DELETE("/leave", ctr.Elder.Leave)
			er.POST("/decide", ctr.Elder.Decide)
		}

		or := apiRouter.Group("organization")
		or.GET("/list", ctr.Organization.GetList)
		or.Use(middleware.CheckRole(1))
		{
			or.POST("/", ctr.Organization.Register)
			or.PUT("/", ctr.Organization.Update)
			or.POST("/decide", ctr.Organization.Decide)
			or.GET("/", ctr.Organization.Get)
		}

		mr := apiRouter.Group("/monitor")
		mr.Use(middleware.CheckRole(1))
		{
			mr.POST("/", ctr.Monitor.Add)
			mr.DELETE("/", ctr.Monitor.Minus)
		}

	}
}
