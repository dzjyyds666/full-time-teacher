package jwt

import (
	"FullTimeTeacher/config"
	"FullTimeTeacher/log/logx"
	"testing"
)

func TestParseJWT(t *testing.T) {

	token := NewJWTUtils(config.JwtConfig{
		SecretKey:      "aaronlikeprogramming",
		ExpirationTime: 24,
	}).CreateJWT("5bdcc349e21e492ba5926dd0c5bac86c")

	//token := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNWJkY2MzNDllMjFlNDkyYmE1OTI2ZGQwYzViYWM4NmMiLCJleHAiOjg2NDAwMDAwfQ.IDw7TOrTzwm6zxlyZpLaINm1PZfcs9bDmIeIaoItjVw"

	logx.GetLogger("logx").Infof("TestParseJWT|token: %v", token)

	claims, err := NewJWTUtils(config.JwtConfig{
		SecretKey:      "aaronlikeprogramming",
		ExpirationTime: 24,
	}).ParseJWT(token)
	if err != nil {
		logx.GetLogger("logx").Errorf("TestParseJWT|err: %v", err)
	}
	logx.GetLogger("logx").Infof("TestParseJWT|claims: %v", claims.UserId)
}
