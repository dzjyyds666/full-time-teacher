package ai

import (
	"FullTimeTeacher/config"
	"FullTimeTeacher/log/logx"
	"FullTimeTeacher/utils/ocr"
	"FullTimeTeacher/utils/result"
	"bytes"
	"io"
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
	for {
		r, err := resp.Recv()
		if err != nil {
			panic(err)
		}
		if resp.IsEnd { // 判断是否结束
			break
		}
		// 把文字流式相应给前端
		c.SSEvent("message", r.Result)
		c.Writer.Flush()
	}
}

func OCR(c *gin.Context) {
	// 获取图片字节数组
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "请上传图片", nil))
		return
	}

	// 读取文件内容
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "请上传图片", nil))
		return
	}
	defer src.Close()

	var imageBytes bytes.Buffer
	io.Copy(&imageBytes, src)

	ocr := ocr.NewOCR()
	text, err := ocr.RecongizeFromBytes(imageBytes.Bytes())
	if err != nil {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "识别失败", nil))
		return
	}

	logx.GetLogger("logx").Infof("text:%s", text)

	//执行ai回答
	// 设置环境变量
	os.Setenv("QIANFAN_ACCESS_KEY", config.GlobalConfig.AI.ApiKey)
	os.Setenv("QIANFAN_SECRET_KEY", config.GlobalConfig.AI.SecretKey)

	chat := qianfan.NewChatCompletion()

	resp, _ := chat.Stream(
		c,
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(text),
			},
		},
	)
	logx.GetLogger("logx").Infof("resp:%v", resp)
	for {
		response, err := resp.Recv()
		if resp.IsEnd {
			break
		}
		if err != nil {
			logx.GetLogger("logx").Errorf("接收消息错误:%v", err)
			break
		}
		// 发送数据到客户端
		c.SSEvent("message", response.Result)
		c.Writer.Flush()
	}
}
