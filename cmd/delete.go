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
	deleteCmd.Flags().BoolP(
		"yes",
		"y",
		false,
		"Don't confirm before deleting",
	)

	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete <subject> <keyalgo>",
	Short: "Deletes certificate from vault",
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

		if y, _ := cmd.Flags().GetBool("yes"); !y {
			prompt := promptui.Prompt{
				Label: fmt.Sprintf(
					"Really delete %s certificate %s (Expire %s)",
					c.PublicKeyAlgorithm(),
					c.SubjectCN(),
					c.NotAfter().Format(time.RFC3339),
				),
				IsConfirm: true,
			}

			_, err = prompt.Run()

			if err != nil {
				log.Info("Not removing certificate")
				return
			}
		}

		err = client.DeleteCertificate(subject, keyalgo)

		if err != nil {
			log.Fatal(err)
		}

		log.Warnf("%s cert %s removed", c.PublicKeyAlgorithm(), c.SubjectCN())
	},
}
