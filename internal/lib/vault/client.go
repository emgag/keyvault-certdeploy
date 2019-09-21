package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/emgag/keyvault-certdeploy/internal/lib/cert"
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

type Client struct {
	VaultClient *keyvault.BaseClient
	VaultURL    string
}

// GetCertificates returns a list of all certificates in vault
func (c *Client) GetCertificates() ([]*cert.Certificate, error) {
	var certs []*cert.Certificate

	req, err := c.VaultClient.GetSecrets(context.Background(), c.VaultURL, nil)

	if err != nil {
		return nil, err
	}

	for ; req.NotDone(); err = req.Next() {
		if err != nil {
			return nil, err
		}

		for _, item := range req.Values() {
			c, err := c.PullCertificate(*item.Tags[TagSubjectCN], *item.Tags[TagKeyAlgo])

			if err != nil {
				return nil, err
			}

			certs = append(certs, c)
		}
	}

	return certs, nil
}

// PushCertificate uploads a certificate to the key vault
func (c *Client) PushCertificate(cert *cert.Certificate) error {
	remote, err := c.PullCertificate(cert.SubjectCN(), cert.PublicKeyAlgorithm())

	if err == nil {
		if cert.Fingerprint() == remote.Fingerprint() {
			return errors.New("Certificate already in keyvault")
		}
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

	_, err = c.VaultClient.SetSecret(
		context.Background(),
		c.VaultURL,
		CertificateIDFromCert(cert),
		ssp,
	)

	return err
}

// PullCertificate fetches a certificate from the key vault
func (c *Client) PullCertificate(subject string, keyalgo string) (*cert.Certificate, error) {
	sb, err := c.VaultClient.GetSecret(
		context.Background(),
		c.VaultURL,
		CertificateID(subject, keyalgo),
		"",
	)

	if err != nil {
		return nil, err
	}

	crt, err := cert.Load(cert.Split([]byte(*sb.Value)))

	if err != nil {
		return nil, err
	}

	return crt, nil
}

// DeleteCertificate deletes a certificate from vault
func (c *Client) DeleteCertificate(subject string, keyalgo string) error {
	_, err := c.VaultClient.DeleteSecret(
		context.Background(),
		c.VaultURL,
		CertificateID(subject, keyalgo),
	)

	if err != nil {
		return err
	}

	return nil
}

// NewClient creates a new key vault client
func NewClient(vaultURL string) (*Client, error) {
	authorizer, err := auth.NewAuthorizerFromEnvironment()

	if err != nil {
		return nil, err
	}

	vc := keyvault.New()
	vc.Authorizer = authorizer

	c := &Client{
		VaultURL:    vaultURL,
		VaultClient: &vc,
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
