package model

type Monitor struct {
	BaseModel

	Openid string `gorm:"column:openid;type:VARCHAR(255);not null" json:"openid"`

	MonitorPhone string `gorm:"column:monitor_phone;not null;type:VARCHAR(11)" json:"monitor_phone"`
	Passed       bool   `gorm:"column:passed;default:false" json:"passed"`
}
