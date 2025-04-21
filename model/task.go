package model

import "time"

type Task struct {
	BaseModel

	Title string `gorm:"column:title;not null;type:VARCHAR(40)" json:"title"`

	Starttime *time.Time `gorm:"column:start_time;not null;" json:"start_time"`
	Endtime   *time.Time `gorm:"column:end_time;not null" json:"end_time"`

	Participants []Volunteer `gorm:"many2many:task_participants"`

	Longitude float64 `gorm:"column:longitude;not null" json:"longitude"`
	Latitude  float64 `gorm:"column:latitude;not null" json:"latitude"`

	Desc string `gorm:"column:desc;not null" json:"desc"`

	Publisher string `gorm:"column:publisher;not null" json:"publisher"`

	Number  uint16 `gorm:"column:number;not null" json:"number"`
	Already uint16 `gorm:"column:already;not null;default:0" json:"already"`
}
