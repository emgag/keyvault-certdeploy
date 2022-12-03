package cert

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// Certificate holds the parsed and the PEM-encoded raw key & cert
type Certificate struct {
	Certificate *tls.Certificate
	Leaf        *x509.Certificate
	Chain       []*x509.Certificate
	RawKey      []byte
	RawCert     []byte
}

// SubjectCN returns the subject CN of the leaf certificate
func (c *Certificate) SubjectCN() string {
	return c.Leaf.Subject.CommonName
}

// NotAfter returns the expire date of the leaf certificate
func (c *Certificate) NotAfter() time.Time {
	return c.Leaf.NotAfter
}

// PublicKeyAlgorithm returns the public key algorithm of the leaf certificate
func (c *Certificate) PublicKeyAlgorithm() string {
	return c.Leaf.PublicKeyAlgorithm.String()
}

// Fingerprint returns the sha256 fingerprint of the leaf certificate
func (c *Certificate) Fingerprint() string {
	return fmt.Sprintf("%x", sha256.Sum256(c.Leaf.Raw))
}

// LeafPEM returns the raw leaf cert as a pem encoded byte slice
func (c *Certificate) LeafPEM() []byte {
	b := pem.Block{Type: "CERTIFICATE", Bytes: c.Leaf.Raw}
	return pem.EncodeToMemory(&b)
}

// ChainPEM returns the raw certificate chain certificates as a pem encoded byte slice
func (c *Certificate) ChainPEM() []byte {
	buf := new(bytes.Buffer)

	for _, cc := range c.Chain {
		b := pem.Block{Type: "CERTIFICATE", Bytes: cc.Raw}
		buf.Write(pem.EncodeToMemory(&b))
	}

	return buf.Bytes()
}

// FullPEM returns the full certificate chain and private key as a pem encoded byte slice
func (c *Certificate) FullPEM() []byte {
	buf := new(bytes.Buffer)
	buf.Write(c.RawCert)
	buf.Write(c.RawKey)
	return buf.Bytes()
}

// String returns the private key concatenated with the whole certificate chain as a PEM-encoded string
func (c *Certificate) String() string {
	sb := strings.Builder{}
	sb.Write(c.RawKey)
	sb.Write(c.RawCert)
	return sb.String()
}

// LoadFromDisk returns a certificate loaded from a key- and certificate chain file
func LoadFromDisk(keyFile string, certFile string) (*Certificate, error) {
	key, err := os.ReadFile(keyFile)

	if err != nil {
		return nil, err
	}

	cert, err := os.ReadFile(certFile)

	if err != nil {
		return nil, err
	}

	return Load(key, cert)
}

// Load returns a certificate from a key- and certificate chain byte slice
func Load(key []byte, cert []byte) (*Certificate, error) {
	certs, err := tls.X509KeyPair(cert, key)

	if err != nil {
		return nil, err
	}

	var leaf *x509.Certificate
	var chain []*x509.Certificate

	for _, c := range certs.Certificate {
		c, err := x509.ParseCertificate(c)

		if err == nil {
			if c.IsCA {
				chain = append(chain, c)
			} else {
				leaf = c
			}
		}
	}

	if leaf == nil {
		return nil, errors.New("No leaf certificate found")
	}

	return &Certificate{
		Certificate: &certs,
		Leaf:        leaf,
		Chain:       chain,
		RawKey:      key,
		RawCert:     cert,
	}, nil
}

// Split returns a separate byte slice for the key and certificate chain from a single combined input
func Split(pemData []byte) (key, cert []byte) {
	k := new(bytes.Buffer)
	c := new(bytes.Buffer)

	for {
		block, rest := pem.Decode(pemData)

		if block == nil {
			break
		}

		if block.Type == "CERTIFICATE" {
			c.Write(pem.EncodeToMemory(block))
		} else {
			k.Write(pem.EncodeToMemory(block))
		}

		pemData = rest
	}

	return k.Bytes(), c.Bytes()
}
