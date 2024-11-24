package main

import (
	"FullTimeTeacher/config"
	"FullTimeTeacher/database"
	"FullTimeTeacher/log/logx"
	"FullTimeTeacher/routers"
	"FullTimeTeacher/utils/result"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//全局异常处理
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logx.GetLogger("logx").Errorf("系统异常:%v", recovered)
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.SystemError, "系统异常,请稍后再试", recovered))
	}))
	// 初始化database
	database.InitDatabase(*config.GlobalConfig)
	// 注册路由
	routers.InitRouter(r, config.GlobalConfig)

	r.Run(":9999")
}
