package winrm

import "fmt"

type Auth string

var (
	Basic Auth = "basic"
	Cert  Auth = "cert"
)

type Endpoint struct {
	Host     string
	Port     int
	HTTPS    bool
	Insecure bool
	Auth     Auth
	CACert   *[]byte
	Cert     []byte
	Key      []byte
}

func (ep *Endpoint) url() string {
	var scheme string
	if ep.HTTPS {
		scheme = "https"
	} else {
		scheme = "http"
	}

	return fmt.Sprintf("%s://%s:%d/wsman", scheme, ep.Host, ep.Port)
}
