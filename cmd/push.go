package cmd

import (
	"fmt"
	"github.com/emgag/keyvault-certdeploy/internal/lib/cert"
	"github.com/emgag/keyvault-certdeploy/internal/lib/vault"
	"github.com/spf13/cobra"
	"os"
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
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		err = vault.PushCertificate(c)

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Printf(
			"Sucessfully pushed certificate %s (%v, %s) to the vault\n",
			c.SubjectCN(),
			c.NotAfter(),
			c.Fingerprint(),
		)
	},
}
