package cmd

import (
	"github.com/emgag/keyvault-certdeploy/internal/lib/vault"
	"github.com/go-playground/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"io/ioutil"
	"os"
	"path"
)

func init() {
	dumpCmd.Flags().StringP(
		"dir",
		"d",
		".",
		"Directory to save files to",
	)

	dumpCmd.Flags().String(
		"key",
		"privkey.pem",
		"Name of the private key file",
	)

	dumpCmd.Flags().String(
		"cert",
		"cert.pem",
		"Name of the leaf certificate file",
	)

	dumpCmd.Flags().String(
		"chain",
		"chain.pem",
		"Name of the certificate chain file",
	)

	dumpCmd.Flags().String(
		"fullchain",
		"fullchain.pem",
		"Name of the full certificate chain file",
	)

	dumpCmd.Flags().String(
		"fullchainprivkey",
		"fullchain.privkey.pem",
		"Name of the full certificate chain + private key file",
	)

	rootCmd.AddCommand(dumpCmd)
}

var dumpCmd = &cobra.Command{
	Use:   "dump <subject> <keyalgo>",
	Short: "Dump certificate and key from vault to current directory or dir, if supplied",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := vault.NewClient(viper.GetString("keyvault.url"))

		if err != nil {
			log.Fatal(err)
		}

		subject := args[0]
		keyalgo := args[1]

		c, err := client.PullCertificate(subject, keyalgo)

		if err != nil {
			log.Fatal(err)
		}

		out, _ := cmd.Flags().GetString("dir")

		files := []struct {
			Param       string
			Data        []byte
			FileMode    os.FileMode
			Description string
		}{
			{"key", c.RawKey, os.FileMode(0400), "private key"},
			{"cert", c.LeafPEM(), os.FileMode(0444), "certificate"},
			{"chain", c.ChainPEM(), os.FileMode(0444), "certificate chain"},
			{"fullchain", c.RawCert, os.FileMode(0444), "full certificate chain"},
			{"fullchainprivkey", c.FullPEM(), os.FileMode(0400), "full certificate chain + private key"},
		}

		for _, f := range files {
			p, _ := cmd.Flags().GetString(f.Param)
			name := path.Join(out, p)

			if _, err := os.Stat(name); err == nil {
				log.Noticef("%s already exists, removing", name)

				if err := os.Remove(name); err != nil {
					log.Warnf("Error removing file %s", name)
				}
			}

			if err := ioutil.WriteFile(name, f.Data, f.FileMode); err != nil {
				log.Alertf("Error writing %s: %s", f.Description, err)
			} else {
				log.Infof("Wrote %s to %s", f.Description, name)
			}
		}

	},
}
