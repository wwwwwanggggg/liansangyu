package model

type Elder struct {
	Openid string `gorm:"column:openid;primaryKey;not null;type:VARCHAR(191)" json:"openid"`

	User User `gorm:"foreignKey:Openid;reference:Openid"`
	TimeModel
	Disease string `gorm:"column:disease" json:"disease"`

	Longitude float64 `gorm:"column:longitude;comment:经度;not null" json:"longitude"`
	Latitude  float64 `gorm:"column:latitude;comment:纬度;not null" json:"latitude"`
}
