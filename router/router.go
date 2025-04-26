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
		apiRouter.Use(middleware.CheckRole(1))
		// userRouter := apiRouter.Group("user")
		// {
		// 	userRouter.POST("/", ctr.User.Login) // 登录

		// userRouter.Use(middleware.CheckRole(1))
		// 	userRouter.PUT("/:type", ctr.User.Update)

		// 	vRouter := userRouter.Group("/volunteer")
		// 	{
		// 		vRouter.PUT("/", ctr.Volunteer.Update)
		// 		vRouter.POST("/signup/:id", ctr.Volunteer.SignUp)     // 报名
		// 		vRouter.DELETE("/signout/:id", ctr.Volunteer.SignOut) // 退选
		// 		vRouter.POST("/in/:id", ctr.Volunteer.Checkin)
		// 		vRouter.POST("/out/:id", ctr.Volunteer.Checkout)
		// 		vRouter.GET("/tasks", ctr.Volunteer.GetTasks) // 获取任务
		// 		vRouter.GET("/", ctr.Volunteer.GetInfo)       // 获取志愿者信息
		// 	}
		// 	eRouter := userRouter.Group("/elder")
		// 	{
		// 		eRouter.PUT("/", ctr.Elder.Update)
		// 		eRouter.GET("/monitor", ctr.Elder.GetMonitor) // 获取被监护人信息
		// 		eRouter.GET("/")
		// 	}
		// 	mRouter := userRouter.Group("/monitor")
		// 	{
		// 		mRouter.PUT("/", ctr.Monitor.Add)
		// 		mRouter.DELETE("/", ctr.Monitor.DeMonitor)
		// 	}

		// }
		// tRouter := apiRouter.Group("task")
		// {
		// 	tRouter.Use(middleware.CheckRole(2))
		// 	tRouter.POST("/", ctr.Task.New) // 创建任务
		// 	tRouter.PUT("/:id", ctr.Task.Update)
		// 	tRouter.DELETE("/:id", ctr.Task.Delete)
		// }
	}
}
