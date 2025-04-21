package model

type Elder struct {
	BaseModel

	Openid  string `gorm:"column:openid;type:VARCHAR(255);not null" json:"openid"`
	Disease string `gorm:"column:disease" json:"disease"`

	Longtitude float64 `gorm:"column:longtitude;comment:经度;not null" json:"longtitude"`
	Latitude   float64 `gorm:"column:latitude;comment:纬度;not null" json:"latitude"`
}
