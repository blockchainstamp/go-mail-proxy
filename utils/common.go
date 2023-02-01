package utils

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

func ReadJsonFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

func WriteJsonFile(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, data, 0666); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func FileExists(fileName string) (os.FileInfo, bool) {

	fileInfo, err := os.Lstat(fileName)

	if fileInfo != nil || (err != nil && !os.IsNotExist(err)) {
		return fileInfo, true
	}
	return nil, false
}
func publicKey(priv any) any {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func GenerateByParam(certPath, keyPath string, days int, ecdsaCurve, host string) error {

	if len(host) == 0 {
		host = "localhost"
	}

	var priv any
	var err error
	switch ecdsaCurve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, 2048)

	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		err = fmt.Errorf("unrecognized elliptic curve: %q", ecdsaCurve)
	}
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	keyUsage := x509.KeyUsageDigitalSignature
	if _, isRSA := priv.(*rsa.PrivateKey); isRSA {
		keyUsage |= x509.KeyUsageKeyEncipherment
	}

	var notBefore = time.Now()

	notAfter := notBefore.Add(time.Duration(days) * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to open cert.pem for writing: %v", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to write data to cert.pem: %v", err)
	}
	if err := certOut.Close(); err != nil {
		return fmt.Errorf("error closing cert.pem: %v", err)
	}

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open key.pem for writing: %v", err)

	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("unable to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("failed to write data to key.pem: %v", err)
	}
	if err := keyOut.Close(); err != nil {
		return fmt.Errorf("error closing key.pem: %v", err)
	}
	return nil
}
