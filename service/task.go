package service

import (
	"errors"
	"fmt"
	"liansangyu/model"
	"time"

	"github.com/jinzhu/copier"
)

type Task struct{}

type UpdateTaskInfo struct {
	Title string `json:"title" binding:"required"`

	Starttime *time.Time `json:"start_time" binding:"required"`
	Endtime   *time.Time `json:"end_time" binding:"required"`
	Longitude float64    `json:"longitude" binding:"required"`
	Latitude  float64    `json:"latitude" binding:"required"`

	Desc   string `json:"desc" binding:"required"`
	Number uint16 `json:"number" binding:"required"`
}

func (Task) New(openid string, info UpdateTaskInfo) error {
	v, err := getV(openid)
	if err != nil {
		return err
	}

	var t model.Task

	copier.Copy(&t, &info)

	t.Publisher = v.Openid
	if err := model.DB.Model(&model.Task{}).Create(&t).Error; err != nil {
		fmt.Println(err)
		return errors.New("创建任务失败")
	}

	return nil
}

func DoAble(openid string, taskID int) (v model.Volunteer, t model.Task, e error) {
	v, err := getV(openid)
	if err != nil {
		return v, t, err
	}
	t, err = getTask(taskID)
	if err != nil {
		return v, t, err
	}

	if t.Publisher != v.Openid {
		return v, t, errors.New("您不是这个任务的发布者，无法修改此任务")
	}

	if time.Now().After(*t.Starttime) {
		return v, t, errors.New("此任务已经开始或结束，无法修改")
	}
	return v, t, nil
}

func (Task) Update(openid string, taskID int, info UpdateTaskInfo) error {
	_, t, err := DoAble(openid, taskID)
	if err != nil {
		return err
	}
	if info.Number < t.Already {
		return errors.New("当前报名的任务已经超过您传入的人数")
	}

	if err := model.DB.Table("tasks").Where("id = ?", taskID).Updates(&info).Error; err != nil {
		fmt.Println(err)
		return errors.New("更新失败")
	}
	return nil

}

func (Task) Delete(openid string, taskID int) error {
	_, t, err := DoAble(openid, taskID)
	if err != nil {
		return err
	}

	if err := model.DB.Delete(&t).Error; err != nil {
		fmt.Println(err)
		return errors.New("删除任务失败")
	}
	return nil
}
