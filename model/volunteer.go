package model

import "time"

type Volunteer struct {
	Openid string `gorm:"not null;primaryKey;type:VARCHAR(191)" json:"openid"`

	School string `gorm:"column:school;type:VARCHAR(40);not null" json:"school"`
	Clazz  string `gorm:"column:clazz;type:VARCHAR(10);not null;comment:班级" json:"clazz"`
	Skills string `gorm:"column:skills;type:json" json:"skills"`

	Hours     uint16     `gorm:"column:hours;not null;default:0" json:"hours"`
	Starttime *time.Time `gorm:"column:start_time"`

	Tasks []Task `gorm:"many2many:task_participants;"`
}
