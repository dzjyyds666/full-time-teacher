package models

type ArticleType struct {
	ArticleTypeID   string `json:"article_type_id,omitempty" gorm:"primaryKey;size:36"` // 问题类型ID，设置为主键，长度为36
	ArticleTypeName string `json:"article_type_name,omitempty" gorm:"size:100;unique"`  // 问题类型名称，设置为唯一，长度为100
	ArticleTypeDesc string `json:"article_type_desc,omitempty" gorm:"size:500"`         // 问题类型描述，长度为500
	IsDeleted       string `json:"is_deleted,omitempty" gorm:"size:1;default:0"`        // 是否删除，长度为1，默认值为0 0：未删除，1：已删除
	ArticleNum      string `json:"article_num,omitempty" gorm:"size:5;default:0"`       // 问题数量，长度为10，默认值为0
	UserNumber      string `json:"user_number,omitempty" gorm:"size:9"`                 // 分类订阅人数

	// 分类和用户人数多对多
	Users []*UserInfo `json:"users,omitempty" gorm:"many2many:user_article_type;"`
}

func (pt *ArticleType) TableName() string {
	return "article_type"
}
