package communication

import (
	"FullTimeTeacher/database"
	"FullTimeTeacher/log/logx"
	"FullTimeTeacher/models"
	"FullTimeTeacher/utils/result"
	"FullTimeTeacher/utils/set"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type article struct {
	ArticleInfo *models.ArticleInfo
	UserInfo    *models.UserInfo
	ArticleType *models.ArticleType
}

type rootInfo struct {
	Username      string                `json:"username,omitempty"`
	Avatar        string                `json:"avatar,omitempty"`
	ArticleReplay *models.ArticleReplay `json:"article_reply,omitempty"`
	ChlidrenInfo  []*chlidrenInfo       `json:"chlidren_info,omitempty"`
}

type chlidrenInfo struct {
	Username      string                `json:"username,omitempty"`
	Avatar        string                `json:"avatar,omitempty"`
	ArticleReplay *models.ArticleReplay `json:"article_reply,omitempty"`
	ToUsername    string                `json:"to_username,omitempty"`
}

type userCache struct {
	UserID   string `json:"user_id,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Username string `json:"username,omitempty"`
}

// Publish 发布文章
func Publish(c *gin.Context) {
	var articleInfo *models.ArticleInfo

	// 绑定json数据
	err := c.ShouldBindJSON(articleInfo)
	if err != nil {
		logx.GetLogger("logx").Errorf("解析请求失败，%v", err)
		panic("解析请求失败，%v" + err.Error())
	}

	articleId, _ := uuid.NewUUID()
	articleInfo.ArticleID = strings.ReplaceAll(articleId.String(), "-", "")

	articleInfo.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	articleInfo.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	userId, _ := c.Get("user_id")
	articleInfo.UserID = userId.(string)

	// 插入数据库
	database.MyDB.Create(articleInfo)
}

// 修改文章
func UpdateArticle(c *gin.Context) {
	var articleInfo *models.ArticleInfo
	err := c.ShouldBindJSON(articleInfo)
	if err != nil {
		logx.GetLogger("logx").Errorf("解析请求失败，%v", err)
		panic("解析请求失败")
	}
	// 先查询到文章信息
	var articleformdatabase *models.ArticleInfo
	res := database.MyDB.First(articleformdatabase, "article_id = ?", articleInfo.ArticleID)
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("修改文章失败，%v", res.Error)
		panic("修改文章失败")
	}

	// 先判断文章是否是本人操作
	userid, _ := c.Get("user_id")
	if articleformdatabase.UserID != userid.(string) {
		logx.GetLogger("logx").Errorf("修改文章失败，%v", "没有权限操作该内容")
		panic("没有权限操作该内容")
	}

	res = database.MyDB.Model(articleformdatabase).
		Updates(models.ArticleInfo{
			ArticleName: articleInfo.ArticleName,
			ArticleDesc: articleInfo.ArticleDesc,
			UpdateTime:  time.Now().Format("2006-01-02 15:04:05"),
		})
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("修改文章失败，%v", res.Error)
		panic("修改文章失败")
	}
	c.JSON(http.StatusOK, result.Result{
		Code: result.EnmuHttptatus.RequestSuccess,
		Msg:  "修改文章成功",
	})
}

// 获取分页文章列表
func GetArticleList(c *gin.Context) {
	// 根据创建时间获取文章列表
	// 传入页码和每页数量
	page, pageSize := c.DefaultQuery("page", "1"), c.DefaultQuery("page_size", "10")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	// 查询数据库
	var articleInfos []*models.ArticleInfo
	res := database.MyDB.
		Where("is_deleted = ?", "0").
		Order("create_time DESC").
		Offset((pageInt - 1) * pageSizeInt).
		Limit(pageSizeInt).Find(&articleInfos)
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("获取文章列表失败，%v", res.Error)
		panic("获取文章列表失败")
	}
	var articleInfoList []*article
	for _, articleInfo := range articleInfos {
		// 查询用户信息
		var articleTmp *article
		articleTmp.ArticleInfo = articleInfo
		var userInfo *models.UserInfo
		res = database.MyDB.Select("username", "avatar").First(userInfo, "user_id = ?", articleInfo.UserID)
		if res.Error != nil {
			logx.GetLogger("logx").Errorf("获取文章列表失败，%v", res.Error)
			panic("获取文章列表失败")
		}
		articleTmp.UserInfo = userInfo
		var articleType *models.ArticleType
		res = database.MyDB.Select("article_type_name").First(articleType, "article_type_id = ?", articleInfo.TypeID)
		if res.Error != nil {
			logx.GetLogger("logx").Errorf("获取文章列表失败，%v", res.Error)
			panic("获取文章列表失败")
		}
		articleTmp.ArticleType = articleType

		articleInfoList = append(articleInfoList, articleTmp)
	}
	c.JSON(http.StatusOK, result.Result{
		Code: result.EnmuHttptatus.RequestSuccess,
		Msg:  "获取文章列表成功",
		Data: articleInfoList,
	})
}

// GetArticleInfo 获取文章信息
func GetArticleInfo(c *gin.Context) {
	// 获取文章信息
	articleId := c.Query("article_id")
	var articleInfo *models.ArticleInfo

	res := database.MyDB.First(&articleInfo, "article_id = ?", articleId)
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("获取文章信息失败，%v", res.Error)
		panic("获取文章信息失败，%v" + res.Error.Error())
	}

	var userinfo *models.UserInfo

	// 使用文章信息查询用户信息
	res = database.MyDB.Select("username", "avatar").First(userinfo, "user_id = ?", articleInfo.UserID)
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("获取用户信息失败，%v", res.Error)
		panic("获取用户信息失败，%v" + res.Error.Error())
	}

	// 使用文章信息查询文章类型
	var articleType *models.ArticleType
	res = database.MyDB.Select("article_type_name").First(articleType, "type_id = ?", articleInfo.TypeID)
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("获取文章类型失败，%v", res.Error)
		panic("获取文章类型失败，%v" + res.Error.Error())
	}
	// 组装数据
	c.JSON(http.StatusOK, result.Result{
		Code: result.EnmuHttptatus.RequestSuccess,
		Msg:  "获取文章信息成功",
		Data: &article{
			UserInfo:    userinfo,
			ArticleInfo: articleInfo,
			ArticleType: articleType,
		},
	})
}

// DeleteArticle 删除文章
func DeleteArticle(c *gin.Context) {
	articleId := c.Query("article_id")

	var article *models.ArticleInfo
	// 先判断文章的用户id与当前用户的id是否相等
	res := database.MyDB.First(&article, "article_id = ?", articleId)
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("获取文章信息失败:%s", res.Error.Error())
		panic("获取文章信息失败")
	}
	userid, _ := c.Get("user_id")
	if article.UserID != userid.(string) {
		c.JSON(http.StatusOK, result.Result{
			Code: result.EnmuHttptatus.RequestFail,
			Msg:  "没有权限操作该内容",
		})
	}

	// 删除文章
	database.MyDB.First(article, "article_id = ?", articleId)
	res = database.MyDB.Model(article).Update("is_deleted", "1")
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("删除文章失败:%s", res.Error.Error())
		panic("删除文章失败")
	}

	// 删除文章下面的回复
	database.MyDB.Model(&models.ArticleReplay{}).
		Where("article_id = ?", articleId).
		Update("is_deleted", "1")

	// 修改分类的文章数目
	database.MyDB.Model(&models.ArticleType{}).
		Where("article_type_id = ?", article.TypeID).
		UpdateColumn("article_num", gorm.Expr("article_num - ?", 1))
}

// 文章评论 ArticleReplay
func ArticleReplay(c *gin.Context) {
	var articleReplay *models.ArticleReplay
	err := c.ShouldBindJSON(articleReplay)
	if err != nil {
		logx.GetLogger("logx").Errorf("参数解析错误:%s", err.Error())
		panic("参数解析错误")
	}

	// 生成评论id
	replayId, _ := uuid.NewUUID()
	articleReplay.ArticleReplayID = strings.ReplaceAll(replayId.String(), "-", "")

	userid, _ := c.Get("user_id")
	articleReplay.UserID = userid.(string)
	articleReplay.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	articleReplay.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	res := database.MyDB.Create(articleReplay)
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("添加评论失败:%s", res.Error.Error())
		panic("添加评论失败")
	}

	c.JSON(http.StatusOK, result.Result{
		Code: result.EnmuHttptatus.RequestSuccess,
		Msg:  "添加评论成功",
	})
}

// 删除评论 DeleteReplay
func DeleteReplay(c *gin.Context) {
	replayId := c.Query("article_replay_id")

	var articleReplay *models.ArticleReplay
	res := database.MyDB.First(articleReplay, "article_replay_id = ?", replayId)
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("获取评论失败:%s", res.Error.Error())
		panic("获取评论失败")
	}

	// 判断评论的用户id与当前用户的id是否相等
	userid, _ := c.Get("user_id")
	if articleReplay.UserID != userid.(string) {
		c.JSON(http.StatusOK, result.Result{
			Code: result.EnmuHttptatus.RequestFail,
			Msg:  "没有权限操作该内容",
		})
	}

	res = database.MyDB.Model(articleReplay).Update("is_deleted", "1")
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("删除评论失败:%s", res.Error.Error())
		panic("删除评论失败")
	}

	c.JSON(http.StatusOK, result.Result{
		Code: result.EnmuHttptatus.RequestSuccess,
		Msg:  "删除评论成功",
	})
}

// GetCommentList 获取评论信息列表
func GetCommentList(c *gin.Context) {
	// 分页获取评论列表
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	articleId := c.Query("article_id")
	var rootsArticleReplayList []*rootInfo
	// 查询评论根节点 parent_id为空
	res := database.MyDB.Where("article_id = ? and is_deleted = ? and parent_id = IS NULL", articleId, "0").
		Order("create_time DESC").
		Offset((pageInt - 1) * pageSizeInt).
		Limit(pageSizeInt).
		Find(&rootsArticleReplayList)
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("获取评论列表失败:%s", res.Error.Error())
		panic("获取评论列表失败")
	}

	// 获取评论的用户id
	stringSet := set.NewStringSet()
	for _, replay := range rootsArticleReplayList {
		stringSet.Add(replay.ArticleReplay.UserID)
		var chlidrensInfo []*chlidrenInfo
		// 查询所有子评论
		err := database.MyDB.
			Where("article_id = ? and is_deleted = ? and parent_id = ?", articleId, "0", replay.ArticleReplay.UserID).
			Order("create_time").
			Find(chlidrensInfo).Error
		if err != nil {
			logx.GetLogger("logx").Errorf("获取子评论失败:%s", err.Error())
			panic("获取子评论失败")
		}
		for _, chlidren := range chlidrensInfo {
			stringSet.Add(chlidren.ArticleReplay.UserID)
		}
		replay.ChlidrenInfo = chlidrensInfo
	}

	// 一次性获取到所有的用户信息
	userInfo := GetCommentUserInfo(stringSet.List())
	for _, replay := range rootsArticleReplayList {
		replay.Avatar = userInfo[replay.ArticleReplay.UserID].Avatar
		replay.Username = userInfo[replay.ArticleReplay.UserID].Username
		for _, chlidren := range replay.ChlidrenInfo {
			chlidren.Avatar = userInfo[chlidren.ArticleReplay.UserID].Avatar
			chlidren.Username = userInfo[chlidren.ArticleReplay.UserID].Username
			chlidren.ToUsername = userInfo[chlidren.ArticleReplay.ToID].Username
		}
	}

	c.JSON(http.StatusOK, result.Result{
		Code: result.EnmuHttptatus.RequestSuccess,
		Msg:  "获取评论列表成功",
		Data: rootsArticleReplayList,
	})
}

func GetCommentUserInfo(userID []string) map[string]userCache {
	var userInfo []*models.UserInfo
	res := database.MyDB.Where("user_id in ?", userID).Select("user_id", "username", "avatar").First(userInfo)
	if res.Error != nil {
		logx.GetLogger("logx").Errorf("获取用户信息失败:%s", res.Error.Error())
		panic("获取用户信息失败")
	}

	userCacheTable := make(map[string]userCache)
	for _, user := range userInfo {
		//组装用户信息 ， 使用map存储，key为user_id，值为userCache
		userCacheTable[user.UserID] = userCache{
			Avatar:   user.Avatar,
			UserID:   user.UserID,
			Username: user.Username,
		}
	}
	return userCacheTable
}
