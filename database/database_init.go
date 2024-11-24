package database

import "FullTimeTeacher/config"

func InitDatabase(config config.Config) {
	InitMySQL(config.MySQL)
	InitRedis(config.Redis)
}
