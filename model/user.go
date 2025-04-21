package model

type User struct {
	BaseModel

	Openid string `gorm:"column:openid;comment:wx openid;type:VARCHAR(255);not null" json:"openid"`
	Name   string `gorm:"column:name;type:VARCHAR(40);not null" json:"name"`
	Phone  string `gorm:"column:phone;type:VARCHAR(11);not null" json:"phone"`

	UserType uint8 `gorm:"column:user_type;not null" json:"user_type"` // 1 volunteer 2 elder 3 monitor 10 normal admin  100 super admin
}
