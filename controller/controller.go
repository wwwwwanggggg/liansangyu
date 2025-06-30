package controller

type Controller struct {
	User
	Volunteer
	Task
	Elder
	Monitor
	Organization
	TestAPI
}

func New() *Controller {
	Controller := &Controller{}
	return Controller
}
