package service

import (
	"errors"
	"fmt"
	"liansangyu/model"

	"gorm.io/gorm"
)

type Elder struct{}

type UpdateElderInfo struct {
	Disease string `json:"disease"`

	Longitude float64 `json:"longitude" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
}

func (Elder) Update(info UpdateElderInfo, openid string) error {
	// 只在有user的情况下更新

	var user model.User

	if err := model.DB.Where("openid = ?", openid).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("没有相应用户信息")
	} else if err != nil {
		fmt.Println(err)
		return errors.New("查询用户信息失败")
	}
	var e model.Elder

	if err := model.DB.Where("openid = ?", user.Openid).First(&e).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// 第一次完善志愿者信息
		e = model.Elder{
			Openid:    user.Openid,
			Disease:   info.Disease,
			Longitude: info.Longitude,
			Latitude:  info.Latitude,
		}
		if err := model.DB.Create(&e).Error; err != nil {
			return errors.New("创建老人信息失败")
		}
	} else if err != nil {
		return errors.New("查询老人信息失败")
	} else {
		// 更新即可
		if err := model.DB.Table("elders").Where("openid = ?", openid).Updates(&info).Error; err != nil {
			return errors.New("更新老人信息失败")
		}
	}
	return nil
}

func (Elder) GetMonitor(openid string) (model.Monitor, error) {

	var user model.User

	if err := model.DB.Where("openid = ?", openid).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Monitor{}, errors.New("没有相应用户信息")
	} else if err != nil {
		fmt.Println(err)
		return model.Monitor{}, errors.New("查询用户信息失败")
	}

	if user.UserType != ELDER {
		return model.Monitor{}, errors.New("您不是老人用户")
	}

	var m model.Monitor

	if err := model.DB.Where("monitor_phone = ?", user.Phone).First(&m).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Monitor{}, errors.New("您还没有绑定监护人")
	} else if err != nil {
		return model.Monitor{}, errors.New("查询监护人信息失败")
	}

	return m, nil

}

// 这里是拿到态度，true即是同意，false即是拒绝
func (Elder) Deal(openid string, atitude bool) error {
	m, err := Elder{}.GetMonitor(openid)
	if err != nil {
		return err
	}

	if m.Passed {
		return errors.New("您的监护人已通过审核，请勿重复添加")
	}

	if atitude {
		// 这里同意绑定此监护人
		if err := model.DB.Model(&m).Update("passed", true).Error; err != nil {
			return errors.New("确认失败")
		}
	} else {
		if err := model.DB.Model(&m).Update("monitor_phone", "").Update("passed", false).Error; err != nil {
			return errors.New("拒绝失败")
		}
	}

	return nil
}
