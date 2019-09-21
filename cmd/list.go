package cmd

import (
	"fmt"
	"github.com/emgag/keyvault-certdeploy/internal/lib/vault"
	"github.com/go-playground/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"strings"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List certificates in vault",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := vault.NewClient(viper.GetString("keyvault.url"))

		if err != nil {
			log.Fatal(err)
		}

		log.Notice("Fetching certificate list")

		certs, err := client.GetCertificates()

		if err != nil {
			log.Fatalf("Error fetching certificate list: %s", err)
		}

		fmt.Print("SubjectCN\tKeyalgo\tExpire\tAlt names\n")

		for _, c := range certs {
			fmt.Printf(
				"%s\t%s\t%s\t%s\n",
				c.SubjectCN(),
				c.PublicKeyAlgorithm(),
				c.NotAfter(),
				strings.Join(c.Leaf.DNSNames, ","),
			)
		}
	},
}
