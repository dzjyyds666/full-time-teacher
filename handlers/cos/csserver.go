package cos

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 直接上传文件
func PutFile(c *gin.Context) {

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": "上传成功",
	})
}
