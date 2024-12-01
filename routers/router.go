package routers

import (
	"FullTimeTeacher/config"
	"FullTimeTeacher/handlers/ai"
	"FullTimeTeacher/handlers/communication"
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

			public.POST("/loginpass", login.LoginByPassword)            // 密码登录
			public.POST("/loginver", login.LoginByVerification)         // 验证码登录
			public.GET("/getArticle", communication.GetArticleInfo)     // 获取文章信息
			public.GET("/getArticleList", communication.GetArticleList) // 获取文章列表

		}

		// 需要认证的路由组
		auth := v1.Group("")
		auth.Use(middlewares.TokenVerify())
		{
			auth.GET("/getStsToken", cos.GetStsToken) // 获取sts临时凭证
			auth.POST("/ocr", ai.OCR)                 // 图片识别
			auth.POST("/textanswer", ai.TextAnswer)   // 文本回答

			auth.POST("/loginout", login.LogOut) // 退出登录

			auth.POST("/publish", communication.Publish)             // 发布文章
			auth.POST("/deleteArticle", communication.DeleteArticle) // 删除文章
			auth.POST("/updateArticle", communication.UpdateArticle) // 更新文章
			auth.POST("/addComment", communication.ArticleReplay)    // 添加评论
			auth.POST("/deleteComment", communication.DeleteReplay)  // 删除评论
		}
	}
}
