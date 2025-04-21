package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"liansangyu/model"
	"time"

	"github.com/umahmood/haversine"
	"gorm.io/gorm"
)

type Volunteer struct{}

type UpdateVInfo struct {
	School string   `json:"school" binding:"required"`
	Clazz  string   `json:"clazz" binding:"required"`
	Skills []string `json:"skills" binding:"required"`
}

func (Volunteer) Update(info UpdateVInfo, openid string) error {
	// 只在有user的情况下更新
	var user model.User
	if err := model.DB.Where("openid = ?", openid).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("没有相应用户信息")
	}

	var v model.Volunteer

	if err := model.DB.Where("openid = ?", openid).First(&v).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// 第一次完善志愿者信息
		m, _ := json.Marshal(info.Skills)
		v = model.Volunteer{
			Openid:    openid,
			School:    info.School,
			Clazz:     info.Clazz,
			Skills:    string(m),
			Hours:     0,
			Starttime: nil,
		}
		if err := model.DB.Create(&v).Error; err != nil {
			fmt.Println(err, openid, info.Skills)
			return errors.New("创建志愿者信息失败")
		}
	} else if err != nil {
		fmt.Println(err)
		return errors.New("查询志愿者信息失败")
	} else {
		// 更新即可
		if err := model.DB.Table("volunteers").Where("openid = ?", openid).Updates(&info).Error; err != nil {
			return errors.New("更新志愿者信息失败")
		}
	}
	return nil

}

func (Volunteer) SignUp(id int, openid string) error {
	var task model.Task

	if err := model.DB.Where("id = ?", id).First(&task).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("没有相应任务信息")
	} else if err != nil {
		return errors.New("查询任务信息失败")
	}

	if time.Now().After(*task.Endtime) {
		return errors.New("任务已经开始或者结束，无法报名")
	}

	var v model.Volunteer
	if err := model.DB.Where("openid = ?", openid).First(&v).Error; err != nil {
		fmt.Println(err, openid)
		return errors.New("查询志愿者信息失败")
	}

	if task.Already >= task.Number {
		return errors.New("报名人数已满")
	}
	if err := model.DB.Model(&task).Association("Participants").Append(&v); err != nil {
		return errors.New("报名失败")
	}
	if err := model.DB.Model(&task).Where("id = ?", id).Update("already", task.Already+1).Error; err != nil {
		return errors.New("更新任务人数失败")
	}
	return nil
}

func (Volunteer) SignOut(id int, openid string) error {
	var task model.Task

	if err := model.DB.Where("id = ?", id).First(&task).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("没有相应任务信息")
	} else if err != nil {
		return errors.New("查询任务信息失败")
	}

	if time.Now().After(*task.Starttime) {
		return errors.New("任务已经开始或者结束，无法退选")
	}

	var v model.Volunteer
	if err := model.DB.Where("openid = ?", openid).First(&v).Error; err != nil {
		fmt.Println(err)
		return errors.New("查询志愿者信息失败")
	}
	var sinfo model.TaskParticipants

	// 先判断是否已经报名
	if err := model.DB.Unscoped().Table("task_participants").Where("task_id = ? and volunteer_id = ?", id, v.ID).First(&sinfo).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("没有报名信息")
	} else if err != nil {
		fmt.Println(err)
		return errors.New("查询报名信息失败")
	}

	if err := model.DB.Model(&task).Association("Participants").Delete(&v); err != nil {
		fmt.Println(err)
		return errors.New("取消报名失败")
	}
	if err := model.DB.Model(&task).Where("id = ?", id).Update("already", task.Already-1).Error; err != nil {
		return errors.New("更新任务人数失败")
	}

	return nil
}

type TaskInfoResp struct {
	NearbyTasks     []model.Task
	Ntotal          int
	FullNumberTasks []model.Task
	Ftotal          int
	FarTasks        []model.Task
	Fartotal        int
}

type Location struct {
	Latitude   float64 `json:"latitude"`
	Longtitude float64 `json:"longtitude"`
}

func (Volunteer) GetTasks(openid string, loc Location) (TaskInfoResp, error) {
	res := TaskInfoResp{
		NearbyTasks:     []model.Task{},
		FullNumberTasks: []model.Task{},
		FarTasks:        []model.Task{},
	}

	if err := model.DB.Where("openid = ?", openid).First(&model.Volunteer{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return TaskInfoResp{}, errors.New("没有相应志愿者信息")
	} else if err != nil {
		return TaskInfoResp{}, errors.New("查询志愿者信息失败")
	}

	var tasks []model.Task
	if err := model.DB.Find(&tasks).Error; err != nil {
		return TaskInfoResp{}, errors.New("查询任务信息失败")
	}
	fmt.Println("location:", loc)
	for _, t := range tasks {
		if t.Already >= t.Number {
			res.FullNumberTasks = append(res.FullNumberTasks, t)
		} else if loc.Latitude == 0 || loc.Longtitude == 0 {
			res.FarTasks = append(res.FarTasks, t)
		} else {
			_, dis := haversine.Distance(
				haversine.Coord{Lat: loc.Latitude, Lon: loc.Longtitude},
				haversine.Coord{Lat: t.Latitude, Lon: t.Longtitude},
			)
			fmt.Println("dis:", dis)
			if dis < 5 {
				// 5 km nearest
				res.NearbyTasks = append(res.NearbyTasks, t)
			} else {
				res.FarTasks = append(res.FarTasks, t)
			}
		}
	}

	res.Ntotal = len(res.NearbyTasks)
	res.Ftotal = len(res.FullNumberTasks)
	res.Fartotal = len(res.FarTasks)
	return res, nil
}

// 检查志愿者是否报名了任务
func checkTaskParticipants(tid int, openid string) (v model.Volunteer, task model.Task, err error) {
	if er := model.DB.Where("openid = ?", openid).First(&v).Error; errors.Is(er, gorm.ErrRecordNotFound) {
		err = errors.New("没有相应志愿者信息")
		return
	} else if er != nil {
		fmt.Println(er)
		err = errors.New("查询志愿者信息失败")
		return
	}

	if er := model.DB.Where("id = ?", tid).First(&task).Error; errors.Is(er, gorm.ErrRecordNotFound) {
		err = errors.New("没有相应任务信息")
		return
	} else if er != nil {
		fmt.Println(err)
		err = errors.New("查询任务信息失败")
		return
	}

	var sinfo model.TaskParticipants

	if er := model.DB.Unscoped().Table("task_participants").Where("task_id = ? and volunteer_id = ?", tid, v.ID).First(&sinfo).Error; errors.Is(er, gorm.ErrRecordNotFound) {
		err = errors.New("没有报名信息")
		return
	} else if er != nil {

		err = errors.New("查询报名信息失败")
		return
	}
	return v, task, nil
}

func (Volunteer) Checkin(t *time.Time, openid string, taksID int) error {
	v, task, err := checkTaskParticipants(taksID, openid)

	if err != nil {
		return err
	}

	if !(t.After(*task.Starttime) && t.Before(*task.Endtime)) {
		return errors.New("签到时间不在任务时间范围内")
	}

	if err := model.DB.Model(&v).Update("start_time", t).Error; err != nil {
		return errors.New("签到失败")
	}

	return nil

}

func durationToUint16(dur time.Duration) uint16 {
	totalMinutes := dur.Minutes()                            // 总分钟数（float64）
	fullHours := int(totalMinutes / 60)                      // 完整的小时数
	remainingMinutes := totalMinutes - float64(fullHours*60) // 余下的分钟数

	if remainingMinutes >= 45 {
		return uint16(fullHours + 1) // 余数≥45分钟，加1小时
	} else {
		return uint16(fullHours) // 余数<45分钟，舍去
	}
}

func (Volunteer) Checkout(t *time.Time, openid string, taksID int) error {
	v, task, err := checkTaskParticipants(taksID, openid)
	if err != nil {
		return err
	}

	if !(t.After(*task.Starttime) && t.Before((*task.Endtime).Add(time.Hour))) {
		return errors.New("签到时间不在任务时间范围内")
	}

	if task.Endtime.Add(-15 * time.Minute).After(*t) {
		return errors.New("任务结束前15分钟内无法签退")
	}

	dur := task.Endtime.Sub(*t)
	hours := durationToUint16(dur)

	if err := model.DB.Model(&v).Update("hours", v.Hours+hours).Update("start_time", nil).Error; err != nil {
		return errors.New("更新失败")
	}

	return nil

}

func (Volunteer) GetInfo(openid string) (model.Volunteer, []model.Task, error) {
	var v model.Volunteer
	if err := model.DB.Where("openid = ?", openid).First(&v).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Volunteer{}, nil, errors.New("没有相应志愿者信息")
	} else if err != nil {
		return model.Volunteer{}, nil, errors.New("查询志愿者信息失败")
	}

	var tasks []model.Task
	if err := model.DB.Model(&v).Association("Tasks").Find(&tasks); err != nil {
		return model.Volunteer{}, nil, errors.New("查询任务信息失败")
	}

	return v, tasks, nil
}
