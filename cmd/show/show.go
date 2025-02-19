package show

import (
	"context"
	"fmt"
	"github.com/jgrigorian/certscan/common"
	"github.com/fatih/color"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/urfave/cli/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math"
	"time"
)

func Certificate(c *cli.Context) {
	secretName := c.String("secret")
	namespace := c.String("namespace")

	client := common.InitClient()
	var secrets *v1.SecretList

	if namespace == "" {
		secrets, _ = client.CoreV1().Secrets("default").List(context.TODO(), metav1.ListOptions{})
	} else {
		secrets, _ = client.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	}

	// Table
	t := table.New().Border(lipgloss.HiddenBorder()).StyleFunc(func(row, col int) lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 2)
	})

	for _, s := range secrets.Items {
		if s.Type == "kubernetes.io/tls" {
			certInfo := common.GetCertInfo(s.Data["tls.crt"])
			timeUntil := time.Until(certInfo.NotAfter)
			daysRemaining := math.Round(timeUntil.Hours() / 24)

			if s.Name == secretName {
				t.Headers(color.HiMagentaString(s.Name))
				t.BorderRow(false)

				t.Row("Namespace", color.WhiteString(s.Namespace))
				t.Row("Valid From", color.WhiteString(certInfo.NotBefore.Format("2006-01-02")))
				t.Row("Valid Until", color.WhiteString(certInfo.NotAfter.Format("2006-01-02")))
				t.Row("Issuer", color.WhiteString(strings.Join(certInfo.Issuer.Organization[0:], "")))

				if daysRemaining <= 10 {
					t.Row("Days Remaining", color.RedString("%v", daysRemaining))
				} else if (daysRemaining > 10) && (daysRemaining < 20) {
					t.Row("Days Remaining", color.YellowString("%v", daysRemaining))
				} else {
					t.Row("Days Remaining", color.GreenString("%v", daysRemaining))
				}

				t.Row("Subject Alternate Name (DNS)",
					color.WhiteString(strings.Join(certInfo.DNSNames, "\n")),
				)
			}
		}
	}

	fmt.Println(t.Render())

}
