package ai

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"log"
	"net/http"
	"sync"
)

//go:embed russian_trusted_root_ca.cert
var rootCA []byte
var lock sync.Once

// patchDefaultTransport
// Patch default transport and http client, if you are using gigachat LLM
func patchDefaultTransport() {
	lock.Do(func() {
		sysPool, err := x509.SystemCertPool()
		if err != nil {
			log.Printf("Warning: failed to load system cert pool: %v; creating a new one", err)
			sysPool = x509.NewCertPool()
		}
		if sysPool == nil {
			sysPool = x509.NewCertPool()
		}

		if ok := sysPool.AppendCertsFromPEM(rootCA); !ok {
			log.Fatal("Failed to append root CA certificates to the pool")
		}

		tlsCfg := &tls.Config{
			RootCAs:    sysPool,
			MinVersion: tls.VersionTLS12,
		}

		customTransport := &http.Transport{
			TLSClientConfig: tlsCfg,
		}

		http.DefaultTransport = customTransport
		http.DefaultClient = &http.Client{Transport: customTransport}
	})

}
