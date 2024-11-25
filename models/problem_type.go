package models

type ProblemType struct {
	ProblemTypeID   string `json:"problem_type_id,omitempty" gorm:"primaryKey;size:36"` // 问题类型ID，设置为主键，长度为36
	ProblemTypeName string `json:"problem_type_name,omitempty" gorm:"size:100;unique"`  // 问题类型名称，设置为唯一，长度为100
	ProblemTypeDesc string `json:"problem_type_desc,omitempty" gorm:"size:500"`         // 问题类型描述，长度为500
	IsDeleted       string `json:"is_deleted,omitempty" gorm:"size:1;default:0"`        // 是否删除，长度为1，默认值为0 0：未删除，1：已删除
	ProblemNum      string `json:"problem_num,omitempty" gorm:"size:5;default:0"`       // 问题数量，长度为10，默认值为0

	ProblemInfos []ProblemInfo `json:"problem_infos,omitempty" gorm:"many2many:problem_type_relation;"` // 问题信息，多对多关系
}

func (pt *ProblemType) TableName() string {
	return "problem_type"
}
