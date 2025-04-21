package service

type Service struct {
	Hello
	User
	Volunteer
	Task
	Elder
	Monitor
}

const (
	VOLUNTEER = iota + 1
	ELDER
	MONITOR
)

func New() *Service {
	service := &Service{}
	return service
}
