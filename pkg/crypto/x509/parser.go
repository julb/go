package x509

import (
	"bytes"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/julb/go/pkg/util/date"
)

type TlsCertificateMetadata struct {
	x509Certificate *x509.Certificate
	IssuerCN        string                 `yaml:"issuerCn" json:"issuerCn"`
	SubjectCN       string                 `yaml:"subjectCn" json:"subjectCn"`
	Sans            []string               `yaml:"sans" json:"sans"`
	SerialNumber    string                 `yaml:"serialNumber" json:"serialNumber"`
	Sha1Fingerprint string                 `yaml:"sha1Fingerprint" json:"sha1Fingerprint"`
	Validity        TlsCertificateValidity `yaml:"validity" json:"validity"`
}

type TlsCertificateValidity struct {
	From          string `yaml:"from" json:"from"`
	To            string `yaml:"to" json:"to"`
	RemainingDays int64  `yaml:"remainingDays" json:"remainingDays"`
	Expired       bool   `yaml:"expired" json:"expired"`
	Valid         bool   `yaml:"valid" json:"valid"`
}

func ParsePemFileAndGetMetadata(path string) (*TlsCertificateMetadata, error) {
	// read certificate bytes
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParsePemContentAndGetMetadata(bytes)
}

func ParsePemContentAndGetMetadata(bytes []byte) (*TlsCertificateMetadata, error) {
	// use pem decoder
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("fail to decode pem file")
	}

	// parse pem to get certificate info
	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	// build metadata
	return &TlsCertificateMetadata{
		x509Certificate: x509Cert,
		IssuerCN:        x509Cert.Issuer.CommonName,
		SubjectCN:       x509Cert.Subject.CommonName,
		Sans:            x509Cert.DNSNames,
		SerialNumber:    hexSerialNumber(x509Cert),
		Sha1Fingerprint: sha1HexFingerprint(x509Cert),
		Validity: TlsCertificateValidity{
			From:          date.DateTimeToString(x509Cert.NotBefore),
			To:            date.DateTimeToString(x509Cert.NotAfter),
			RemainingDays: int64(time.Until(x509Cert.NotAfter).Hours() / 24),
			Valid:         x509Cert.NotBefore.Before(time.Now()) && x509Cert.NotAfter.After(time.Now()),
			Expired:       x509Cert.NotAfter.Before(time.Now()),
		},
	}, nil
}

func hexSerialNumber(x509Cert *x509.Certificate) string {
	return fmt.Sprintf("%x", x509Cert.SerialNumber)
}

func sha1HexFingerprint(x509Cert *x509.Certificate) string {
	fingerprint := sha1.Sum(x509Cert.Raw)

	var buf bytes.Buffer
	for _, f := range fingerprint {
		fmt.Fprintf(&buf, "%02x", f)
	}
	return buf.String()
}
