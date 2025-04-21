package service

import (
	"errors"
	"fmt"
	"liansangyu/model"

	"gorm.io/gorm"
)

type User struct{}

type UpdateUserInfo struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required,numeric,len=11"`
}

func (User) GetLevel(openid string) (uint8, error) {
	var user model.User

	if err := model.DB.Where("openid = ?", openid).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return 1, nil
	} else if err != nil {
		return 0, errors.New("查询用户信息失败")
	}

	return user.UserType, nil

}

func (User) Update(info UpdateUserInfo, openid string, userType uint8) error {
	var user model.User
	if err := model.DB.Where("openid = ?", openid).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有用户数据，创建用户

		ui := model.User{
			Name:  info.Name,
			Phone: info.Phone,

			Openid:   openid,
			UserType: userType,
		}

		if err := model.DB.Create(&ui).Error; err != nil {
			return errors.New("创建用户失败")
		}
	} else if err != nil {
		return errors.New("查询用户失败")
	} else {
		// 存在数据
		fmt.Println(user)
		if err := model.DB.Table("users").Where("openid = ?", openid).Updates(&info).Error; err != nil {
			return errors.New("更新用户失败")
		}
	}
	return nil
}
