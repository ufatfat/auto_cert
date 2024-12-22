package config

type Conf struct {
	Persistence  PersistenceType  `json:"persistence" yaml:"persistence"`
	DNSChallenge DNSChallengeType `json:"dns_challenge"`
	CA           CAType           `json:"ca" yaml:"ca"`
	CertUser     CertUserType     `json:"cert_user" yaml:"cert_user"`
	CDN          CDNType          `json:"cdn" yaml:"cdn"`
}

type PersistenceType struct {
	Path string `json:"path" yaml:"path"`
}
type DNSChallengeType struct {
	APIKey    string `json:"api_key" yaml:"api_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
}
type CAType struct {
	URL string `json:"url" yaml:"url"`
}
type CertUserType struct {
	Email string `json:"email" yaml:"email"`
	Key   struct {
		Type  string `json:"type" yaml:"type"`
		Value string `json:"value" yaml:"value"`
	} `json:"key" yaml:"key"`
}
type CDNType struct {
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
}
