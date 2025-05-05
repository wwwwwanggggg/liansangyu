package service

import (
	"errors"
	"fmt"
	"liansangyu/model"

	"gorm.io/gorm"
)

type Monitor struct{}

func (Monitor) Add(openid string, elderPhone string) error {
	u, err := getU(openid)

	if u.IsElder {
		return errors.New("您不可以是监护人身份")
	} else if u.IsMonitor {
		return errors.New("您已经是监护人了")
	}

	if err != nil {
		return err
	}

	if u.IsMonitor {
		return errors.New("您已经是某个老人的监护人了")
	}

	var e model.User

	if err := model.DB.
		Where("phone = ?", elderPhone).
		First(&e).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("当前没有这个老人的信息")
	} else if err != nil {
		fmt.Println(err)
		return errors.New("查询老人信息出错")
	}

	if !e.IsElder {
		return errors.New("此号码对应的用户不是老人")
	}

	elder, err := getE(openid)

	if err != nil {
		return err
	}

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := model.DB.Model(&u).Update("is_monitor", true).Error; err != nil {
			fmt.Println(err)
			return errors.New("更新用户信息出错")
		}
		monitor := model.Monitor{
			Openid: openid,

			Elder:       elder,
			ElderOpenid: elder.Openid,
		}
		if err := model.DB.Model(&u).Association("monitor").Append(&monitor); err != nil {
			fmt.Println(err)
			return errors.New("更新用户信息出错")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (Monitor) Minus(openid string) error {
	u, err := getU(openid)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	if !u.IsMonitor {
		return errors.New("您不是监护人")
	}

	// if err := model.DB.Model(&model.Monitor{}).Association("elder").Delete(&e); err != nil {
	// 	fmt.Println(err)
	// 	return errors.New("接触关系绑定失败")
	// }

	return nil
}
