package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
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

// Client represents a certdeploy client
type Client struct {
	VaultClient *azsecrets.Client
	VaultURL    string
}

// GetCertificates returns a list of all certificates in vault
func (c *Client) GetCertificates() ([]*cert.Certificate, error) {
	var certs []*cert.Certificate

	pager := c.VaultClient.NewListSecretsPager(nil)

	for pager.More() {
		page, err := pager.NextPage(context.Background())

		if err != nil {
			return nil, err
		}

		for _, secret := range page.Value {
			c, err := c.PullCertificate(*secret.Tags[TagSubjectCN], *secret.Tags[TagKeyAlgo])

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

	stringPtr := func(s string) *string { return &s }

	ssp := azsecrets.SetSecretParameters{
		Tags: map[string]*string{
			TagFingerprint: stringPtr(cert.Fingerprint()),
			TagKeyAlgo:     stringPtr(cert.PublicKeyAlgorithm()),
			TagNotAfter:    stringPtr(fmt.Sprintf("%d", cert.NotAfter().Unix())),
			TagSubjectCN:   stringPtr(cert.SubjectCN()),
		},
		ContentType: stringPtr(mimeTypePEM),
		Value:       stringPtr(cert.String()),
	}

	_, err = c.VaultClient.SetSecret(
		context.Background(),
		CertificateIDFromCert(cert),
		ssp,
		nil,
	)

	return err
}

// PullCertificate fetches a certificate from the key vault
func (c *Client) PullCertificate(subject string, keyalgo string) (*cert.Certificate, error) {
	sb, err := c.VaultClient.GetSecret(
		context.Background(),
		CertificateID(subject, keyalgo),
		"",
		nil,
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
		CertificateID(subject, keyalgo),
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

// NewClient creates a new vault client
func NewClient(vaultURL string) (*Client, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	vc, err := azsecrets.NewClient(vaultURL, cred, nil)

	if err != nil {
		return nil, err
	}

	c := &Client{
		VaultURL:    vaultURL,
		VaultClient: vc,
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
