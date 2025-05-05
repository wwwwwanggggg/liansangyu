package service

import (
	"errors"
	"fmt"
	"liansangyu/model"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Elder struct{}

type UpdateEInfo struct {
	Disease   string  `json:"disease" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
}

func (Elder) Register(openid string, info UpdateEInfo) error {
	e, err := getE(openid)
	if err != nil && !(err.Error() == "没有相应老人信息") {
		return err
	} else if err == nil {
		return errors.New("您已经有老人账号了")
	}

	u, err := getU(openid)
	if err != nil {
		return err
	}

	if u.IsVolunteer || u.IsMonitor || u.IsOrganization {
		return errors.New("您不可以拥有老人身份")
	}

	copier.Copy(&e, &info)
	e.Openid = openid

	if err := model.DB.Create(&e).Error; err != nil {
		fmt.Println(err)
		return errors.New("创建信息失败")
	}

	if err := model.DB.Table("users").Where("openid = ?", openid).Update("is_elder", true).Error; err != nil {
		fmt.Println(err)
		return errors.New("更新用户信息失败")
	}

	return nil
}

func (Elder) Update(openid string, info UpdateEInfo) error {
	_, err := getE(openid)
	if err != nil {
		return err
	}

	if err := model.DB.Table("elders").Where("openid = ?", openid).Updates(&info).Error; err != nil {
		fmt.Println(err)
		return errors.New("更新数据失败")
	}
	return nil
}

func (Elder) Join(Openid string, OName string) error {
	v, err := getE(Openid)
	if err != nil {
		return err
	}
	isO, err := elderIsInOrganization(Openid, OName)
	if err != nil {
		return err
	}

	if isO {
		return errors.New("您已经加入这个组织了")
	}

	if err := model.DB.Model(&model.Organization{}).Association("Elder").Append(&v); err != nil {
		fmt.Println(err)
		return errors.New("加入组织失败")
	}

	return nil
}

func (Elder) Leave(Openid string, OName string) error {
	v, err := getE(Openid)
	if err != nil {
		return err
	}

	isO, err := elderIsInOrganization(Openid, OName)
	if err != nil {
		return err
	}

	if !isO {
		return errors.New("您没有加入这个组织,无法退出")
	}

	if err := model.DB.Model(&model.Organization{}).Association("Elder").Delete(&v); err != nil {
		fmt.Println(err)
		return errors.New("退出组织失败")
	}

	return nil
}

// 决定是否通过
func (Elder) Decide(openid string) error {
	_, err := getE(openid)
	if err != nil {
		return err
	}

	if err := model.DB.Table("monitors").Where("elder_openid = ?", openid).Update("passed", true).Error; err != nil {
		fmt.Println(err)
		return errors.New("更新用户信息失败")
	}

	return nil
}

func getE(openid string) (model.Elder, error) {
	var e model.Elder
	if err := model.DB.
		Where("openid = ?", openid).
		First(&e).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return e, errors.New("没有相应老人信息")
	} else if err != nil {
		fmt.Println(err)
		return e, errors.New("查询老人信息出错")
	}

	return e, nil
}
