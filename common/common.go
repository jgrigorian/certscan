package common

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	secretsClient coreV1Types.SecretInterface
	certCache     = struct {
		sync.RWMutex
		certs map[string]*x509.Certificate
	}{
		certs: make(map[string]*x509.Certificate),
	}
)

func getKubeconfig() string {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
	}
	return kubeconfig
}

func InitClient() (*kubernetes.Clientset, error) {
	kubeconfig := getKubeconfig()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config: %w", err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	secretsClient = client.CoreV1().Secrets("default")
	return client, nil
}

func GetCertInfo(pemData []byte) (*x509.Certificate, error) {
	if len(pemData) == 0 {
		return nil, errors.New("empty certificate data")
	}

	// Generate a cache key from the PEM data
	cacheKey := string(pemData)

	// Check cache first
	certCache.RLock()
	if cert, exists := certCache.certs[cacheKey]; exists {
		certCache.RUnlock()
		return cert, nil
	}
	certCache.RUnlock()

	// Parse certificate
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Cache the result
	certCache.Lock()
	certCache.certs[cacheKey] = cert
	certCache.Unlock()

	return cert, nil
}

// ClearCertCache clears the certificate cache
func ClearCertCache() {
	certCache.Lock()
	certCache.certs = make(map[string]*x509.Certificate)
	certCache.Unlock()
}
