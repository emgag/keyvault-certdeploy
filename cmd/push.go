package cmd

import (
	"github.com/emgag/keyvault-certdeploy/internal/lib/cert"
	"github.com/emgag/keyvault-certdeploy/internal/lib/vault"
	"github.com/go-playground/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:   "push <privkey.pem> <fullchain.pem>",
	Short: "Push a certificate to the vault",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		keyFile := args[0]
		certFile := args[1]

		c, err := cert.LoadFromDisk(keyFile, certFile)

		if err != nil {
			log.Fatalf("Error loading cert: %s", err)
		}

		err = vault.PushCertificate(c)

		if err != nil {
			log.Errorf("Error pushing cert: %s", err)
		} else {
			log.Noticef(
				"Successfully pushed %s certificate %s (%v, %s) to the vault",
				c.PublicKeyAlgorithm(),
				c.SubjectCN(),
				c.NotAfter(),
				c.Fingerprint(),
			)
		}
	},
}
