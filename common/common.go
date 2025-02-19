package common

import (
	"crypto/x509"
	pem2 "encoding/pem"
	"k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var secretsClient coreV1Types.SecretInterface

func getKubeconfig() string {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
	}
	return kubeconfig
}

func InitClient() *kubernetes.Clientset {
	kubeconfig := getKubeconfig()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	secretsClient = client.CoreV1().Secrets("default")
	return client
}

func GetCertInfo(pem []byte) *x509.Certificate {
	pemFile, _ := pem2.Decode(pem)
	c, _ := x509.ParseCertificate(pemFile.Bytes)
	return c
}
