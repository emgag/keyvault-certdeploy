package cmd

import (
	"fmt"
	"github.com/emgag/keyvault-certdeploy/internal/lib/vault"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
)

func init() {
	rootCmd.AddCommand(dumpCmd)
}

var dumpCmd = &cobra.Command{
	Use:   "dump <subject> <keyalgo> [dir]",
	Short: "Dump certificate from vault to current directory or dir, if supplied",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		subject := args[0]
		keyalgo := args[1]
		out := "."

		if len(args) > 2 {
			out = args[2]
		}

		c, err := vault.PullCertificate(subject, keyalgo)

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		keyFile := path.Join(out, "privkey.pem")

		if err := ioutil.WriteFile(keyFile, c.RawKey, os.FileMode(0400)); err != nil {
			fmt.Printf("Error writing private key: %v\n", err)
			os.Exit(1)
		}

		chainFile := path.Join(out, "fullchain.pem")

		if err := ioutil.WriteFile(chainFile, c.RawCert, os.FileMode(0444)); err != nil {
			fmt.Printf("Error writing certificate chain: %v\n", err)
			os.Exit(1)
		}
	},
}
