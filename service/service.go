package service

import (
	"errors"
	"liansangyu/model"

	"gorm.io/gorm"
)

type Service struct {
	Hello
	User
	Volunteer
	Task
	Elder
	Monitor
	Organization
}

const (
	VOLUNTEER = iota + 1
	ELDER
	MONITOR
)

var Vequals = func(a, b model.Volunteer) bool { return a.Openid == b.Openid }

func New() *Service {
	service := &Service{}
	return service
}

func Find[T any](s []T, target T, equals func(a, b T) bool) (index int) {
	for i, v := range s {
		if equals(target, v) {
			index = i
			return
		}
	}
	return -1
}

func getO(OName string) (model.Organization, error) {
	var o model.Organization

	if err := model.DB.Preload("Admin").Preload("Volunteer").Preload("Elder").
		Where("name = ?", OName).
		First(&o).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return o, errors.New("没有相应组织信息")
	} else if err != nil {
		return o, errors.New("查找组织信息失败")
	}

	return o, nil
}
