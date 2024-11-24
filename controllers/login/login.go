package login

import (
	"FullTimeTeacher/config"
	"FullTimeTeacher/database"
	"FullTimeTeacher/log/logx"
	"FullTimeTeacher/models"
	"FullTimeTeacher/utils/jwt"
	"FullTimeTeacher/utils/result"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type LoginInfo struct {
	Email        string `json:"email,omitempty"`
	Password     string `json:"password,omitempty"`
	Verification string `json:"verification,omitempty"`
}

// LoginByPassword 密码登录
func LoginByPassword(c *gin.Context) {

	var loginInfo LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		// 参数错误
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, err.Error(), nil))
		return
	}

	logx.GetLogger("logx").Infof("LoginByPassword|loginInfo: %v", loginInfo)
	var user models.UserInfo
	Result := database.MyDB.Where("email = ?", loginInfo.Email).First(&user)
	if Result.Error != nil {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, Result.Error.Error(), nil))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password)); err != nil {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "密码错误", nil))
		return
	}

	//生成token
	token := jwt.NewJWTUtils(config.GlobalConfig.JWT).CreateJWT(user.UserID)

	//把token存入redis
	err := database.RDB.Set(c, fmt.Sprintf(database.Redis_Token_Key, user.UserID), token, time.Hour*time.Duration(config.GlobalConfig.JWT.ExpirationTime)).Err()
	if err != nil {
		panic("redis错误:" + err.Error())
	}

	// 设置响应头
	c.Header("Authorization", token)

	c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.RequestSuccess, "登录成功", models.UserInfo{
		UserID:        user.UserID,
		Email:         user.Email,
		Username:      user.Username,
		Avatar:        user.Avatar,
		Sex:           user.Sex,
		Bio:           user.Bio,
		LastLoginIp:   user.LastLoginIp,
		LastLoginTime: user.LastLoginTime,
		Experience:    user.Experience,
	}))
}

func LoginByVerification(c *gin.Context) {

	var loginInfo LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		panic("参数错误:" + err.Error())
	}

	// 验证验证码
	verificationCode, err := database.RDB.Get(c, loginInfo.Email).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "验证码已过期", nil))
			return
		} else {
			panic("redis错误:" + err.Error())
		}
	}

	if verificationCode != loginInfo.Verification {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "验证码错误", nil))
		return
	}

	// 查询用户
	var user models.UserInfo
	Result := database.MyDB.Where("email = ?", loginInfo.Email).First(&user)
	if Result.Error != nil {
		panic("数据库错误:" + Result.Error.Error())
	}

	//生成token
	token := jwt.NewJWTUtils(config.GlobalConfig.JWT).CreateJWT(user.UserID)

	//把token存入redis
	err = database.RDB.Set(c, fmt.Sprintf(database.Redis_Token_Key, user.UserID), token, time.Hour*time.Duration(config.GlobalConfig.JWT.ExpirationTime)).Err()
	if err != nil {
		panic("redis错误:" + err.Error())
	}

	// 设置响应头
	c.Header("Authorization", token)
	c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.RequestSuccess, "登录成功", models.UserInfo{
		UserID:        user.UserID,
		Email:         user.Email,
		Username:      user.Username,
		Avatar:        user.Avatar,
		Sex:           user.Sex,
		Bio:           user.Bio,
		LastLoginIp:   user.LastLoginIp,
		LastLoginTime: user.LastLoginTime,
		Experience:    user.Experience,
	}))
}
