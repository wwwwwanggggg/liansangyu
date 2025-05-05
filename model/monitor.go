package model

type Monitor struct {
	Openid string `gorm:"column:openid;not null;primaryKey" json:"openid"`

	User User `gorm:"foreignKey:Openid;reference:Openid"`

	TimeModel
	Elder       Elder  `gorm:"foreignKey:ElderOpenid;reference:Openid"`
	ElderOpenid string `gorm:"not null;type:VARCHAR(191)"`
	Passed      bool   `gorm:"not null;default:false"`
}
