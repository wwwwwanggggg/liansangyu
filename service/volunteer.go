package service

import (
	"errors"
	"fmt"
	"liansangyu/model"
	"sort"
	"time"

	"github.com/umahmood/haversine"
	"gorm.io/gorm"
)

type Volunteer struct {
}

type UVInfo struct {
	School string `json:"school" binding:"required"`
	Clazz  string `json:"clazz" binding:"required"`
	Skills string `json:"skills" binding:"required"`
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

// 作为接口暴露，完善个人信息之后使用
func (Volunteer) Register(Openid string, info UVInfo) error {
	u, err := getU(Openid)
	if err != nil {
		return err
	}

	if u.IsElder {
		return errors.New("老人身份和志愿者身份不能共存")
	}

	if u.IsVolunteer {
		return errors.New("此账号已经注册过志愿者了")
	}

	v := model.Volunteer{
		Openid: Openid,
		School: info.School,
		Clazz:  info.Clazz,
		Skills: info.Skills,
	}

	if err := model.DB.Create(&v).Error; err != nil {
		fmt.Println(err)
		return errors.New("创建志愿者信息失败")
	}

	if err := model.DB.Table("users").Where("openid = ?", Openid).Update("is_volunteer", true).Error; err != nil {
		fmt.Println(err)
		return errors.New("更新用户信息失败")
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
	_, err := getV(Openid)

	if err != nil {
		return err
	}

	t, err := getTask(taskID)
	if err != nil {
		return err
	}

	if t.Already >= t.Number {
		return errors.New("此任务人数已满")
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

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&model.TaskParticipants{
			Taskid:          int64(taskID),
			VolunteerOpenid: Openid,
		}).Error; err != nil {
			return errors.New("报名失败")
		}

		if err := tx.Model(&t).Update("already", t.Already+1).Error; err != nil {
			fmt.Println(err)
			return errors.New("更新任务失败")
		}
		return nil

	}); err != nil {
		return err
	}
	return nil
}

func (Volunteer) Signout(Openid string, taskID int) error {
	_, err := getV(Openid)
	if err != nil {
		return err
	}

	t, err := getTask(taskID)
	if err != nil {
		return err
	}

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

	if err := model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(
			&model.TaskParticipants{},
			"task_id = ? AND volunteer_openid = ?",
			taskID,
			Openid,
		).Error; err != nil {
			return fmt.Errorf("退出报名失败: %v", err)
		}
		if err := tx.Model(&model.Task{}).Where("id = ?", taskID).Update("already", t.Already-1).Error; err != nil {
			return errors.New("更新任务信息失败")
		}

		return nil

	}); err != nil {
		return err
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

func elderIsInOrganization(Openid string, OName string) (bool, error) {

	v, err := getE(Openid)

	if err != nil {
		return false, err
	}

	o, err := getO(OName)

	if err != nil {
		return false, err
	}

	index := Find(o.Elder, v, Eequals)
	if index != -1 {
		return true, nil
	}
	return false, nil
}

func (Volunteer) Join(Openid string, OName string) error {
	_, err := getV(Openid)
	if err != nil {
		return err
	}

	o, err := getO(OName)

	if o.Openid == Openid {
		return errors.New("您是这个组织的创建者")
	}
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

	if err := model.DB.Model(&model.OrganizationVolunteers{}).Create(&model.OrganizationVolunteers{
		OrganizationOpenid: o.Openid,
		VolunteerOpenid:    Openid,
	}).Error; err != nil {
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

	o, err := getO(OName)
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

	if err := model.DB.Model(&model.OrganizationVolunteers{}).Where("organization_openid = ? AND volunteer_openid = ?", o.Openid, Openid).Delete(&v); err != nil {
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

	task, err := getTask(taskID)
	if err != nil {
		return err
	}

	if !(t.After(*task.Starttime) && t.Before(task.Endtime.Add(time.Hour))) {
		return errors.New("您不在任务开始之后，结束一小时之内签到")
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

type Location struct {
	Longitide float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Resp struct {
	NearBy       []model.Task `json:"near_by"`
	Far          []model.Task `json:"far"`
	Full         []model.Task `json:"full"`
	Participated []model.Task `json:"participated"`
}

func (Resp) New() Resp {
	r := Resp{}
	r.NearBy = []model.Task{}
	r.Far = []model.Task{}
	r.Full = []model.Task{}
	r.Participated = []model.Task{}
	return r
}

type ByNumber []model.Task

func (h ByNumber) Len() int {
	return len(h)
}

func (h ByNumber) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// 从大到小
func (h ByNumber) Less(i, j int) bool {
	return h[i].Number > h[j].Number
}

func (Volunteer) GetTaskList(openid string, loc Location) (Resp, error) {
	v, err := getV(openid)
	if err != nil && err.Error() != "没有对应志愿者信息" {
		return Resp{}, err
	}

	var tasks []model.Task
	if err := model.DB.Find(&tasks).Error; err != nil {
		fmt.Println(err)
		return Resp{}, errors.New("查询任务失败")
	}
	r := Resp{}.New()
	sort.Sort(ByNumber(tasks))
	if loc.Latitude == 0 || loc.Longitide == 0 {
		r.Far = tasks
		return r, nil
	}

	for _, t := range tasks {
		if Find(t.Participants, v, Vequals) != -1 {
			r.Participated = append(r.Participated, t)
			continue
		}
		if t.Already >= t.Number {
			r.Full = append(r.Full, t)
			continue
		}
		_, km := haversine.Distance(haversine.Coord{
			Lat: loc.Latitude,
			Lon: loc.Longitide,
		}, haversine.Coord{
			Lat: t.Latitude,
			Lon: t.Longitude,
		})
		if km <= 5 {
			r.NearBy = append(r.NearBy, t)
		} else {
			r.Far = append(r.Far, t)
		}
	}

	return r, nil
}
