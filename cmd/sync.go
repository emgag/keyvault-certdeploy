package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/emgag/keyvault-certdeploy/internal/lib/cert"
	"github.com/emgag/keyvault-certdeploy/internal/lib/config"
	"github.com/emgag/keyvault-certdeploy/internal/lib/vault"
	"github.com/go-playground/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	syncCmd.Flags().Bool(
		"nohooks",
		false,
		"Disable running hooks after cert update",
	)

	syncCmd.Flags().BoolP(
		"force",
		"f",
		false,
		"Force update even if version on disk matches the one in vault",
	)

	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync configured certificates from vault to system",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var certs []config.CertList
		err := viper.UnmarshalKey("certs", &certs)

		if err != nil {
			log.Fatalf("Error loading config: %s", err)
		}

		client, err := vault.NewClient(viper.GetString("keyvault.url"))

		if err != nil {
			log.Fatal(err)
		}

		hooks := make(map[string]bool)

		for _, c := range certs {
			log.Infof("Fetching %s cert %s", c.KeyAlgo, c.SubjectCN)

			rc, err := client.PullCertificate(c.SubjectCN, c.KeyAlgo)

			if err != nil {
				log.Errorf("Error fetching %s cert %s: %s", c.KeyAlgo, c.SubjectCN, err)
				continue
			}

			lc, err := cert.LoadFromDisk(c.PrivKey, c.FullChain)

			if err == nil && rc.Fingerprint() == lc.Fingerprint() {
				f, _ := cmd.Flags().GetBool("force")

				if f {
					log.Noticef("%s cert %s: already up to date, forcing update", c.KeyAlgo, c.SubjectCN)
				} else {
					log.Noticef("%s cert %s: already up to date", c.KeyAlgo, c.SubjectCN)
					continue
				}
			}

			files := []struct {
				Name        string
				Data        []byte
				FileMode    os.FileMode
				Description string
			}{
				{c.PrivKey, rc.RawKey, os.FileMode(0400), "private key"},
				{c.Cert, rc.LeafPEM(), os.FileMode(0444), "certificate"},
				{c.Chain, rc.ChainPEM(), os.FileMode(0444), "certificate chain"},
				{c.FullChain, rc.RawCert, os.FileMode(0444), "full certificate chain"},
				{c.FullChainPrivKey, rc.FullPEM(), os.FileMode(0400), "full certificate chain + private key"},
			}

			for _, f := range files {
				if f.Name == "" {
					log.Noticef("No filename for %s defined, skipping", f.Description)
					continue
				}

				if _, err := os.Stat(f.Name); err == nil {
					log.Noticef("%s already exists, removing", f.Name)

					if err := os.Remove(f.Name); err != nil {
						log.Warnf("Error removing file %s", f.Name)
					}
				}

				if err := ioutil.WriteFile(f.Name, f.Data, f.FileMode); err != nil {
					log.Alertf("Error writing %s: %s", f.Description, err)
				} else {
					log.Infof("Wrote %s to %s", f.Description, f.Name)
				}
			}

			for _, h := range c.Hooks {
				hooks[strings.TrimSpace(h)] = true
			}
		}

		if skip, _ := cmd.Flags().GetBool("nohooks"); !skip {
			for h := range hooks {
				log.Noticef("Run hook %s", h)
				c := strings.Split(h, " ")
				command := exec.Command(c[0], c[1:]...)

				if out, err := command.CombinedOutput(); err != nil {
					log.Errorf("Error running hook %s: %s", h, out)
				}
			}
		} else {
			log.Notice("Skipping update hooks")
		}
	},
}
