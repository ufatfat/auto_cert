package alicdn

import (
	"auto_cert/config"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type CDN struct {
}

// UpdateCDN implements manager.CDNClient interface
func (c *CDN) UpdateCDN(certificate, privateKey, domainName string) (err error) {

	sdkConfig := sdk.NewConfig()

	credential := credentials.NewAccessKeyCredential(config.CDN.AccessKey, config.CDN.SecretKey)
	cdnClient, err := sdk.NewClientWithOptions("cn-qingdao", sdkConfig, credential)
	if err != nil {
		return
	}

	request := requests.NewCommonRequest()

	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "cdn.aliyuncs.com"
	request.Version = "2018-05-10"
	request.ApiName = "SetCdnDomainSSLCertificate"
	request.QueryParams["SSLProtocol"] = "on"
	request.QueryParams["SSLPub"] = certificate
	request.QueryParams["SSLPri"] = privateKey
	request.QueryParams["DomainName"] = domainName

	response, err := cdnClient.ProcessCommonRequest(request)
	if err != nil {
		return
	}
	if !response.IsSuccess() {
		return errors.New(response.GetHttpContentString())
	}
	return
}
