package fcw

import (
	"fmt"
	cert9util "github.com/eyedeekay/cert9util/lib"
	"path/filepath"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// CertManager handles Firefox certificate operations
type CertManager struct {
	profileDir string
	db         *cert9util.CertificateDB9
}

// NewCertManager creates a certificate manager for the given Firefox profile
func NewCertManager(profileDir string) (*CertManager, error) {
	certDB := filepath.Join(profileDir, "cert9.db")
	db, err := cert9util.NewCertificateDB9(certDB)
	if err != nil {
		return nil, fmt.Errorf("failed to open cert9.db: %w", err)
	}

	return &CertManager{
		profileDir: profileDir,
		db:         db,
	}, nil
}

// AddCertificate imports a certificate from a PEM file
func (cm *CertManager) AddCertificate(certPath, nickname string) error {
	cert, err := LoadCertificate(certPath)
	if err != nil {
		return fmt.Errorf("failed to load certificate: %w", err)
	}
	return cm.db.AddCertificate(cert, nickname, cert9util.TrustAttributes{})
}

// LoadCertificate loads a certificate from a PEM file
func LoadCertificate(certPath string) (*x509.Certificate, error) {
	pemData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM data")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert, nil
}

// RemoveCertificate removes a certificate by its nickname
func (cm *CertManager) RemoveCertificate(nickname string) error {
	return cm.db.RemoveCertificate(nickname)
}

// ListCertificates returns all certificates in the database
func (cm *CertManager) ListCertificates() ([]cert9util.CertificateInfo, error) {
	return cm.db.ListCertificates()
}
