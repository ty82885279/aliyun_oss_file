package main

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"net/http"
)

type aliyun = string

const (
	Region          aliyun = "cn-shanghai"
	AccessKeyId     aliyun = "LTAI4FyRs7GZ3vEqecvyDXiK"
	AccessKeySecret aliyun = "aY7gchBnav9AACmoxXEgW3RHkZQZNM"
	RoleArn         aliyun = "acs:ram::1934830090765555:role/u1"
	RoleSessionName aliyun = "test"
)

func getAliyunToken(c *gin.Context) {
	//构建一个阿里云客户端, 用于发起请求。
	//构建阿里云客户端时，需要设置AccessKey ID和AccessKey Secret。
	client, err := sts.NewClientWithAccessKey(Region, AccessKeyId, AccessKeySecret)

	//构建请求对象。
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	//设置参数。关于参数含义和设置方法，请参见API参考。
	request.RoleArn = RoleArn
	request.RoleSessionName = RoleSessionName

	response, err := client.AssumeRole(request)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"StatusCode":   500,
			"ErrorCode":    "InvalidAccessKeyId.NotFound",
			"ErrorMessage": "Specified access key is not found.",
		})
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"StatusCode":      200,
		"AccessKeyId":     response.Credentials.AccessKeyId,
		"AccessKeySecret": response.Credentials.AccessKeySecret,
		"SecurityToken":   response.Credentials.SecurityToken,
		"Expiration":      response.Credentials.Expiration,
	})
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:8003",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}

func callBack(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func main() {

	r := gin.Default()
	r.Use(TlsHandler())
	r.GET("token", getAliyunToken)
	r.POST("callback", callBack)
	r.RunTLS(":8003", "./cert/1_lee2code.com_bundle.crt", "./cert/2_lee2code.com.key")
	//r.Run(":8003")

}
