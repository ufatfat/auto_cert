package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"os"
	"os/exec"
	"strings"
)

const (
	CertificateFilePath = "/var/snap/lego/common/.lego/certificates/"
)

var entries map[string]cron.EntryID
var timer *cron.Cron

func main() {
	timer = cron.New()
	entries = make(map[string]cron.EntryID)
	e := initRouter()
	e.Run(":8081")
}

func challenge(c *gin.Context) {
	chal := c.Param("challenge")
	data, err := os.ReadFile("./challenges/.well-known/acme-challenge/" + chal)
	if err != nil {
		panic(err)
	}
	c.Header("Content-Type", "text/plain")
	c.String(200, string(data))
}
func obtain(c *gin.Context) {
	dn := c.Query("dn")
	if dn == "" {
		c.Status(400)
		return
	}

	cmd := exec.Command("lego", "--email", "lzff@ufatfat.com", "--domains", dn, "--dns", "alidns", "run")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		return
	}
	if !strings.Contains(string(out), "Server responded with a certificate.") {
		c.Status(500)
		fmt.Printf("Failed to obtain certificate:\n%s\n", string(out))
		// todo 失败处理
		return
	}
	e, err := timer.AddFunc("@every 1800h", func() {
		renew(dn)
	})
	if err != nil {
		c.Status(500)
		fmt.Println("Failed to create cron job of renewing certificates of domain name", dn)
		return
	}
	entries[dn] = e

	updateCDNCert(dn)
}
func _renew(c *gin.Context) {
	dn := c.Query("dn")
	if dn == "" {
		c.Status(400)
		return
	}

	renew(dn)
}

func initRouter() *gin.Engine {
	router := gin.Default()

	root := router.Group("/")
	{
		root.GET(".well-known/acme-challenge/:challenge", challenge)
		root.POST("obtain", obtain)
		root.POST("renew", _renew)
	}

	return router
}

func renew(dn string) {
	cmd := exec.Command("lego", "--email", "lzff@ufatfat.com", "--domains", dn, "--dns", "alidns", "renew")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
	}
	if !strings.Contains(string(out), "Server responded with a certificate.") {
		// todo 失败处理
	}

	updateCDNCert(dn)
}

func updateCDNCert(dn string) {
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
