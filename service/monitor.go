package service

import (
	"errors"
	"fmt"
	"liansangyu/model"

	"gorm.io/gorm"
)

type Monitor struct{}

type UpdateMonitorInfo struct {
	MonitorPhone string `json:"monitor_phone" binding:"required,len=11,numeric"`
}

func (Monitor) Add(info UpdateMonitorInfo, openid string) error {
	var user model.User

	if err := model.DB.Where("phone = ?", info.MonitorPhone).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("没有相应被监护人信息")
	} else if err != nil {
		return errors.New("查询被监护人信息失败")
	}

	fmt.Println("ELDER:", ELDER, "user-type:", user.UserType)
	if user.UserType != ELDER {
		return errors.New("被监护人不是老人")
	}

	var m model.Monitor
	if err := model.DB.Where("openid = ?", openid).First(&m).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果没有记录，则创建新的记录
		m = model.Monitor{
			Openid:       openid,
			MonitorPhone: info.MonitorPhone,
		}

		if err := model.DB.Create(&m).Error; err != nil {
			return errors.New("创建记录失败")
		}
	} else if err != nil {
		return errors.New("查询记录失败")
	}

	if m.Passed {
		return errors.New("您已经绑定了被监护人，请先解绑")
	}

	if err := model.DB.Model(&m).Update("monitor_phone", info.MonitorPhone).Error; err != nil {
		return errors.New("更新记录失败")
	}
	return nil
}

func (Monitor) DeMonitor(openid string) error {
	var m model.Monitor

	if err := model.DB.Where("openid = ?", openid).First(&m).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("没有找到记录")
	} else if err != nil {
		return errors.New("查询记录失败")
	}

	if err := model.DB.Model(&m).Update("monitor_phone", "").Update("passed", false).Error; err != nil {
		return errors.New("解绑失败")
	}

	return nil
}
