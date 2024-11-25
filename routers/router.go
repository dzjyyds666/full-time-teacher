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

			public.POST("/loginpass", login.LoginByPassword)
			public.POST("/loginver", login.LoginByVerification)

			public.POST("/ocr", ai.OCR)
			public.POST("/textanswer", ai.TextAnswer)

		}

		// 需要认证的路由组
		auth := v1.Group("")
		auth.Use(middlewares.TokenVerify())
		{
			auth.POST("/putfile", cos.PutFile)
		}
	}
}
