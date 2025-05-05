package controller

type Controller struct {
	User
	Volunteer
	Task
	Elder
	Monitor
	Organization
}

func New() *Controller {
	Controller := &Controller{}
	return Controller
}
