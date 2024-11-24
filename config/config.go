package config

import (
	"FullTimeTeacher/log/logx"

	"github.com/spf13/viper"
)

// Config 结构体表示应用程序的配置
type Config struct {
	AppName    string `mapstructure:"app_name"`
	AppVersion string `mapstructure:"app_version"`
	ServerPort int    `mapstructure:"server_port"`
	MySQL      MySQLConfig
	JWT        JwtConfig
	Redis      RedisConfig
	Email      EmailConfig
	AI         AIConfig
}

// MySQLConfig 结构体表示 MySQL 的配置
type MySQLConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DBName   string `mapstructure:"db_name"`
}

// JwtConfig 结构体表示 JWT 的配置
type JwtConfig struct {
	SecretKey      string `mapstructure:"secret_key"`
	ExpirationTime int64  `mapstructure:"expiration_time"`
}

// redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type EmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
	Alias    string `mapstructure:"alias"`
}

type AIConfig struct {
	ApiKey    string `mapstructure:"api_key"`
	SecretKey string `mapstructure:"secret_key"`
	Endpoint  string `mapstructure:"endpoint"`
}

var GlobalConfig *Config

func init() {
	logx.GetLogger("logx").Info("初始化配置")
	// 读取配置文件
	cfg, err := LoggingConfig()
	if err != nil {
		logx.GetLogger("logx").Errorf("读取配置文件失败: %v", err)
		return
	}
	// 将配置文件绑定结构体
	GlobalConfig = cfg
}

func LoggingConfig() (*Config, error) {

	configFilePath := "./config/config.toml"

	// 使用viper读取配置文件
	viper.SetConfigType("toml")
	viper.SetConfigFile(configFilePath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// 将配置文件绑定结构体
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
