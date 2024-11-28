package sdk

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gabriel-vasile/mimetype"
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
func (c *AwsClient) PutObject(ctx context.Context, fid *string, body io.Reader) error {
	// 读取文件前512字节 判断文件的类型
	header := make([]byte, 512)
	l, err := io.ReadFull(body, header)
	if err != nil {
		return fmt.Errorf("读取文件头部失败: %w", err)
	}
	contentType := mimetype.Detect(header[:l]).String()
	if strings.HasPrefix(contentType, "video") {
		// 给文件路径添加前缀
		fid = aws.String(fmt.Sprintf("video/%s", aws.ToString(fid)))
	} else if strings.HasPrefix(contentType, "image") {
		fid = aws.String(fmt.Sprintf("image/%s", aws.ToString(fid)))
	} else if strings.HasPrefix(contentType, "audio") {
		fid = aws.String(fmt.Sprintf("audio/%s", aws.ToString(fid)))
	} else {
		fid = aws.String(fmt.Sprintf("file/%s", aws.ToString(fid)))
	}
	// 先计算文件的md5值
	md5hash := md5.New()
	tee := io.TeeReader(io.MultiReader(bytes.NewReader(header[:l]), body), md5hash) // 使用 TeeReader 避免重复读取
	// 使用fid作为文件名
	_, err = c.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:     aws.String(c.Buckets),
		Key:        fid,
		Body:       tee,
		ContentMD5: aws.String(hex.EncodeToString(md5hash.Sum(nil))),
		// 可以根据文件扩展名设置 ContentType
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("上传文件失败: %w", err)
	}
	return nil
}

func (c *AwsClient) DisplayUrl(ctx context.Context, fid *string, contentType string) string {
	path := fmt.Sprintf("%s/%s/%s/%s", c.Endpoint, c.Buckets, strings.Split(contentType, "/")[0], aws.ToString(fid))
	return path
}
