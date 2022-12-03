package cmd

import (
	"fmt"
	"time"

	"github.com/emgag/keyvault-certdeploy/internal/lib/vault"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	pruneCmd.Flags().IntP(
		"days",
		"d",
		7,
		"Delete certificates after this many days of being expired",
	)

	pruneCmd.Flags().BoolP(
		"noop",
		"n",
		false,
		"Just list expired certificates, don't actually remove the certs",
	)

	pruneCmd.Flags().BoolP(
		"yes",
		"y",
		false,
		"Don't confirm before pruning",
	)

	rootCmd.AddCommand(pruneCmd)
}

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove expired certificates from vault",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := vault.NewClient(viper.GetString("keyvault.url"))

		if err != nil {
			log.Fatal(err)
		}

		log.Info("Fetching certificate list")

		certs, err := client.GetCertificates()

		if err != nil {
			log.Fatalf("Error fetching certificate list: %s", err)
		}

		days, err := cmd.Flags().GetInt("days")

		if err != nil {
			log.Fatal(err)
		}

		// filter certificates
		i := 0

		for _, c := range certs {
			diff := time.Since(c.NotAfter())

			if (diff.Hours()/float64(24) - float64(days)) > 0 {
				certs[i] = c
				i++
			}
		}

		certs = certs[:i]

		if len(certs) == 0 {
			log.Warn("No expired certificates found")
			return
		}

		if n, _ := cmd.Flags().GetBool("noop"); n {
			fmt.Println("Following certificates are expired and would be removed:")

			for _, c := range certs {
				fmt.Printf(
					"* %s (%s, %s)\n",
					c.SubjectCN(),
					c.PublicKeyAlgorithm(),
					c.NotAfter().Format(time.RFC3339),
				)
			}

		} else {
			if y, _ := cmd.Flags().GetBool("yes"); !y {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf("Really delete %d expired certificates", len(certs)),
					IsConfirm: true,
				}

				_, err = prompt.Run()

				if err != nil {
					log.Info("Not removing any certificates")
					return
				}
			}

			for _, c := range certs {
				err = client.DeleteCertificate(c.SubjectCN(), c.PublicKeyAlgorithm())

				if err != nil {
					log.Errorf("Failed deleting certificate %s (%s, %s): %s",
						c.SubjectCN(),
						c.PublicKeyAlgorithm(),
						c.NotAfter().Format(time.RFC3339),
						err,
					)
				} else {
					log.Warnf(
						"Certificate %s (%s, %s) deleted",
						c.SubjectCN(),
						c.PublicKeyAlgorithm(),
						c.NotAfter().Format(time.RFC3339),
					)
				}
			}
		}

	},
}
