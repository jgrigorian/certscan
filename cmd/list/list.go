package list

import (
	"context"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jgrigorian/certscan/common"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

type certInfo struct {
	name          string
	namespace     string
	expiration    time.Time
	issuer        string
	daysRemaining float64
	err           error
}

func processCertificates(secrets []v1.Secret, expiring bool, results chan<- certInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, s := range secrets {
		if s.Type != "kubernetes.io/tls" {
			continue
		}

		certData, exists := s.Data["tls.crt"]
		if !exists {
			results <- certInfo{err: fmt.Errorf("no certificate data found in secret %s/%s", s.Namespace, s.Name)}
			continue
		}

		cert, err := common.GetCertInfo(certData)
		if err != nil {
			results <- certInfo{err: fmt.Errorf("failed to parse certificate in %s/%s: %w", s.Namespace, s.Name, err)}
			continue
		}

		daysRemaining := math.Round(time.Until(cert.NotAfter).Hours() / 24)
		if expiring && daysRemaining > 10 {
			continue
		}

		results <- certInfo{
			name:          s.Name,
			namespace:     s.Namespace,
			expiration:    cert.NotAfter,
			issuer:        strings.Join(cert.Issuer.Organization, ""),
			daysRemaining: daysRemaining,
		}
	}
}

func Certificates(c *cli.Context) error {
	namespace := c.String("namespace")
	allNamespaces := c.Bool("all-namespaces")
	expiring := c.Bool("expiring")

	client, err := common.InitClient()
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}

	var secrets *v1.SecretList
	if allNamespaces {
		secrets, err = client.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	} else {
		secrets, err = client.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	}
	if err != nil {
		return fmt.Errorf("failed to list secrets: %w", err)
	}

	// Create a channel to collect results
	results := make(chan certInfo, len(secrets.Items))
	var wg sync.WaitGroup

	// Process certificates concurrently
	numWorkers := 4 // Adjust based on your system's capabilities
	secretsPerWorker := len(secrets.Items) / numWorkers
	if secretsPerWorker == 0 {
		secretsPerWorker = 1
	}

	for i := 0; i < numWorkers; i++ {
		start := i * secretsPerWorker
		end := start + secretsPerWorker
		if i == numWorkers-1 {
			end = len(secrets.Items)
		}

		wg.Add(1)
		go processCertificates(secrets.Items[start:end], expiring, results, &wg)
	}

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Create table
	t := table.New().Border(lipgloss.HiddenBorder()).StyleFunc(func(row, col int) lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 2)
	})

	t.Headers(
		color.HiMagentaString("Secret Name"),
		color.HiMagentaString("Namespace"),
		color.HiMagentaString("Expiration"),
		color.HiMagentaString("Issuer"),
		color.HiMagentaString("Days Remaining"),
	)

	certCount := 0
	var errors []string

	// Process results
	for result := range results {
		if result.err != nil {
			errors = append(errors, result.err.Error())
			continue
		}

		certCount++
		daysStr := fmt.Sprintf("%v", result.daysRemaining)
		if result.daysRemaining <= 10 {
			daysStr = color.RedString("%v", result.daysRemaining)
		}

		t.Row(
			result.name,
			result.namespace,
			result.expiration.Format("2006-01-02"),
			result.issuer,
			daysStr,
		)
	}

	if certCount == 0 {
		if len(errors) > 0 {
			fmt.Fprintf(os.Stderr, "Errors encountered:\n%s\n", strings.Join(errors, "\n"))
		}
		color.Red("No certificates found in namespace \"%s\"", namespace)
		return cli.Exit("", 1)
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Warnings:\n%s\n", strings.Join(errors, "\n"))
	}

	fmt.Println(t.Render())
	return nil
}
