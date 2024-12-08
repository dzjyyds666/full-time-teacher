package communication

import (
	"FullTimeTeacher/database"
	"FullTimeTeacher/log/logx"
	"FullTimeTeacher/models"
	"FullTimeTeacher/utils/result"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func CreateType(c *gin.Context) {
	var typeInfo models.ArticleType
	if err := c.ShouldBindJSON(&typeInfo); err != nil {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, err.Error(), nil))
		return
	}

	newUuid, _ := uuid.NewUUID()
	typeId := strings.ReplaceAll(newUuid.String(), "-", "")
	typeInfo.ArticleTypeID = typeId

	var typeInfotmp models.ArticleType
	res := database.MyDB.First(&typeInfotmp, "article_type_name = ? or article_type_id = ?", typeInfo.ArticleTypeName, typeInfo.ArticleTypeID)
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		logx.GetLogger("logx").Info("文章分类已经存在")
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ArticleTypeExists, "文章分类已经存在", nil))
		return
	}

	res = database.MyDB.Create(&typeInfo)
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("数据库写入异常:%s", res.Error.Error())
		panic("数据库写入异常" + res.Error.Error())
	}

	c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.RequestSuccess, "添加成功", nil))
}
