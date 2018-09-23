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
	mimeTypePEM = "application/x-pem-file"
	// TagFingerprint is the key for the fingerprint meta data in a vault object
	TagFingerprint = "fingerprint"
	// TagKeyAlgo is the key for the PublicKeyAlgorithm meta data in a vault object
	TagKeyAlgo = "keyalgo"
	// TagNotAfter is the key for the expiry date meta data in a vault object
	TagNotAfter = "notafter"
	// TagSubjectCN is the key for the certificate subject common name meta data in a vault object
	TagSubjectCN = "subjectcn"
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

// GetCertificates returns a list of all certificates in vault
func GetCertificates() ([]*cert.Certificate, error) {
	vc, err := NewClient()

	if err != nil {
		return nil, err
	}

	var certs []*cert.Certificate

	req, err := vc.GetSecrets(context.Background(), viper.GetString("keyvault.url"), nil)

	if err != nil {
		return nil, err
	}

	for ; req.NotDone(); err = req.Next() {
		if err != nil {
			return nil, err
		}

		for _, item := range req.Values() {
			c, err := PullCertificate(*item.Tags[TagSubjectCN], *item.Tags[TagKeyAlgo])

			if err != nil {
				return nil, err
			}

			certs = append(certs, c)
		}
	}

	return certs, nil
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
			TagFingerprint: to.StringPtr(cert.Fingerprint()),
			TagKeyAlgo:     to.StringPtr(cert.PublicKeyAlgorithm()),
			TagNotAfter:    to.StringPtr(fmt.Sprintf("%d", cert.NotAfter().Unix())),
			TagSubjectCN:   to.StringPtr(cert.SubjectCN()),
		},
		ContentType: to.StringPtr(mimeTypePEM),
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
