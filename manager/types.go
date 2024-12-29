package manager

import (
	"crypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/registration"
)

type CDNClient interface {
	UpdateCDN(string, string, string) error
}

type Domain struct {
	DomainName          string
	certificateResource *certificate.Resource
	CDNClient
}

type User struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

type DomainConfig struct{}
