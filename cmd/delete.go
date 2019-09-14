package cmd

import (
	"fmt"
	"time"

	"github.com/emgag/keyvault-certdeploy/internal/lib/vault"
	"github.com/go-playground/log"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
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
		subject := args[0]
		keyalgo := args[1]

		c, err := vault.PullCertificate(subject, keyalgo)

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
				log.Notice("Not removing certificate")
				return
			}
		}

		err = vault.DeleteCertificate(subject, keyalgo)

		if err != nil {
			log.Fatal(err)
		}

		log.Warnf("%s cert %s removed", c.PublicKeyAlgorithm(), c.SubjectCN())
	},
}
