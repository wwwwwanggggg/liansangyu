package model

type Notification struct {
	BaseModel

	User       User   `gorm:"foreignKey:UserOpenid;reference:Openid" json:"user"`
	UserOpenid string `gorm:"not null;type:VARCHAR(191)" json:"user_openid"`

	Title   string `json:"title"`
	Content string `json:"content"`
}
