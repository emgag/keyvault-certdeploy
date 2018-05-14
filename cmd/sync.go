package cmd

import (
	"fmt"
	"github.com/emgag/keyvault-certdeploy/internal/lib/cert"
	"github.com/emgag/keyvault-certdeploy/internal/lib/config"
	"github.com/emgag/keyvault-certdeploy/internal/lib/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync configured certificates from vault to system",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		certs := []config.CertList{}
		err := viper.UnmarshalKey("certs", &certs)

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		hooks := make(map[string]string)
		_ = hooks

		for _, c := range certs {
			rc, err := vault.PullCertificate(c.SubjectCN, c.KeyAlgo)

			if err != nil {
				fmt.Printf("%+v\n", err)
				continue
			}

			lc, err := cert.LoadFromDisk(c.PrivKey, c.FullChain)

			if err == nil && rc.Fingerprint() == lc.Fingerprint() {
				fmt.Printf("%s: Certificate already up to date\n", rc.SubjectCN())
				continue
			}

			if err := ioutil.WriteFile(c.PrivKey, rc.RawKey, os.FileMode(0400)); err != nil {
				fmt.Printf("%s: Error writing private key file %s: %v\n", rc.SubjectCN(), c.PrivKey, err)
				continue
			}

			if err := ioutil.WriteFile(c.FullChain, rc.RawCert, os.FileMode(0444)); err != nil {
				fmt.Printf("%s: Error writing fullchain file %s: %v\n", rc.SubjectCN(), c.FullChain, err)
				continue
			}
		}

	},
}
