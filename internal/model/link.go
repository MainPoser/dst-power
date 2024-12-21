package model

type Link struct {
	Model
	Name   string `gorm:"column:name;NOT NULL" json:"name"`               // 友链名称
	Url    string `gorm:"column:url;NOT NULL" json:"url"`                 // 友链地址
	Image  string `gorm:"column:image" json:"image"`                      // logo
	Status int    `gorm:"column:status;default:1;NOT NULL" json:"status"` // 状态 1-启用 2-禁用
	Sort   int    `gorm:"column:sort;default:0;NOT NULL" json:"sort"`     // 排序
}

func (l *Link) TableName() string {
	return "link"
}
