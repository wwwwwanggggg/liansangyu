package model

type Organization struct {
	Openid string `gorm:"not null;primaryKey;type:VARCHAR(255)" json:"openid"`

	User User `gorm:"foreignKey:Openid;reference:Openid"`

	TimeModel
	Name string `gorm:"not null" json:"name"`
	Logo string

	Admin     []User      `gorm:"many2many:organization_admins"`
	Volunteer []Volunteer `gorm:"many2many:organization_volunteers"`
	Elder     []Elder     `gorm:"many2many:organization_elders"`
}
