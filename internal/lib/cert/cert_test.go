package cert_test

import (
	"bytes"
	"fmt"
	"github.com/emgag/keyvault-certdeploy/internal/lib/cert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type CertInfo struct {
	KeyFile            string
	KeyData            []byte
	CertFile           string
	CertData           []byte
	ChainFile          string
	ChainData          []byte
	FullChainFile      string
	FullChainData      []byte
	SubjectCN          string
	NotAfter           time.Time
	PublicKeyAlgorithm string
	Fingerprint        string
}

var testData = "../../../testdata"

var testCerts = []*CertInfo{
	{
		SubjectCN:          "c1.example.org",
		PublicKeyAlgorithm: "rsa",
		NotAfter:           time.Date(2023, 6, 13, 15, 48, 0, 0, time.UTC),
		Fingerprint:        "d8a713b1b7f1da027c59cc636a017d94df4f0e70176975d5704e34d18c4f2803",
	}, {
		SubjectCN:          "c1.example.org",
		PublicKeyAlgorithm: "ecdsa",
		NotAfter:           time.Date(2023, 6, 13, 15, 48, 0, 0, time.UTC),
		Fingerprint:        "bf182f09bda97ae3b91c2d26c37220138891a3130489b26326d763e74f0aef81",
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
		c.ChainFile = filepath.Join(testData, "chain.pem")
		c.KeyData, _ = ioutil.ReadFile(c.KeyFile)
		c.CertData, _ = ioutil.ReadFile(c.CertFile)
		c.ChainData, _ = ioutil.ReadFile(c.ChainFile)
		c.FullChainData, _ = ioutil.ReadFile(c.FullChainFile)
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

		if bytes.Compare(tc.CertData, c.LeafPEM()) != 0 {
			t.Error("LeafCert() does not return correct cert data")
		}

		if bytes.Compare(tc.ChainData, c.ChainPEM()) != 0 {
			t.Error("ChainPEM() does not return correct cert data")
		}

		//if bytes.Compare(tc.FullChainData, c.FullPEM()) != 0 {
		//	t.Error("FullPEM() does not return correct cert data")
		//}

	}
}
