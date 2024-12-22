package util

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"os"
)

func UpdateCDNCert(dn string) (err error) {
	config := sdk.NewConfig()

	// Please ensure that the environment variables ALIBABA_CLOUD_ACCESS_KEY_ID and ALIBABA_CLOUD_ACCESS_KEY_SECRET are set.
	credential := credentials.NewAccessKeyCredential(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"))
	client, err := sdk.NewClientWithOptions("cn-qingdao", config, credential)
	if err != nil {
		panic(err)
	}
	if dn[0] == '*' {
		dn = "_" + dn[1:]
	}

	request := requests.NewCommonRequest()

	pub, err := os.ReadFile(CertificateFilePath + dn + ".crt")
	priv, err := os.ReadFile(CertificateFilePath + dn + ".key")

	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "cdn.aliyuncs.com"
	request.Version = "2018-05-10"
	request.ApiName = "SetCdnDomainSSLCertificate"
	request.QueryParams["SSLProtocol"] = "on"
	request.QueryParams["SSLPub"] = string(pub)
	request.QueryParams["SSLPri"] = string(priv)
	request.QueryParams["DomainName"] = dn

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpContentString())
}
