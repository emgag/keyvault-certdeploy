package cert_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/emgag/keyvault-certdeploy/internal/lib/cert"
)

type certInfo struct {
	KeyFile            string
	KeyData            []byte
	CertFile           string
	CertData           []byte
	ChainFile          string
	ChainData          []byte
	FullChainFile      string
	FullChainData      []byte
	FullChainKeyFile   string
	FullChainKeyData   []byte
	SubjectCN          string
	NotAfter           time.Time
	PublicKeyAlgorithm string
	Fingerprint        string
}

var testData = "../../../testdata"

var testCerts = []*certInfo{
	{
		SubjectCN:          "c1.example.org",
		PublicKeyAlgorithm: "rsa",
		NotAfter:           time.Date(2027, 12, 03, 17, 19, 0, 0, time.UTC),
		Fingerprint:        "e6473ebc403e4a6f664dcbb98eb6886dd92232382ac641b06a731e737c586392",
	}, {
		SubjectCN:          "c1.example.org",
		PublicKeyAlgorithm: "ecdsa",
		NotAfter:           time.Date(2027, 12, 03, 17, 18, 0, 0, time.UTC),
		Fingerprint:        "ee27bb29e9563efd2830902b091e084fc55f797334cc04ede14d6daec7fae005",
	}, {
		SubjectCN:          "c2.example.org",
		PublicKeyAlgorithm: "rsa",
		NotAfter:           time.Date(2017, 6, 13, 16, 05, 0, 0, time.UTC),
		Fingerprint:        "cecb46eb661418e3dafc9d2ee21e7f30865862b3c0823c9466f9678026b55026",
	}, {
		SubjectCN:          "c2.example.org",
		PublicKeyAlgorithm: "ecdsa",
		NotAfter:           time.Date(2017, 6, 13, 16, 11, 0, 0, time.UTC),
		Fingerprint:        "be74f5f0ba56cec82d6371898b179a3790c1dca657e1bea78ef9555ba37e0e0d",
	},
}

func setup() {
	for _, c := range testCerts {
		c.KeyFile = filepath.Join(testData, fmt.Sprintf("%s.privkey.%s.pem", c.SubjectCN, c.PublicKeyAlgorithm))
		c.CertFile = filepath.Join(testData, fmt.Sprintf("%s.cert.%s.pem", c.SubjectCN, c.PublicKeyAlgorithm))
		c.FullChainFile = filepath.Join(testData, fmt.Sprintf("%s.fullchain.%s.pem", c.SubjectCN, c.PublicKeyAlgorithm))
		c.FullChainKeyFile = filepath.Join(testData, fmt.Sprintf("%s.fullchain.key.%s.pem", c.SubjectCN, c.PublicKeyAlgorithm))
		c.ChainFile = filepath.Join(testData, "chain.pem")
		c.KeyData, _ = os.ReadFile(c.KeyFile)
		c.CertData, _ = os.ReadFile(c.CertFile)
		c.ChainData, _ = os.ReadFile(c.ChainFile)
		c.FullChainData, _ = os.ReadFile(c.FullChainFile)
		c.FullChainKeyData, _ = os.ReadFile(c.FullChainKeyFile)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	//shutdown()
	os.Exit(code)
}

func TestNonExisting(t *testing.T) {
	_, err := cert.LoadFromDisk("nonexistant", "nonexistant")

	if err == nil {
		t.Error("loading from disk did not return an error for unknown files")
	}
}

func TestCerts(t *testing.T) {
	for _, tc := range testCerts {
		_, err := cert.LoadFromDisk(tc.FullChainFile, tc.KeyFile)

		if err == nil {
			t.Error("switching inputs did not return an error")
		}

		c, err := cert.LoadFromDisk(tc.KeyFile, tc.FullChainFile)
		_ = c

		if err != nil {
			t.Error(err)
		}

		if tc.SubjectCN != c.SubjectCN() {
			t.Errorf("SubjectCN does not match, got %s expected %s", c.SubjectCN(), tc.SubjectCN)
		}

		if tc.Fingerprint != c.Fingerprint() {
			t.Errorf("Fingerprint does not match, got %s expected %s", c.Fingerprint(), tc.Fingerprint)
		}

		if tc.NotAfter != c.NotAfter() {
			t.Errorf("Expire date does not match, got %s expected %s", c.NotAfter(), tc.NotAfter)
		}

		if !bytes.Equal(tc.CertData, c.LeafPEM()) {
			t.Error("LeafCert() does not return correct cert data")
		}

		if !bytes.Equal(tc.ChainData, c.ChainPEM()) {
			t.Error("ChainPEM() does not return correct cert data")
		}

		if !bytes.Equal(tc.FullChainKeyData, c.FullPEM()) {
			t.Error("FullPEM() does not return correct cert data")
		}

	}
}
