package sdk

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsClient struct {
	AccessKeyID     string     `json:"access_key_id,omitempty"`     // 访问密钥ID
	SecretAccessKey string     `json:"secret_access_key,omitempty"` // 密钥
	Region          string     `json:"region,omitempty"`            // 区域
	Buckets         string     `json:"buckets,omitempty"`           // 桶
	Endpoint        string     `json:"endpoint,omitempty"`          // 端点
	S3Client        *s3.Client // s3客户端
}

func NewAwsClient(config *AwsClient) *AwsClient {
	// 自定义http客户端
	hcli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// 解析地址
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}
				// 尝试解析dns
				ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
				if err != nil || len(ips) == 0 {
					// 如果解析失败，则使用localhost
					addr = net.JoinHostPort("127.0.0.1", port)
				} else {
					//使用解析出来的地址
					addr = net.JoinHostPort(ips[0].String(), port)
				}
				// 连接
				return net.Dial(network, addr)
			},
		},
		// 设置超时
		Timeout: 30 * time.Second,
	}

	// 创建aws凭证
	credentials := credentials.NewStaticCredentialsProvider(config.AccessKeyID, config.SecretAccessKey, "")
	// 创建aws配置
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithCredentialsProvider(credentials),
		awsconfig.WithRegion(config.Region),
	)
	if err != nil {
		panic(err)
	}
	// 创建s3客户端
	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.HTTPClient = hcli
	})
	config.S3Client = client
	return config
}

func (c *AwsClient) WithBucket(bucket string) *AwsClient {
	c.Buckets = bucket
	return c
}

func (c *AwsClient) WithEndpoint(endpoint string) *AwsClient {
	c.Endpoint = endpoint
	return c
}

func (c *AwsClient) WithAccessKeyID(accessKeyID string) *AwsClient {
	c.AccessKeyID = accessKeyID
	return c
}

func (c *AwsClient) WithSecretAccessKey(secretAccessKey string) *AwsClient {
	c.SecretAccessKey = secretAccessKey
	return c
}

func (c *AwsClient) WithRegion(region string) *AwsClient {
	c.Region = region
	return c
}

// 直传文件
func (c *AwsClient) PutObject(ctx context.Context, key string, body io.Reader) error {

	return nil
}
