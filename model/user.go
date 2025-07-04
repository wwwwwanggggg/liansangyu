package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Openid string `gorm:"comment:wx openid;not null;primaryKey;type:VARCHAR(255)" json:"openid"`
	TimeModel

	Name  string `gorm:"column:name;type:VARCHAR(40);not null" json:"name"`
	Phone string `gorm:"column:phone;type:VARCHAR(11);not null;unique" json:"phone"`

	UserType uint8 `gorm:"column:user_type;not null" json:"user_type"` // 1 volunteer 2 elder 3 monitor 10 normal admin  100 super admin

	// Elder          Elder        `gorm:"foreignKey:Openid;reference:Openid" json:"elder"`
	// Monitor        Monitor      `gorm:"foreignKey:Openid;reference:Openid" json:"monitor"`
	// Organization   Organization `gorm:"foreignKey:Openid;reference:Openid" json:"organization"`
	IsVolunteer    bool `gorm:"not null;default:false"`
	IsElder        bool `gorm:"not null;default:false"`
	IsMonitor      bool `gorm:"not null;default:false"`
	IsOrganization bool `gorm:"not null;default:false"`
}

func (u *User) AfterUpdate(tx *gorm.DB) error {
	err := tx.Model(&u).Update("updated_at", time.Now()).Error
	return err
}
