package model

type OrganizationAdmins struct {
	Organization Organization `gorm:"foreignKey:OrganizationOpenid;reference:Openid"`
	User         User         `gorm:"foreignKey:AdminOpenid;reference:Openid"`

	OrganizationOpenid string `gorm:"not null;type:VARCHAR(191)"`
	AdminOpenid        string `gorm:"not null;type:VARCHAR(191)"`
	Passed             bool   `gorm:"not null;default:false"`
}

type OrganizationVolunteers struct {
	Organization Organization `gorm:"foreignKey:OrganizationOpenid;reference:Openid"`
	Volunteer    Volunteer    `gorm:"foreignKey:VolunteerOpenid;reference:Openid"`

	OrganizationOpenid string `gorm:"not null;type:VARCHAR(191)"`
	VolunteerOpenid    string `gorm:"not null;type:VARCHAR(191)"`
	Passed             bool   `gorm:"not null;default:false"`
}

type OrganizationElders struct {
	Organization Organization `gorm:"foreignKey:OrganizationOpenid;reference:Openid"`
	Elder        Elder        `gorm:"foreignKey:ElderOpenid;reference:Openid"`

	OrganizationOpenid string `gorm:"not null;type:VARCHAR(191)"`
	ElderOpenid        string `gorm:"not null;type:VARCHAR(191)"`
	Passed             bool   `gorm:"not null;default:false"`
}

type TaskParticipants struct {
	Task            Task      `gorm:"foreignKey:Taskid;reference:ID"`
	Volunteer       Volunteer `gorm:"foreignKey:VolunteerOpenid;reference:Openid"`
	Taskid          int64     `gorm:"column:task_id"`
	VolunteerOpenid string    `gorm:"not null;type:VARCHAR(191)"`
}
