package model

type Elder struct {
	BaseModel

	Openid  string `gorm:"column:openid;type:VARCHAR(255);not null" json:"openid"`
	Disease string `gorm:"column:disease" json:"disease"`

	Longitude float64 `gorm:"column:longitude;comment:经度;not null" json:"longitude"`
	Latitude  float64 `gorm:"column:latitude;comment:纬度;not null" json:"latitude"`
}
