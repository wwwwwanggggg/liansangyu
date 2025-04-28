package service

import (
	"errors"
	"fmt"
	"liansangyu/model"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type User struct{}

type UpdateUserInfo struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"len=11,numeric"`
}

func getU(openid string) (model.User, error) {
	var u model.User
	if err := model.DB.Where("openid = ?", openid).First(&u).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return u, errors.New("没有相应用户信息")
	} else if err != nil {
		fmt.Println(err)
		return u, errors.New("查询用户信息失败")
	}

	return u, nil
}

// 完善基本信息
func (User) Register(openid string, info UpdateUserInfo) error {
	u, err := getU(openid)
	if err == nil {
		return errors.New("此微信号已经注册过用户了")
	} else if err.Error() != "没有相应用户信息" {
		return err
	}
	copier.Copy(&u, info)

	if err := model.DB.Create(&u).Error; err != nil {
		fmt.Println(err)
		return errors.New("创建用户失败")
	}
	return nil
}

func (User) Update(openid string, info UpdateUserInfo) error {
	_, err := getU(openid)
	if err != nil {
		return err
	}

	if err := model.DB.Table("users").Where("openid = ?", openid).Updates(&info).Error; err != nil {
		fmt.Println(err)
		return errors.New("更新失败")
	}

	return nil
}

/*
前端调用建议：

	这个返回所有相关的个人信息，建议前端写一个定时触发的函数，专门调用这个接口，
*/
func (User) Get(openid string) (u model.User, v model.Volunteer, m model.Monitor, o model.Organization, e error) {
	u, err := getU(openid)
	if err != nil {
		e = err
		return
	}
	db := model.DB.Model(&u)
	if u.IsVolunteer {
		db.Preload("Volunteer")
		vo, err := getV("openid")
		if err != nil {
			e = err
			return
		}
		v = vo
	}
	if u.IsElder {
		db.Preload("Elder")
	}
	if u.IsMonitor {
		db.Preload("Monitor")
	}
	if u.IsOrganization {
		db.Preload("Organization")
	}
	if err := db.First(&u).Error; err != nil {
		fmt.Println(err)
		e = errors.New("查询失败")
		return
	}
	return
}
