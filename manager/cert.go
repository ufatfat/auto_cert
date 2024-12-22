package manager

import (
	"auto_cert/config"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/go-acme/lego/v4/certificate"
	"os"
)

func (d *Domain) Obtain() (err error) {
	obtainRequest := certificate.ObtainRequest{
		Domains: []string{d.DomainName},
		Bundle:  true,
	}

	certificateResource, err := client.Certificate.Obtain(obtainRequest)
	if err != nil {
		return
	}
	d.certificateResource = certificateResource

	return nil
}

func (d *Domain) Renew() (err error) {
	if d.certificateResource == nil || d.certificateResource.Certificate == nil {
		return errors.New("no certificate found, can not renew")
	}
	certificateResource, err := client.Certificate.RenewWithOptions(*d.certificateResource, &certificate.RenewOptions{
		Bundle: true,
	})
	if err != nil {
		return
	}
	d.certificateResource = certificateResource
	return nil
}

func (d *Domain) GetCertificate() (cert []byte, err error) {
	return d.certificateResource.Certificate, nil
}

func (d *Domain) GetPrivateKey() (key []byte, err error) {
	return d.certificateResource.PrivateKey, nil
}

func (d *Domain) WriteCertificate(path string) (err error) {
	return os.WriteFile(fmt.Sprintf("./%s/%s.crt", path, d.DomainName), d.certificateResource.Certificate, 0755)
}

func (d *Domain) WritePrivateKey(path string) (err error) {
	return os.WriteFile(fmt.Sprintf("./%s/%s.key", path, d.DomainName), d.certificateResource.PrivateKey, 0755)
}

func (d *Domain) UpdateCDN() (err error) {
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
	request.QueryParams["SSLPub"] = string(d.certificateResource.Certificate)
	request.QueryParams["SSLPri"] = string(d.certificateResource.PrivateKey)
	request.QueryParams["DomainName"] = d.DomainName

	response, err := cdnClient.ProcessCommonRequest(request)
	if err != nil {
		return
	}
	if !response.IsSuccess() {
		return errors.New(response.GetHttpContentString())
	}
	return
}
