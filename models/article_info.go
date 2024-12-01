package models

type ArticleInfo struct {
	ArticleID   string `json:"article_id,omitempty" gorm:"primaryKey;size:36"` // 问题ID，设置为主键，长度为36
	ArticleName string `json:"article_name,omitempty" gorm:"size:100;unique"`  // 问题名称，设置为唯一，长度为100
	ArticleDesc string `json:"article_desc,omitempty" gorm:"size:500"`         // 问题描述，长度为500
	UserID      string `json:"user_id,omitempty" gorm:"size:36"`               // 用户ID，长度为36
	CreateTime  string `json:"create_time,omitempty" gorm:"size:20"`           // 创建时间，长度为20
	UpdateTime  string `json:"update_time,omitempty" gorm:"size:20"`           // 更新时间，长度为20
	Status      string `json:"status,omitempty" gorm:"size:1;default:1"`       // 状态，长度为1，默认值为1 0：下架，1：上架
	IsDeleted   string `json:"is_deleted,omitempty" gorm:"size:1;default:0"`   // 是否删除，长度为1，默认值为0 0：未删除，1：已删除
	TypeID      string `json:"type_id,omitempty" gorm:"size:36"`               // 分类id
}

func (pi *ArticleInfo) TableName() string {
	return "article_info"
}
