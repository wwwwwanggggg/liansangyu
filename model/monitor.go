package model

type Monitor struct {
	Openid string `gorm:"column:openid;not null;primaryKey" json:"openid"`

	Elder       Elder  `gorm:"foreignKey:ElderOpenid;reference:Openid"`
	ElderOpenid string `gorm:"not null;type:VARCHAR(191)"`
}
