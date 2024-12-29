package manager

import (
	"auto_cert/config"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"github.com/go-acme/lego/v4/registration"
	"os"
)

var client *lego.Client

//var dnsProvider *alidns.DNSProvider

func init() {
	var err error
	var privateKey = new(crypto.PrivateKey)

	if config.User.Key.Value != "" {
		var ecdsaPrivateKey *ecdsa.PrivateKey
		switch config.User.Key.Type {
		case "string":
			ecdsaPrivateKey, err = x509.ParseECPrivateKey([]byte(config.User.Key.Value))
			if err != nil {
				panic(err)
			}
		case "file":
			_key, err := os.ReadFile(config.User.Key.Value)
			if err != nil {
				panic(err)
			}
			block, _ := pem.Decode(_key)
			ecdsaPrivateKey, err = x509.ParseECPrivateKey(block.Bytes)
			if err != nil {
				panic(err)
			}
		}
		*privateKey = crypto.PrivateKey(ecdsaPrivateKey)
	} else {
		var ecdsaPrivateKey *ecdsa.PrivateKey
		ecdsaPrivateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			panic(err)
		}
		keyData, err := x509.MarshalECPrivateKey(ecdsaPrivateKey)
		if err != nil {
			panic(err)
		}

		keyFile, err := os.Create("./private.key")
		if err != nil {
			panic(err)
		}
		defer keyFile.Close()

		block := pem.Block{
			Type:    "EC PRIVATE KEY",
			Headers: nil,
			Bytes:   keyData,
		}
		if err = pem.Encode(keyFile, &block); err != nil {
			panic(err)
		}

		fmt.Printf("ECDSA private key generated, key file written to: ./private.key\n")

		*privateKey = crypto.PrivateKey(ecdsaPrivateKey)
	}

	user := User{
		Email: config.User.Email,
		key:   *privateKey,
	}

	legoConfig := lego.NewConfig(&user)
	legoConfig.CADirURL = config.CA.URL
	legoConfig.Certificate.KeyType = certcrypto.RSA2048

	client, err = lego.NewClient(legoConfig)
	if err != nil {
		panic(err)
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		panic(err)
	}
	user.Registration = reg

	dnsConfig := alidns.NewDefaultConfig()
	dnsConfig.APIKey = config.DNSChallenge.APIKey
	dnsConfig.SecretKey = config.DNSChallenge.SecretKey

	dnsProvider, err := alidns.NewDNSProviderConfig(dnsConfig)
	if err != nil {
		panic(err)
	}
	if err = client.Challenge.SetDNS01Provider(dnsProvider); err != nil {
		panic(err)
	}
}

func NewDomain(dn string) *Domain {
	return &Domain{
		DomainName: dn,
	}
}

func NewDomainWithConfig(dn string, config ...DomainConfig) *Domain {
	return &Domain{
		DomainName: dn,
	}
}
