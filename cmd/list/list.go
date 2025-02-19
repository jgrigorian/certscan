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
	"time"
)

func Certificates(c *cli.Context) {
	namespace := c.String("namespace")
	allNamespaces := c.Bool("all-namespaces")
	expiring := c.Bool("expiring")

	client := common.InitClient()
	var secrets *v1.SecretList
	certRequest := make(chan *v1.SecretList)
	certCount := 0

	if allNamespaces == true {
		secrets, _ = client.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	} else {
		secrets, _ = client.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	}

	//secrets, _ = client.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	go func() { certRequest <- secrets }()
	certResult := <-certRequest

	// Table
	t := table.New().Border(lipgloss.HiddenBorder()).StyleFunc(func(row, col int) lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 2)
	})

	t.Headers(color.HiMagentaString("Secret Name"), color.HiMagentaString("Namespace"), color.HiMagentaString("Expiration"), color.HiMagentaString("Issuer"), color.HiMagentaString("Days Remaining"))

	for _, s := range certResult.Items {
		if expiring == true && s.Type == "kubernetes.io/tls" {
			certCount += 1
			certExp := common.GetCertInfo(s.Data["tls.crt"]).NotAfter
			certIssuer := common.GetCertInfo(s.Data["tls.crt"]).Issuer.Organization[0:]
			daysRemaining := math.Round(time.Until(certExp).Hours() / 24)
			if daysRemaining <= 10 {
				t.Row(s.Name, s.Namespace, certExp.Format("2006-01-02"), strings.Join(certIssuer, ""), color.RedString("%v", daysRemaining))
			}
		} else if s.Type == "kubernetes.io/tls" {
			certCount += 1
			certExp := common.GetCertInfo(s.Data["tls.crt"]).NotAfter
			certIssuer := common.GetCertInfo(s.Data["tls.crt"]).Issuer.Organization[0:]
			daysRemaining := math.Round(time.Until(certExp).Hours() / 24)
			if daysRemaining <= 10 {
				t.Row(s.Name, s.Namespace, certExp.Format("2006-01-02"), strings.Join(certIssuer, ""), color.RedString("%v", daysRemaining))
				//t.Row(s.Name, s.Namespace, certExp.Format("2006-01-02"), color.RedString("%v", daysRemaining))
			} else {
				t.Row(s.Name, s.Namespace, certExp.Format("2006-01-02"), strings.Join(certIssuer, ""), fmt.Sprintf("%v", daysRemaining))
				//t.Row(s.Name, s.Namespace, certExp.Format("2006-01-02"), fmt.Sprintf("%v", daysRemaining))
			}
		}
	}
	if certCount == 0 {
		color.Red("No certificates found in namespace \"%s\"", namespace)
		os.Exit(1)
	}
	fmt.Println(t.Render())
}
