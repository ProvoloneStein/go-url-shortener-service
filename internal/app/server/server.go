package server

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	serverReadTimeout  = 15 * time.Second
	serverWriteTimeout = 30 * time.Second
	maxHeaderBytes     = 1 << 20 // 1 MB
)

func initTLS(certName, keyName string) error {
	certFile, err := os.OpenFile(certName, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer certFile.Close()

	keyFile, err := os.OpenFile(keyName, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"CustomeApp"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certWriter := bufio.NewWriter(certFile)
	_, err = certWriter.Write(certPEM.Bytes())
	if err != nil {
		return err
	}
	if err := certWriter.Flush(); err != nil {
		return fmt.Errorf("ошибка записи буфера в файл : %w", err)
	}

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	keyWriter := bufio.NewWriter(keyFile)
	_, err = keyWriter.Write(privateKeyPEM.Bytes())
	if err != nil {
		return err
	}
	if err := keyWriter.Flush(); err != nil {
		return fmt.Errorf("ошибка записи буфера в файл : %w", err)
	}
	return nil
}

func Run(enableTLS bool, addr string, handler http.Handler) error {
	httpServer := &http.Server{
		Addr:           addr,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    serverReadTimeout,
		WriteTimeout:   serverWriteTimeout,
	}
	if enableTLS {
		certName := "cert.pem"
		keyName := "key.pem"
		if err := initTLS(certName, keyName); err != nil {
			return fmt.Errorf("server initTLS:  %w", err)
		}
		if err := httpServer.ListenAndServeTLS(certName, keyName); err != nil {
			return fmt.Errorf("server run:  %w", err)
		}
	} else {
		if err := httpServer.ListenAndServe(); err != nil {
			return fmt.Errorf("server run:  %w", err)
		}
	}
	return nil
}
