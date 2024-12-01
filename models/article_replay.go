package models

type ArticleReplay struct {
	ArticleReplayID string `json:"article_replay_id,omitempty" gorm:"primaryKey;size:36"` // 问题回复ID，设置为主键，长度为36
	ArticleID       string `json:"article_id,omitempty" gorm:"size:36"`                   // 问题ID，长度为36
	UserID          string `json:"user_id,omitempty" gorm:"size:36"`                      // 用户ID，长度为36
	ParentID        string `json:"parent_id,omitempty" gorm:"size:36"`                    // 父回复id
	ToID            string `json:"to_id,omitempty" gorm:"size:36"`                        // 回复对象id
	ReplayContent   string `json:"replay_content,omitempty" gorm:"size:500"`              // 回复内容，长度为500
	CreateTime      string `json:"create_time,omitempty" gorm:"size:20"`                  // 创建时间，长度为20
	UpdateTime      string `json:"update_time,omitempty" gorm:"size:20"`                  // 更新时间，长度为20
	Status          string `json:"status,omitempty" gorm:"size:1;default:1"`              // 状态，长度为1，默认值为1 0：下架，1：上架
	IsDeleted       string `json:"is_deleted,omitempty" gorm:"size:1;default:0"`          // 是否删除，长度为1，默认值为0 0：未删除，1：已删除
}

func (pr *ArticleReplay) TableName() string {
	return "article_replay"
}
