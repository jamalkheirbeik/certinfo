package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"
)

type Certificate struct {
	NotBefore    time.Time `json:"not_before"`
	NotAfter     time.Time `json:"not_after"`
	Issuer       string    `json:"issuer"`
	Subject      string    `json:"subject"`
	IsCA         bool      `json:"is_ca"`
	SerialNumber string    `json:"serial_number"`
	Fingerprint  string    `json:"fingerprint_sha256"`
}

func newCertificate(c x509.Certificate) *Certificate {
	hash := sha256.Sum256(c.Raw)
	return &Certificate{
		NotBefore:    c.NotBefore,
		NotAfter:     c.NotAfter,
		Issuer:       c.Issuer.String(),
		Subject:      c.Subject.String(),
		IsCA:         c.IsCA,
		SerialNumber: c.SerialNumber.String(),
		Fingerprint:  hex.EncodeToString(hash[:]),
	}
}

func parseCertificateChain(pemData []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for block, rest := pem.Decode(pemData); block != nil; block, rest = pem.Decode(rest) {
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			certs = append(certs, cert)
		}
	}
	if len(certs) == 0 {
		return nil, errors.New("could not find certificates to parse")
	}
	return certs, nil
}

func verifyCertificateChain(cs []*x509.Certificate, skipRoot bool) error {
	if len(cs) == 0 {
		return errors.New("empty certificate chain")
	}
	for i := 0; i < len(cs)-1; i++ {
		child := cs[i]
		parent := cs[i+1]
		// Basic Constraints (CA Rules) (only certificates with IsCA=true can sign other certificates)
		if !parent.IsCA {
			return errors.New("invalid certificate chain. parent is not a certified authority.")
		}
		// Signature Integrity (Cryptographic Validity) (Leaf  → signed by Intermediate, Intermediate → signed by Root, Root → self-signed)
		err := child.CheckSignatureFrom(parent)
		if err != nil {
			return err
		}
		// Issuer / Subject Matching (child.Issuer == parent.Subject)
		if !reflect.DeepEqual(child.Issuer, parent.Subject) {
			return errors.New("Issuer/Subject Mismatch. Child's issuer != Parent's Subject.")
		}
	}
	// Valid Time Window (NotBefore <= Now <= NotAfter)
	currentTime := time.Now()
	for i, c := range cs {
		if currentTime.Before(c.NotBefore) || currentTime.After(c.NotAfter) {
			return fmt.Errorf("certificate %d is expired. invalid chain.\n", i)
		}
	}
	// Trust Anchor (the root must be trusted. we check if the root certificate exists in the system trust store or not)
	if !skipRoot {
		root := cs[len(cs)-1]
		roots, err := x509.SystemCertPool()
		if err != nil {
			return err
		}
		opts := x509.VerifyOptions{Roots: roots}
		_, err = root.Verify(opts)
		if err != nil {
			return fmt.Errorf("Root certificate is not trusted by the system: %s\n", err)
		}
	}

	return nil
}

func printCertificates(cs []*x509.Certificate) {
	newCs := make([]*Certificate, len(cs))
	for i, c := range cs {
		newCs[i] = newCertificate(*c)
	}
	b, err := json.Marshal(newCs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}
