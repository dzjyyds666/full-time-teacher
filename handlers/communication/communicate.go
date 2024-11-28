package communication

import (
	"FullTimeTeacher/database"
	"FullTimeTeacher/models"
	"FullTimeTeacher/utils/result"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Publish 发布问题
func Publish(c *gin.Context) {
	var problemInfo models.ProblemInfo

	err := c.ShouldBindJSON(&problemInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Result{
			Code: result.EnmuHttptatus.ParamError,
			Msg:  err.Error(),
		})
		return
	}

	// 问题id
	problemInfo.ProblemID = uuid.New().String()
	problemInfo.ProblemID = strings.ReplaceAll(problemInfo.ProblemID, "-", "")

	// 用户id
	userID, _ := c.Get("user_id")
	problemInfo.UserID = userID.(string)

	// 创建时间
	problemInfo.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	problemInfo.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	// 插入数据库
	err = database.MyDB.Create(&problemInfo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Result{
			Code: result.EnmuHttptatus.SystemError,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result.Result{
		Code: result.EnmuHttptatus.RequestSuccess,
		Msg:  "发布成功",
	})
}
