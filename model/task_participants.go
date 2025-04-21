package model

type TaskParticipants struct {
	TaskID      uint64 `gorm:"column:task_id;not null" json:"task_id"`
	VolunteerID uint64 `gorm:"column:volunteer_id;not null" json:"volunteer_id"`
}
