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

	"go.uber.org/zap"
)

const (
	serverReadTimeout  = 15 * time.Second
	serverWriteTimeout = 30 * time.Second
	maxHeaderBytes     = 1 << 20 // 1 MB
	defaultFilePerm    = 0600
	serialNumber       = 1234
	defaultTLSPath     = "tls"
)

func initTLS(logger *zap.Logger, certName, keyName string) error {
	certFile, err := os.OpenFile(defaultTLSPath+"/"+certName, os.O_CREATE|os.O_RDWR, defaultFilePerm)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла сертификата: %w", err)
	}
	defer func() {
		if deferErr := certFile.Close(); deferErr != nil {
			logger.Error("ошибка при закрытии файла сертификата.", zap.Error(deferErr))
		}
	}()

	keyFile, err := os.OpenFile(defaultTLSPath+"/"+keyName, os.O_CREATE|os.O_RDWR, defaultFilePerm)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла ключа: %w", err)
	}
	defer func() {
		if deferErr := keyFile.Close(); deferErr != nil {
			logger.Error("ошибка при закрытии файла ключа.", zap.Error(deferErr))
		}
	}()

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(serialNumber),
		Subject: pkix.Name{
			Organization: []string{"CustomeApp"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback}, //nolint:gomnd // ip
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0), //nolint:gomnd // cert time
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096) //nolint:gomnd // GenerateKey
	if err != nil {
		return fmt.Errorf("ошибка генерации ключа: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("ошибка создания сертификата: %w", err)
	}

	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	if err != nil {
		return fmt.Errorf("ошибка записи сертификата в буфер: %w", err)
	}

	certWriter := bufio.NewWriter(certFile)
	_, err = certWriter.Write(certPEM.Bytes())
	if err != nil {
		return fmt.Errorf("ошибка записи в буфер: %w", err)
	}
	if err := certWriter.Flush(); err != nil {
		return fmt.Errorf("ошибка записи буфера в файл: %w", err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return fmt.Errorf("ошибка записи сертификата ключа в буфер: %w", err)
	}
	keyWriter := bufio.NewWriter(keyFile)
	_, err = keyWriter.Write(privateKeyPEM.Bytes())
	if err != nil {
		return fmt.Errorf("ошибка записи в буфер: %w", err)
	}
	if err := keyWriter.Flush(); err != nil {
		return fmt.Errorf("ошибка записи буфера в файл : %w", err)
	}
	return nil
}

func InitServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:           addr,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    serverReadTimeout,
		WriteTimeout:   serverWriteTimeout,
	}
}

func Run(logger *zap.Logger, enableTLS bool, serv *http.Server) error {
	if enableTLS {
		certName := "cert.pem"
		keyName := "key.pem"
		if err := initTLS(logger, certName, keyName); err != nil {
			return fmt.Errorf("server initTLS:  %w", err)
		}
		if err := serv.ListenAndServeTLS(certName, keyName); err != nil {
			return fmt.Errorf("server run:  %w", err)
		}
	} else {
		if err := serv.ListenAndServe(); err != nil {
			return fmt.Errorf("server run:  %w", err)
		}
	}
	return nil
}
