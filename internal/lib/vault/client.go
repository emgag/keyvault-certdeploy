package vault

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/emgag/keyvault-certdeploy/internal/lib/cert"
	"github.com/spf13/viper"
	"strings"
)

const (
	MIME_TYPE_PEM   = "application/x-pem-file"
	TAG_FINGERPRINT = "fingerprint"
	TAG_KEYALGO     = "keyalgo"
	TAG_NOTAFTER    = "notafter"
	TAG_SUBJECT     = "subject"
)

// NewClient creates a new key vault client
func NewClient() (*keyvault.BaseClient, error) {
	authorizer, err := auth.NewAuthorizerFromEnvironment()

	if err != nil {
		return nil, err
	}

	vc := keyvault.New()
	vc.Authorizer = authorizer

	return &vc, nil
}

// PushCertificate uploads a certificate to the key vault
func PushCertificate(cert *cert.Certificate) error {
	remote, err := PullCertificate(cert.SubjectCN(), cert.PublicKeyAlgorithm())

	if err == nil {
		if cert.Fingerprint() == remote.Fingerprint() {
			return errors.New("Certificate already in keyvault")
		}
	}

	vc, err := NewClient()

	if err != nil {
		return err
	}

	ssp := keyvault.SecretSetParameters{
		Tags: map[string]*string{
			TAG_FINGERPRINT: to.StringPtr(cert.Fingerprint()),
			TAG_KEYALGO:     to.StringPtr(cert.PublicKeyAlgorithm()),
			TAG_NOTAFTER:    to.StringPtr(fmt.Sprintf("%d", cert.NotAfter().Unix())),
			TAG_SUBJECT:     to.StringPtr(cert.SubjectCN()),
		},
		ContentType: to.StringPtr(MIME_TYPE_PEM),
		Value:       to.StringPtr(cert.String()),
	}

	_, err = vc.SetSecret(
		context.Background(),
		viper.GetString("keyvault.url"),
		CertificateIDFromCert(cert),
		ssp,
	)

	return err
}

// PullCertificate fetches a certificate from the key vault
func PullCertificate(subject string, keyalgo string) (*cert.Certificate, error) {
	vc, err := NewClient()

	if err != nil {
		return nil, err
	}

	sb, err := vc.GetSecret(
		context.Background(),
		viper.GetString("keyvault.url"),
		CertificateID(subject, keyalgo),
		"",
	)

	if err != nil {
		return nil, err
	}

	c, err := cert.Load(cert.Split([]byte(*sb.Value)))

	if err != nil {
		return nil, err
	}

	return c, nil
}

// CertificateID generates an object id to be used as an identifier for the cert in the vault
func CertificateID(subject string, keyalgo string) string {
	return fmt.Sprintf(
		"%s-%s",
		strings.Replace(subject, ".", "-", -1),
		strings.ToLower(keyalgo),
	)
}

// CertificateIDFromCert wraps CertificateID() for a certificate type
func CertificateIDFromCert(c *cert.Certificate) string {
	return CertificateID(c.SubjectCN(), c.PublicKeyAlgorithm())
}
