package database

import (
	"FullTimeTeacher/config"
	"FullTimeTeacher/log/logx"
	"FullTimeTeacher/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局mysql客户端
var MyDB *gorm.DB

func InitMySQL(config config.MySQLConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName)

	var err error
	MyDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logx.GetLogger("logx").Errorf("数据库连接失败: %v", err)
	}

	logx.GetLogger("logx").Info("MySQL连接成功")

	MyDB.AutoMigrate(&models.ArticleInfo{}, &models.ArticleReplay{}, &models.ArticleType{}, &models.UserInfo{})
}
