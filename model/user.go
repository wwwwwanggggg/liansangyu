package model

type User struct {
	Openid string `gorm:"comment:wx openid;not null;primaryKey;type:VARCHAR(255)" json:"openid"`
	Name   string `gorm:"column:name;type:VARCHAR(40);not null" json:"name"`
	Phone  string `gorm:"column:phone;type:VARCHAR(11);not null" json:"phone"`

	UserType uint8 `gorm:"column:user_type;not null" json:"user_type"` // 1 volunteer 2 elder 3 monitor 10 normal admin  100 super admin

	Volunteer    Volunteer    `gorm:"foreignKey:Openid;reference:Openid" json:"volunteer"`
	Elder        Elder        `gorm:"foreignKey:Openid;reference:Openid" json:"elder"`
	Monitor      Monitor      `gorm:"foreignKey:Openid;reference:Openid" json:"monitor"`
	Organization Organization `gorm:"foreignKey:Openid;reference:Openid" json:"organization"`
}
