package login

import (
	"FullTimeTeacher/database"
	"FullTimeTeacher/log/logx"
	"FullTimeTeacher/models"
	"FullTimeTeacher/utils/email"
	"FullTimeTeacher/utils/result"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

type RegisterInfo struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captcha_id"`
	Code      string `json:"code"`
}

// Register 注册
func Register(c *gin.Context) {

	var registerInfo RegisterInfo
	if err := c.ShouldBindJSON(&registerInfo); err != nil {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, err.Error(), nil))
		return
	}
	logx.GetLogger("logx").Infof("Register|RegisterInfo: %v", registerInfo)

	if len(registerInfo.Email) == 0 || len(registerInfo.Password) == 0 {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "邮箱或密码不能为空", nil))
		return
	}

	if len(registerInfo.Captcha) == 0 || len(registerInfo.CaptchaId) == 0 {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "图片验证码不能为空", nil))
		return
	}

	if len(registerInfo.Code) == 0 {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "邮箱验证码不能为空", nil))
		return
	}

	// 验证图片验证码
	captcha, err := database.RDB.Get(c, fmt.Sprintf(database.Redis_Captcha_Key, registerInfo.CaptchaId)).Result()
	if err != nil {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "图片验证码已过期", nil))
		return
	}
	if captcha != registerInfo.Captcha {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "图片验证码错误", nil))
		return
	}
	// 验证邮箱验证码
	verificationCode, err := database.RDB.Get(c, fmt.Sprintf(database.Redis_Verification_Code_Key, registerInfo.Email)).Result()
	if err != nil {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "邮箱验证码错误", nil))
		return
	}
	if verificationCode != registerInfo.Code {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "邮箱验证码错误", nil))
		return
	}

	// uuid生成用户id
	userID := uuid.New().String()
	// 去除uuid的-
	userID = strings.Replace(userID, "-", "", -1)

	//密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(registerInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.ParamError, "密码加密失败", nil))
		return
	}

	// 设置默认用户名
	username := fmt.Sprintf("user_%s", userID[:8])

	// 注册用户
	database.MyDB.Create(&models.UserInfo{
		UserID:     userID,
		Username:   username,
		Email:      registerInfo.Email,
		Password:   string(hash),
		Avatar:     "https://img.tuxiangyan.com/uploads/allimg/2021082810/rd22b0qzue1.jpg",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	})

	c.JSON(http.StatusOK, result.NewResult(result.EnmuHttptatus.RequestSuccess, "注册成功", nil))
}

// 获取图形验证码
func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 200, 5, 0.8, 75)
	store := base64Captcha.DefaultMemStore

	//生成图形验证码
	captcha := base64Captcha.NewCaptcha(driver, store)

	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		panic("获取图片验证码失败" + err.Error())
	}

	// 使用redis存取验证码
	err = database.RDB.Set(context.Background(), fmt.Sprintf(database.Redis_Captcha_Key, id), answer, time.Minute*5).Err()

	if err != nil {
		panic("redis存取验证码失败" + err.Error())
	}

	c.JSON(http.StatusOK, result.Result{
		Code: result.EnmuHttptatus.RequestSuccess,
		Msg:  "获取验证码成功",
		Data: map[string]interface{}{
			"id":    id,
			"image": b64s,
		},
	})
}

// 发送注册验证码
func SendEmail(c *gin.Context) {

	to := c.Query("email")

	// 生成随机验证码
	verificationCode := GenerateVerificationCode(6)

	// 把验证码存入redis
	ok, err := database.RDB.SetNX(c, fmt.Sprintf(database.Redis_Verification_Code_Key, to), verificationCode, time.Minute*5).Result()
	if err != nil {
		panic("redis存取验证码失败" + err.Error())
	}

	if !ok {
		c.JSON(http.StatusOK, result.Result{
			Code: result.EnmuHttptatus.RequestSuccess,
			Msg:  "验证码已发送，请稍后再试",
		})
		return
	}

	subject := "验证码"
	body := fmt.Sprintf(
		`<p>您的验证码是: %s</p>
		<p>请在5分钟内完成注册</p>
		<p>请不要回复此邮件</p>`, verificationCode)

	err = email.SendEmail(to, subject, body)
	if err != nil {
		panic("发送邮件失败:" + err.Error())
	}

	c.JSON(http.StatusOK, result.Result{
		Code: result.EnmuHttptatus.RequestSuccess,
		Msg:  "发送成功，请及时查收",
	})
}

func GenerateVerificationCode(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	rand.Seed(uint64(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
