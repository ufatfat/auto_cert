package manager

import (
	"errors"
	"fmt"
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

func (d *Domain) SetCDNClient(cdnClient CDNClient) {
	d.CDNClient = cdnClient
}
