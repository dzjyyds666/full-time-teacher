package cos

import (
	"context"
	"net/http"

	globalConfig "FullTimeTeacher/config"
	"FullTimeTeacher/utils/result"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/gin-gonic/gin"
)

// 生成s3的sts临时凭证
func GetStsToken(c *gin.Context) {

	//加载配置
	creds := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(globalConfig.GlobalConfig.Cos.AccessId, globalConfig.GlobalConfig.Cos.SecretKey, ""),
	)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("cn-north"),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		panic("load config error:" + err.Error())
	}

	// 创建sts客户端
	stsClient := sts.NewFromConfig(cfg)
	sessionName := "cos-sts"

	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(""),
		RoleSessionName: aws.String(sessionName),
	}

	res, err := stsClient.AssumeRole(context.TODO(), input)
	if err != nil {
		panic("assume role error:" + err.Error())
	}

	c.JSON(http.StatusOK, result.Result{
		Code: result.EnmuHttptatus.RequestSuccess,
		Msg:  "success",
		Data: res,
	})

}
