package certbot

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/ovh"
	"github.com/go-acme/lego/v4/registration"
)

type CertInfo struct {
	Domain     string    `json:"domain"`
	NotBefore  time.Time `json:"notBefore"`
	NotAfter   time.Time `json:"notAfter"`
	DaysLeft   int       `json:"daysLeft"`
	NeedsRenew bool      `json:"needsRenew"`
	CertPath   string    `json:"certPath"`
	KeyPath    string    `json:"keyPath"`
	Issuer     string    `json:"issuer"`
	Serial     string    `json:"serial"`
	DNSNames   []string  `json:"dnsNames"`
}

type CertRequest struct {
	Domain    string
	SANs      []string
	Email     string
	CertDir   string
	Staging   bool
	RenewDays int
}

type Certbot struct {
	ovhEndpoint     string
	ovhAppKey       string
	ovhAppSecret    string
	ovhConsumerKey  string
	ovhClientID     string
	ovhClientSecret string
}

func New(ovhEndpoint, appKey, appSecret, consumerKey string) *Certbot {
	return &Certbot{
		ovhEndpoint:    ovhEndpoint,
		ovhAppKey:      appKey,
		ovhAppSecret:   appSecret,
		ovhConsumerKey: consumerKey,
	}
}

func NewOAuth2(ovhEndpoint, clientID, clientSecret string) *Certbot {
	return &Certbot{
		ovhEndpoint:     ovhEndpoint,
		ovhClientID:     clientID,
		ovhClientSecret: clientSecret,
	}
}

type certUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *certUser) GetEmail() string {
	return u.Email
}

func (u *certUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *certUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func (c *Certbot) getOVHProvider() (*ovh.DNSProvider, error) {
	config := ovh.NewDefaultConfig()
	config.APIEndpoint = c.ovhEndpoint

	if c.ovhClientID != "" && c.ovhClientSecret != "" {
		config.OAuth2Config = &ovh.OAuth2Config{
			ClientID:     c.ovhClientID,
			ClientSecret: c.ovhClientSecret,
		}
	} else {
		config.ApplicationKey = c.ovhAppKey
		config.ApplicationSecret = c.ovhAppSecret
		config.ConsumerKey = c.ovhConsumerKey
	}

	return ovh.NewDNSProviderConfig(config)
}

func (c *Certbot) RequestCertificate(req CertRequest) (*CertInfo, error) {
	if err := os.MkdirAll(req.CertDir, 0700); err != nil {
		return nil, fmt.Errorf("création répertoire certificats: %w", err)
	}

	user, err := c.createUser(req.Email, req.CertDir)
	if err != nil {
		return nil, fmt.Errorf("création utilisateur ACME: %w", err)
	}

	legoConfig := lego.NewConfig(user)
	legoConfig.CADirURL = lego.LEDirectoryProduction
	if req.Staging {
		legoConfig.CADirURL = lego.LEDirectoryStaging
	}

	client, err := lego.NewClient(legoConfig)
	if err != nil {
		return nil, fmt.Errorf("création client ACME: %w", err)
	}

	ovhProvider, err := c.getOVHProvider()
	if err != nil {
		return nil, fmt.Errorf("configuration provider OVH: %w", err)
	}

	client.Challenge.SetDNS01Provider(ovhProvider,
		dns01.AddDNSTimeout(10*time.Minute),
		dns01.AddRecursiveNameservers([]string{"dns16.ovh.net:53", "ns16.ovh.net:53"}),
		dns01.RecursiveNSsPropagationRequirement(),
	)

	if user.Registration == nil {
		reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		if err != nil {
			return nil, fmt.Errorf("enregistrement ACME: %w", err)
		}
		user.Registration = reg
		_ = saveRegistration(user, req.CertDir)
	}

	domains := append([]string{req.Domain}, req.SANs...)
	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, fmt.Errorf("obtention certificat: %w", err)
	}

	domainSafe := strings.ReplaceAll(req.Domain, "*", "_")
	certPath := filepath.Join(req.CertDir, domainSafe+".crt")
	keyPath := filepath.Join(req.CertDir, domainSafe+".key")

	if err := os.WriteFile(certPath, certificates.Certificate, 0600); err != nil {
		return nil, fmt.Errorf("sauvegarde certificat: %w", err)
	}
	if err := os.WriteFile(keyPath, certificates.PrivateKey, 0600); err != nil {
		return nil, fmt.Errorf("sauvegarde clé privée: %w", err)
	}

	return c.GetCertInfo(req.Domain, req.CertDir, req.RenewDays)
}

func (c *Certbot) RenewCertificate(req CertRequest) (*CertInfo, error) {
	if err := os.MkdirAll(req.CertDir, 0700); err != nil {
		return nil, fmt.Errorf("création répertoire certificats: %w", err)
	}

	user, err := c.createUser(req.Email, req.CertDir)
	if err != nil {
		return nil, fmt.Errorf("création utilisateur ACME: %w", err)
	}

	legoConfig := lego.NewConfig(user)
	legoConfig.CADirURL = lego.LEDirectoryProduction
	if req.Staging {
		legoConfig.CADirURL = lego.LEDirectoryStaging
	}

	client, err := lego.NewClient(legoConfig)
	if err != nil {
		return nil, fmt.Errorf("création client ACME: %w", err)
	}

	ovhProvider, err := c.getOVHProvider()
	if err != nil {
		return nil, fmt.Errorf("configuration provider OVH: %w", err)
	}

	client.Challenge.SetDNS01Provider(ovhProvider,
		dns01.AddDNSTimeout(10*time.Minute),
		dns01.AddRecursiveNameservers([]string{"dns16.ovh.net:53", "ns16.ovh.net:53"}),
		dns01.RecursiveNSsPropagationRequirement(),
	)

	domainSafe := strings.ReplaceAll(req.Domain, "*", "_")
	certPath := filepath.Join(req.CertDir, domainSafe+".crt")
	keyPath := filepath.Join(req.CertDir, domainSafe+".key")

	certData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("lecture certificat: %w", err)
	}
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("lecture clé privée: %w", err)
	}

	certificates, err := client.Certificate.Renew(certificate.Resource{
		Domain:      req.Domain,
		Certificate: certData,
		PrivateKey:  keyData,
	}, true, false, "")
	if err != nil {
		return nil, fmt.Errorf("renouvellement certificat: %w", err)
	}

	if err := os.WriteFile(certPath, certificates.Certificate, 0600); err != nil {
		return nil, fmt.Errorf("sauvegarde certificat: %w", err)
	}
	if err := os.WriteFile(keyPath, certificates.PrivateKey, 0600); err != nil {
		return nil, fmt.Errorf("sauvegarde clé privée: %w", err)
	}

	return c.GetCertInfo(req.Domain, req.CertDir, req.RenewDays)
}

func (c *Certbot) GetCertInfo(domain, certDir string, renewDays int) (*CertInfo, error) {
	domainSafe := strings.ReplaceAll(domain, "*", "_")
	certPath := filepath.Join(certDir, domainSafe+".crt")
	keyPath := filepath.Join(certDir, domainSafe+".key")

	certData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("lecture certificat: %w", err)
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, fmt.Errorf("décodage PEM échoué")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing certificat: %w", err)
	}

	daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)
	needsRenew := daysLeft <= renewDays

	return &CertInfo{
		Domain:     domain,
		NotBefore:  cert.NotBefore,
		NotAfter:   cert.NotAfter,
		DaysLeft:   daysLeft,
		NeedsRenew: needsRenew,
		CertPath:   certPath,
		KeyPath:    keyPath,
		Issuer:     cert.Issuer.CommonName,
		Serial:     cert.SerialNumber.String(),
		DNSNames:   cert.DNSNames,
	}, nil
}

func (c *Certbot) ListCertificates(certDir string) ([]CertInfo, error) {
	entries, err := os.ReadDir(certDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []CertInfo{}, nil
		}
		return nil, fmt.Errorf("lecture répertoire: %w", err)
	}

	var certs []CertInfo
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".crt" {
			domain := strings.TrimSuffix(entry.Name(), ".crt")
			domain = strings.ReplaceAll(domain, "_", "*")

			info, err := c.GetCertInfo(domain, certDir, 30)
			if err != nil {
				continue
			}
			certs = append(certs, *info)
		}
	}

	return certs, nil
}

func (c *Certbot) createUser(email, certDir string) (*certUser, error) {
	keyPath := filepath.Join(certDir, ".acme-user.key")
	regPath := filepath.Join(certDir, ".acme-user-reg.json")

	keyData, err := os.ReadFile(keyPath)
	if err == nil {
		block, _ := pem.Decode(keyData)
		if block != nil {
			privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err == nil {
				user := &certUser{Email: email, key: privKey}
				if regData, err := os.ReadFile(regPath); err == nil {
					var reg registration.Resource
					if json.Unmarshal(regData, &reg) == nil {
						user.Registration = &reg
					}
				}
				return user, nil
			}
		}
	}

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("génération clé privée: %w", err)
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(privKey)

	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	})

	if err := os.WriteFile(keyPath, pemKey, 0600); err != nil {
		return nil, fmt.Errorf("sauvegarde clé utilisateur: %w", err)
	}

	return &certUser{Email: email, key: privKey}, nil
}

func saveRegistration(user *certUser, certDir string) error {
	if user.Registration == nil {
		return nil
	}
	regPath := filepath.Join(certDir, ".acme-user-reg.json")
	data, err := json.Marshal(user.Registration)
	if err != nil {
		return err
	}
	return os.WriteFile(regPath, data, 0600)
}
