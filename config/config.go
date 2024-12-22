package config

import (
	"encoding/json"
	"os"
)

var (
	conf *Conf

	Persistence  *PersistenceType
	DNSChallenge *DNSChallengeType
	CA           *CAType
	User         *CertUserType
	CDN          *CDNType

	Path string = "./config.json"
)

func init() {
	_config, err := os.ReadFile(Path)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(_config, &conf); err != nil {
		panic(err)
	}

	Persistence = &conf.Persistence
	DNSChallenge = &conf.DNSChallenge
	CA = &conf.CA
	User = &conf.CertUser
	CDN = &conf.CDN
}
