package routers

import (
	"FullTimeTeacher/config"
	"FullTimeTeacher/handlers/ai"
	"FullTimeTeacher/handlers/cos"
	"FullTimeTeacher/handlers/login"
	"FullTimeTeacher/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, config *config.Config) {

	v1 := r.Group("/api/v1")
	{
		//无需认证的路由组
		public := v1.Group("")
		{
			public.GET("/captcha", login.GetCaptcha)
			public.POST("/register", login.Register)
			public.GET("/verification", login.SendEmail)

			public.POST("/loginpass", login.LoginByPassword)    // 密码登录
			public.POST("/loginver", login.LoginByVerification) // 验证码登录

		}

		// 需要认证的路由组
		auth := v1.Group("")
		auth.Use(middlewares.TokenVerify())
		{
			auth.GET("/getStsToken", cos.GetStsToken) // 获取sts临时凭证
			auth.POST("/ocr", ai.OCR)                 // 图片识别
			auth.POST("/textanswer", ai.TextAnswer)   // 文本回答
			auth.POST("/loginout", login.LogOut)      // 退出登录
		}
	}
}
