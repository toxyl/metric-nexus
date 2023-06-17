package metrics

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	reNonASCII = regexp.MustCompile(`[^a-zA-Z0-9\-\_]+`)
)

func interfaceToFloat64(i interface{}) (float64, bool) {
	v := float64(0)
	switch c := i.(type) {
	case []byte:
		v, _ = strconv.ParseFloat(string(c), 64)
	case string:
		v, _ = strconv.ParseFloat(c, 64)
	case int:
		v = float64(c)
	case int8:
		v = float64(c)
	case int16:
		v = float64(c)
	case int32:
		v = float64(c)
	case int64:
		v = float64(c)
	case uint:
		v = float64(c)
	case uint8:
		v = float64(c)
	case uint16:
		v = float64(c)
	case uint32:
		v = float64(c)
	case uint64:
		v = float64(c)
	case float32:
		v = float64(c)
	case float64:
		v = c
	default:
		return v, false
	}
	return v, true
}

func fileExists(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	_, err = file.Stat()
	return err == nil
}

func sanitizeKey(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, " ", "_")
	name = reNonASCII.ReplaceAllString(name, "")
	return strings.TrimSpace(name)
}

func generateSelfSignedCertificate(commonName, organization, keyFile, certFile string) (pathKey string, pathCert string, err error) {
	if keyFile == "" {
		keyFile = filepath.Join(os.TempDir(), "nexus-tls.key")
	}

	if certFile == "" {
		certFile = filepath.Join(os.TempDir(), "nexus-tls.cert")
	}

	if fileExists(keyFile) && fileExists(certFile) {
		return keyFile, certFile, errors.New("files already exist")
	}
	if fileExists(keyFile) {
		return keyFile, certFile, errors.New("key file already exists but cert file is missing")
	}
	if fileExists(certFile) {
		return keyFile, certFile, errors.New("cert file already exists but key file is missing")
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return keyFile, certFile, fmt.Errorf("generating private key failed: %s", err.Error())
	}

	keyBuf := new(bytes.Buffer)
	err = pem.Encode(keyBuf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	if err != nil {
		return keyFile, certFile, fmt.Errorf("encoding private key failed: %s", err.Error())
	}
	_ = os.WriteFile(keyFile, keyBuf.Bytes(), 0600)

	tml := x509.Certificate{
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(5, 0, 0),
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{organization},
		},
		BasicConstraintsValid: true,
	}
	cert, err := x509.CreateCertificate(rand.Reader, &tml, &tml, &key.PublicKey, key)
	if err != nil {
		return keyFile, certFile, fmt.Errorf("creating certificate failed: %s", err.Error())
	}

	certBuf := new(bytes.Buffer)
	err = pem.Encode(certBuf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	if err != nil {
		return keyFile, certFile, fmt.Errorf("encoding certificate failed: %s", err.Error())
	}

	_ = os.WriteFile(certFile, certBuf.Bytes(), 0600)

	return keyFile, certFile, nil
}
