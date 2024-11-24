package ai

import (
	"FullTimeTeacher/config"
	"FullTimeTeacher/log/logx"
	"FullTimeTeacher/utils/result"
	"net/http"
	"os"

	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/gin-gonic/gin"
)

// 文本回答
func TextAnswer(c *gin.Context) {
	message := c.PostForm("message")
	if message == "" {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "请输入问题", nil))
		return
	}

	logx.GetLogger("logx").Infof("message:%s", message)

	// 设置环境变量
	os.Setenv("QIANFAN_ACCESS_KEY", config.GlobalConfig.AI.ApiKey)
	os.Setenv("QIANFAN_SECRET_KEY", config.GlobalConfig.AI.SecretKey)

	chat := qianfan.NewChatCompletion()

	resp, _ := chat.Stream(
		c,
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(message),
			},
		},
	)
	logx.GetLogger("logx").Infof("resp:%v", resp)

	// for {
	// 	// r, err := resp.Recv()
	// 	// if err != nil {
	// 	// 	panic("ai错误:" + err.Error())
	// 	// }
	// 	// if resp.IsEnd { // 判断是否结束
	// 	// 	break
	// 	// }
	// 	// // 把文字流式相应给前端
	// 	// c.Stream(func(w io.Writer) bool {
	// 	// 	bw := bufio.NewWriter(w)

	// 	// })
	// }
}
