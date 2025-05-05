package service

import (
	"errors"
	"fmt"
	"liansangyu/model"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Organization struct{}

type UpdateOInfo struct {
	Name string `json:"name" binding:"required"`
	Logo string `json:"logo"`
}

func (Organization) Register(openid string, info UpdateOInfo) error {
	u, err := getU(openid)
	if err != nil {
		return err
	}

	if u.IsElder {
		return errors.New("您不能创建组织")
	}

	if u.IsOrganization {
		return errors.New("您最多是一个组织的管理员")
	}

	var o model.Organization
	copier.Copy(&o, &info)
	o.Openid = openid
	if err := model.DB.Create(&o).Error; err != nil {
		fmt.Println(err)
		return errors.New("创建失败")
	}

	if err := model.DB.Table("users").Where("openid = ?", openid).Update("is_organization", true).Error; err != nil {
		fmt.Println(err)
		return errors.New("更新用户信息失败")
	}

	return nil

}

func (Organization) Update(openid string, info UpdateOInfo) error {
	_, err := getU(openid)
	if err != nil {
		return err
	}

	if err := model.DB.Table("organizations").Where("openid = ?", openid).Updates(&info).Error; err != nil {
		fmt.Println(err)
		return errors.New("更新组织信息失败")
	}

	return nil
}

func (Organization) Get(openid string) (model.Organization, error) {
	var o model.Organization

	u, err := getU(openid)
	if err != nil {
		return o, err
	}

	if !u.IsOrganization {
		return o, errors.New("您不是一个组织的最高管理员")
	}

	if err := model.DB.Model(&model.Organization{}).
		Where("openid = ?", openid).
		Preload("Elder").Preload("Admin").Preload("Volunteer").
		First(&o).Error; err != nil {
		fmt.Println(err)
		return o, errors.New("查询组织信息失败")
	}

	return o, nil

}

// 决定通过志愿者，管理员，老人
func (Organization) Decide(openid string, query []string) error {
	u, err := getU(openid)
	if err != nil {
		return err
	}
	if !u.IsOrganization {
		return errors.New("您不是一个组织的最高管理员")
	}

	_, err = getOO(openid)
	if err != nil {
		return err
	}

	if err := model.DB.
		Table("organization_volunteers").
		Where("organization_openid = ? AND volunteer_openid IN ?", openid, query).
		Update("passed", true).Error; err != nil {
		fmt.Println(err)
		return errors.New("更新用户信息失败")
	}
	if err := model.DB.
		Table("organization_admins").
		Where("organization_openid = ? AND admin_openid IN ?", openid, query).
		Update("passed", true).Error; err != nil {
		fmt.Println(err)
		return errors.New("更新用户信息失败")
	}

	return nil
}

func (Organization) GetList() ([]model.Organization, error) {
	var d []model.Organization

	if err := model.DB.Model(&model.Organization{}).Preload("Admin").Preload("Volunteer").Preload("Elder").Find(&d).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return d, errors.New("查询组织信息失败")
	}

	return d, nil
}
