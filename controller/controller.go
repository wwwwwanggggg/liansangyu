package controller

type Controller struct {
	User
	Volunteer
	Task
	Elder
	Monitor
}

func New() *Controller {
	Controller := &Controller{}
	return Controller
}
