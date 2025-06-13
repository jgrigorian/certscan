package show

import (
	"context"
	"fmt"

	"github.com/jgrigorian/certscan/common"

	"github.com/urfave/cli/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Certificate(c *cli.Context) error {
	secretName := c.String("secret")
	namespace := c.String("namespace")

	if secretName == "" {
		return cli.Exit("Secret name is required", 1)
	}

	client, err := common.InitClient()
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}

	var s *v1.Secret
	if namespace == "" {
		s, err = client.CoreV1().Secrets("default").Get(context.TODO(), secretName, metav1.GetOptions{})
	} else {
		s, err = client.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	}
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	if s.Type != "kubernetes.io/tls" {
		return cli.Exit("Secret is not a TLS secret", 1)
	}

	certData, exists := s.Data["tls.crt"]
	if !exists {
		return cli.Exit("No certificate data found in secret", 1)
	}

	cert, err := common.GetCertInfo(certData)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	fmt.Printf("Certificate Information for %s/%s:\n", s.Namespace, s.Name)
	fmt.Printf("Subject: %s\n", cert.Subject.CommonName)
	fmt.Printf("Issuer: %s\n", cert.Issuer.CommonName)
	fmt.Printf("Valid from: %s\n", cert.NotBefore.Format("2006-01-02 15:04:05"))
	fmt.Printf("Valid until: %s\n", cert.NotAfter.Format("2006-01-02 15:04:05"))
	fmt.Printf("Serial Number: %s\n", cert.SerialNumber.String())
	fmt.Printf("DNS Names: %v\n", cert.DNSNames)
	fmt.Printf("IP Addresses: %v\n", cert.IPAddresses)

	return nil
}
