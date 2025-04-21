package service

import (
	"errors"
	"liansangyu/model"
	"time"
)

type Task struct{}

type NewTaskInfo struct {
	Title     string     `json:"title" binding:"required"`
	Starttime *time.Time `json:"start_time" binding:"required"`
	Endtime   *time.Time `json:"end_time" binding:"required"`

	Longtitude float64 `json:"longtitude" binding:"required"`
	Latitude   float64 `json:"latitude" binding:"required"`
	Desc       string  `json:"desc" binding:"required"`

	Number uint16 `json:"number" binding:"required"`
}

func (Task) New(info NewTaskInfo, openid string) error {
	if err := model.DB.Create(&model.Task{
		Title:      info.Title,
		Starttime:  info.Starttime,
		Endtime:    info.Endtime,
		Longtitude: info.Longtitude,
		Latitude:   info.Latitude,
		Desc:       info.Desc,
		Publisher:  openid,
		Number:     info.Number,
		Already:    0,
	}).Error; err != nil {
		return errors.New("创建任务失败")
	}

	return nil
}

type UpdateTaskInfo struct {
	Title string `json:"title" binding:"required"`

	Starttime *time.Time `json:"start_time" binding:"required"`
	Endtime   *time.Time `json:"end_time" binding:"required"`

	Longtitude float64 `json:"longtitude" binding:"required"`
	Latitude   float64 `json:"latitude" binding:"required"`

	Desc   string `json:"desc" binding:"required"`
	Number uint16 `json:"number" binding:"required"`
}

func (Task) Update(info UpdateTaskInfo, id int, openid string) error {
	var task model.Task
	if err := model.DB.Where("id = ?", id).First(&task).Error; err != nil {
		return errors.New("任务不存在")
	}

	if task.Publisher != openid {
		return errors.New("你不是该任务的发布者")
	}

	if task.Starttime.Before(time.Now()) {
		return errors.New("任务已开始，无法修改")
	}
	if task.Endtime.Before(time.Now()) {
		return errors.New("任务已结束，无法修改")
	}
	if task.Already > info.Number {
		return errors.New("任务参与人数已超过上限，无法修改")
	}

	if info.Starttime.Before(time.Now()) {
		return errors.New("任务开始时间不能早于当前时间")
	}
	if info.Endtime.Before(time.Now()) {
		return errors.New("任务结束时间不能早于当前时间")
	}

	if err := model.DB.Model(&task).Updates(model.Task{
		Title:      info.Title,
		Starttime:  info.Starttime,
		Endtime:    info.Endtime,
		Longtitude: info.Longtitude,
		Latitude:   info.Latitude,
		Desc:       info.Desc,
	}).Error; err != nil {
		return errors.New("更新任务失败")
	}
	return nil
}

func (Task) Delete(openid string, id int) error {
	var task model.Task
	if err := model.DB.Where("id = ?", id).First(&task).Error; err != nil {
		return errors.New("任务不存在")
	}

	if task.Publisher != openid {
		return errors.New("你不是该任务的发布者")
	}

	if err := model.DB.Delete(&task).Error; err != nil {
		return errors.New("删除任务失败")
	}
	return nil
}
