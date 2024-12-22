package manager

import (
	"crypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/registration"
)

type Domain struct {
	DomainName          string
	certificateResource *certificate.Resource
}

type User struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}
