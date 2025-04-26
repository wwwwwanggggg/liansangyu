package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"liansangyu/model"
	"time"

	"gorm.io/gorm"
)

type Volunteer struct {
}

type UVInfo struct {
	School string   `json:"school" binding:"required"`
	Clazz  string   `json:"clazz" binding:"required"`
	Skills []string `json:"skills" binding:"required"`
}

func getV(openid string) (model.Volunteer, error) {
	// 不可能有多个同Openid的Volunteer
	// Openid 引用自 User,User.Openid primaryKey

	var v model.Volunteer

	if err := model.DB.Where("openid = ?", openid).First(&v).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return v, errors.New("没有对应志愿者信息")
	} else if err != nil {
		fmt.Println(err)

		return v, errors.New("查询志愿者信息出错")
	}

	return v, nil
}

// 注册不作为接口暴露，由User处调用
func (Volunteer) Register(Openid string, info UVInfo) error {
	_, err := getV(Openid)
	if err != nil && !(err.Error() == "没有对应志愿者信息") {
		return err
	} else if err == nil {
		return errors.New("此微信号已经创建过志愿者身份了")
	}

	// 序列化json不会出错的
	j, _ := json.Marshal(info.Skills)
	v := model.Volunteer{
		Openid: Openid,
		School: info.School,
		Clazz:  info.Clazz,
		Skills: string(j),
	}

	if err := model.DB.Create(&v).Error; err != nil {
		fmt.Println(err)
		return errors.New("创建志愿者信息失败")
	}

	return nil
}

func (Volunteer) Update(Openid string, info UVInfo) error {
	_, err := getV(Openid)

	if err != nil {
		return err
	}

	if err := model.DB.Table("volunteers").Where("openid = ?", Openid).Updates(&info).Error; err != nil {
		fmt.Println(err)
		return errors.New("更新志愿者信息失败")
	}

	return nil
}

// 获取所有和志愿者，任务相关的信息
func (Volunteer) Get(Openid string) (u model.User, err error) {
	var user model.User

	if err := model.DB.
		Model(&model.User{}).
		Preload("Volunteer").
		Where("openid = ?", Openid).
		First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return u, errors.New("没有相应志愿者信息")
	} else if err != nil {
		fmt.Println(err)
		return u, errors.New("查询志愿者信息出错")
	}
	return user, nil
}

// 不允许注销
// func (Volunteer) Delete(Openid string) error {

// }

func getTask(taskID int) (mt model.Task, e error) {
	var t model.Task

	if err := model.DB.
		Preload("Participants").
		Where("id = ?", taskID).First(&t).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return mt, errors.New("没有相应的任务信息")
	} else if err != nil {
		fmt.Println(err)
		return mt, errors.New("查询任务信息失败")
	}

	return t, nil
}

func isParticipanted(Openid string, taskID int) (bool, error) {
	t, err := getTask(taskID)
	if err != nil {
		return false, err

	}

	v, err := getV(Openid)

	if err != nil {
		return false, err
	}

	index := Find(t.Participants, v, Vequals)
	if index != -1 {
		return true, nil
	}
	return false, nil
}

func (Volunteer) Signin(Openid string, taskID int) error {
	v, err := getV(Openid)

	if err != nil {
		return err
	}

	t, err := getTask(taskID)
	if err != nil {
		return err
	}

	isP, err := isParticipanted(Openid, taskID)
	if err != nil {
		return err
	}

	if isP {
		return errors.New("您已经报名或者参加过这个任务了，无法再次报名")
	}

	if time.Now().Add(-time.Hour).After(*t.Starttime) {
		return errors.New("任务即将开始或已经开始,无法报名")
	}

	if err := model.DB.Model(&Task{}).Association("Participants").Append(&v); err != nil {
		fmt.Println(err)
		return errors.New("报名任务失败")
	}
	return nil
}

func (Volunteer) Signout(Openid string, taskID int) error {
	v, err := getV(Openid)

	if err != nil {
		return err
	}

	isP, err := isParticipanted(Openid, taskID)
	if err != nil {
		return err
	}
	if !isP {
		return errors.New("您没有参加这个任务,无法退出报名")
	}

	if err := model.DB.Model(&model.Task{}).Association("Participants").Delete(&v); err != nil {
		fmt.Println(err)
		return errors.New("退出报名失败")
	}
	return nil
}

func isInOrganization(Openid string, OName string) (bool, error) {

	v, err := getV(Openid)

	if err != nil {
		return false, err
	}

	o, err := getO(OName)

	if err != nil {
		return false, err
	}

	index := Find(o.Volunteer, v, Vequals)
	if index != -1 {
		return true, nil
	}
	return false, nil
}

func (Volunteer) Join(Openid string, OName string) error {
	v, err := getV(Openid)
	if err != nil {
		return err
	}
	isO, err := isInOrganization(Openid, OName)
	if err != nil {
		return err
	}

	if isO {
		return errors.New("您已经加入这个组织了")
	}

	if err := model.DB.Model(&model.Organization{}).Association("Volunteer").Append(&v); err != nil {
		fmt.Println(err)
		return errors.New("加入组织失败")
	}

	return nil
}

func (Volunteer) Leave(Openid string, OName string) error {
	v, err := getV(Openid)
	if err != nil {
		return err
	}

	isO, err := isInOrganization(Openid, OName)
	if err != nil {
		return err
	}

	if !isO {
		return errors.New("您没有加入这个组织,无法退出")
	}

	if err := model.DB.Model(&model.Organization{}).Association("Volunteer").Delete(&v); err != nil {
		fmt.Println(err)
		return errors.New("退出组织失败")
	}

	return nil
}

func (Volunteer) Checkin(Openid string, taskID int) error {
	t := time.Now()

	v, err := getV(Openid)
	if err != nil {
		return err
	}
	isP, err := isParticipanted(Openid, taskID)
	if err != nil {
		return err
	}
	if !isP {
		return errors.New("您没有报名这个任务,无法签到")
	}

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&v).Update("start_time", t).Error; err != nil {
			fmt.Println(err)
			return errors.New("签到失败")
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func Timer(s, e time.Time) (int, error) {
	d := e.Sub(s)

	if d < 0 {
		return 0, errors.New("签退时间比签到时间还早?")
	}

	// 计算总分钟数
	totalMinutes := int(d.Minutes())

	// 计算完整的小时数
	hours := totalMinutes / 60

	// 计算剩余的分钟数
	remainingMinutes := totalMinutes % 60

	// 如果剩余分钟 > 45，则小时数 +1
	if remainingMinutes > 45 {
		hours += 1
	}

	// 返回 time.Duration 类型的小时数
	return hours, nil
}

func (Volunteer) Checkout(Openid string, taskID int) error {
	n := time.Now()

	v, err := getV(Openid)

	if err != nil {
		return err
	}

	t, err := getTask(taskID)
	if err != nil {
		return err
	}

	isp, err := isParticipanted(Openid, taskID)
	if err != nil {
		return err
	}

	if !isp {
		return errors.New("您没有参加这个任务，无法签退")
	}

	if n.After((*t.Endtime).Add(time.Hour)) || n.Before((*t.Endtime).Add(-15*time.Minute)) {
		return errors.New("您不在合理的签退时间之内")
	}

	d, err := Timer(*v.Starttime, *t.Endtime)
	if err != nil {
		return err
	}
	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&v).Update("start_time", nil).Update("hours", v.Hours+uint16(d)).Error; err != nil {
			fmt.Println(err)
			return errors.New("更新志愿者信息失败")
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (Volunteer) GetTaskList() {}
