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

func (User) Login(openid string) (model.User, error) {
	u, err := getU(openid)
	return u, err
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

	u.Openid = openid
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
type UserInfo struct {
	User         model.User         `json:"user"`
	Volunteer    model.Volunteer    `json:"volunteer"`
	Elder        model.Elder        `json:"elder"`
	Monitor      model.Monitor      `json:"monitor"`
	Organization model.Organization `json:"organization"`
}

func getM(openid string) (model.Monitor, error) {
	var m model.Monitor
	if err := model.DB.Where("openid = ?", openid).First(&m).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return m, errors.New("不存在监护人信息")
	} else if err != nil {
		fmt.Println(err)
		return m, errors.New("查询监护人信息出错")
	}
	return m, nil
}

func (User) Get(openid string) (UserInfo, error) {
	var res UserInfo
	u, err := getU(openid)
	if err != nil {
		return res, err
	}
	e, err := getE(openid)
	if err != nil && err.Error() != "没有相应老人信息" {
		return res, err
	}
	v, err := getV(openid)
	if err != nil && err.Error() != "没有对应志愿者信息" {
		return res, err
	}

	m, err := getM(openid)
	if err != nil && err.Error() != "不存在监护人信息" {
		return res, err
	}

	o, err := getOO(openid)
	if err != nil && err.Error() != "没有相应组织信息" {
		return res, err
	}

	res.User = u
	res.Volunteer = v
	res.Monitor = m
	res.Organization = o
	res.Elder = e

	return res, nil
}

type DecideInfo struct {
	Openid   string `json:"openid"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	InfoType string `json:"info_type"`
}

func (User) GetDecideList(openid string) ([]DecideInfo, error) {
	u, err := getU(openid)
	r := []DecideInfo{}
	if err != nil {
		return r, err
	}

	if u.IsElder {
		var m model.Monitor
		if err := model.DB.Model(&model.Monitor{}).Preload("User").Where("elder_openid = ?", u.Openid).First(&m).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			return r, errors.New("没有相应监护人数据")
		} else if err != nil {
			return r, errors.New("查询信息失败")
		}

		r = append(r, DecideInfo{
			Openid:   m.Openid,
			Phone:    m.User.Phone,
			Name:     m.User.Name,
			InfoType: "老人-监护人",
		})

		return r, nil
	} else if u.IsOrganization {
		var o model.Organization
		if err := model.DB.
			Model(&model.Organization{}).
			Preload("Volunteer").
			Preload("Elder").Where("openid = ?", openid).First(&o).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			return r, errors.New("没有相应组织信息")
		} else if err != nil {
			return r, errors.New("查询信息失败")
		}

		for _, v := range o.Volunteer {
			u, err := getU(v.Openid)
			if err != nil {
				return []DecideInfo{}, err
			}
			r = append(r, DecideInfo{
				Openid:   u.Openid,
				Phone:    u.Phone,
				Name:     u.Name,
				InfoType: "组织-志愿者",
			})
		}

		for _, e := range o.Elder {
			u, err := getU(e.Openid)

			if err != nil {
				return []DecideInfo{}, err
			}
			r = append(r, DecideInfo{
				Openid:   u.Openid,
				Phone:    u.Phone,
				Name:     u.Name,
				InfoType: "组织-老人",
			})
		}
	}
	return r, nil
}
