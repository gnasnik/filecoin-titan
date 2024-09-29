package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/Filecoin-Titan/titan/node/config"
	"github.com/Filecoin-Titan/titan/node/repo"
)

var tlsCfg *tls.Config

func fetchTlsConfigFromRemote(acmeAddress string) (cfg *tls.Config, domainSuffix string) {

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
		NextProtos: []string{"h2", "h3"},
		// InsecureSkipVerify: true,
		// GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
		// 	return nil, nil
		// },
		// VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		// 	return nil
		// },
		// GetConfigForClient: func(chi *tls.ClientHelloInfo) (*tls.Config, error) {
		// 	return nil, nil
		// },
	}

	type Acme struct {
		Certificate string    `json:"certificate"`
		PrivateKey  string    `json:"private_key"`
		CreatedAt   time.Time `json:"created_at"`
		ExpireAt    time.Time `json:"expire_at"`
	}
	resp, err := http.Get(acmeAddress)
	if err != nil {
		log.Errorf("fetch tls config failed, error:%s", err.Error())
		return nil, ""
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorf("tls server error, code:%s", resp.Status)
		return nil, ""
	}

	var acme Acme
	err = json.NewDecoder(resp.Body).Decode(&acme)
	if err != nil {
		log.Errorf("read tls config failed, error:%s", err.Error())
		return nil, ""
	}

	cert, err := tls.X509KeyPair([]byte(acme.Certificate), []byte(acme.PrivateKey))
	if err != nil {
		log.Errorf("load tls config failed, error:%s", err.Error())
		return nil, ""
	}

	if len(cert.Certificate) == 0 {
		log.Error("no certificate found in the provided tls.CertificatePath")
		return nil, ""
	}

	parsedCert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		log.Errorf("parse certificate failed, error:%s", err.Error())
		return nil, ""
	}

	if parsedCert.NotAfter.Before(time.Now()) {
		log.Error("remote certificate was expired")
		return nil, ""
	}

	for _, v := range parsedCert.DNSNames {
		if strings.HasPrefix(v, "*") {
			domainSuffix = v
			tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
		}
	}

	if err := mergeLocalTlsConfig(tlsConfig); err != nil {
		log.Error("merge local tls config error", err)
		return nil, ""
	}

	tlsCfg = tlsConfig

	return tlsConfig, domainSuffix
}

func refreshTlsConfig(acmeAddress string, lr repo.LockedRepo, candidateCfg *config.CandidateCfg) {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		tlsConfig, suffix := fetchTlsConfigFromRemote(acmeAddress)
		if suffix == "" {
			log.Error("fail to refresh tls config, remote host not exist, will use default config")
			continue
		}
		log.Infof("refreshing tls config with suffix:%s", suffix)
		if err := flushConfig(lr, tlsConfig, candidateCfg); err != nil {
			log.Errorf("flush config: %v", err)
		}
	}
}

func mergeLocalTlsConfig(tlsConfig *tls.Config) error {
	localCfg, err := defaultTLSConfig()
	if err != nil {
		return err
	}
	tlsConfig.Certificates = append(tlsConfig.Certificates, localCfg.Certificates...)
	return nil
}

func defaultTLSConfig() (*tls.Config, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		// DNSNames:     []string{"localhost"},
		// IPAddresses:  []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("0.0.0.0")},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		Certificates:       []tls.Certificate{tlsCert},
		NextProtos:         []string{"h2", "h3"},
		InsecureSkipVerify: true, //nolint:gosec // skip verify in default config
	}, nil
}

func flushConfig(lr repo.LockedRepo, tlsConfig *tls.Config, cfg *config.CandidateCfg) error {
	// TODO flush cert/key pair to external to disk config
	if tlsConfig == nil {
		return nil
	}

	// Check if there are any certificates
	if len(tlsConfig.Certificates) > 0 {
		cert := tlsConfig.Certificates[0]

		var certBytes []byte
		for _, c := range cert.Certificate {
			certBytes = append(certBytes, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c})...)
		}

		var certKeyBytes []byte

		privBytes, err := x509.MarshalECPrivateKey(cert.PrivateKey.(*ecdsa.PrivateKey))
		if err == nil {
			certKeyBytes = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
		} else {
			privBytes, err = x509.MarshalPKCS8PrivateKey(cert.PrivateKey)
			if err != nil {
				return err
			}
			certKeyBytes = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
		}

		if err := lr.SetCertificate(certBytes, certKeyBytes); err != nil {
			log.Errorf("SetCertificate: %v", err)
		}

	}

	return lr.SetConfig(func(raw interface{}) {
		scfg, ok := raw.(*config.CandidateCfg)
		if !ok {
			return
		}
		scfg.ExternalURL = cfg.ExternalURL
		scfg.IngressHostName = cfg.IngressHostName
		scfg.IngressCertificatePath = cfg.IngressCertificatePath
		scfg.IngressCertificateKeyPath = cfg.IngressCertificateKeyPath
	})
}