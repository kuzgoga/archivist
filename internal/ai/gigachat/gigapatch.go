package gigachat

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"log"
	"net/http"
	"strings"
	"sync"
)

//go:embed russian_trusted_root_ca.cert
var rootCA []byte
var lock sync.Once

// patchDefaultClient Patch default transport and http client, if you are using gigachat LLM
func patchDefaultClient() {
	lock.Do(func() {
		sysTr := http.DefaultTransport

		minTr := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    makeCombinedPool(),
				MinVersion: tls.VersionTLS12,
			},
		}

		smartTr := &switchingTransport{
			defaultTransport:  sysTr,
			gigachatTransport: minTr,
		}

		http.DefaultClient = &http.Client{Transport: smartTr}
	})
}

var trustedHosts = []string{
	"ngw.devices.sberbank.ru",
	"gigachat.devices.sberbank.ru",
}

type switchingTransport struct {
	defaultTransport  http.RoundTripper
	gigachatTransport http.RoundTripper
}

func makeCombinedPool() *x509.CertPool {
	pool, err := x509.SystemCertPool()
	if err != nil || pool == nil {
		log.Printf("Warning: failed to load system cert pool: %v", err)
		pool = x509.NewCertPool()
	}

	if ok := pool.AppendCertsFromPEM(rootCA); !ok {
		log.Fatal("Failed to append MinCifra certificates")
	}

	return pool
}

func (s *switchingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	host := req.URL.Hostname()
	for _, trusted := range trustedHosts {
		if strings.HasSuffix(host, trusted) {
			return s.gigachatTransport.RoundTrip(req)
		}
	}
	return s.defaultTransport.RoundTrip(req)
}
