package cmd

import (
	"errors"
	"fmt"
	"path"

	"github.com/emgag/keyvault-certdeploy/internal/lib/version"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	logQuiet   bool
	logVerbose bool
)

var rootCmd = &cobra.Command{
	Use:     "keyvault-certdeploy",
	Short:   "X.509-Certificate deployment helper for Azure Key Vault",
	Version: fmt.Sprintf("%s -- %s", version.Version, version.Commit),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if logQuiet && logVerbose {
			return errors.New("quiet and verbose are mutually exclusive")
		}

		if viper.GetString("keyvault.name") == "" || viper.GetString("keyvault.url") == "" {
			return errors.New("Invalid config: at least keyvault.name and keyvault.url need to be set")
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	_ = rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initLogger, initConfig)

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		"",
		"Config file (default locations are $HOME/.config/keyvault-certdeploy.yml, /etc/keyvault-certdeploy/keyvault-certdeploy.yml, $PWD/keyvault-certdeploy.yml)",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&logVerbose,
		"verbose",
		"v",
		false,
		"Be more verbose",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&logQuiet,
		"quiet",
		"q",
		false,
		"Be quiet",
	)
}

// initLogger sets loglevels based on flags
func initLogger() {
	level := log.WarnLevel

	if logVerbose {
		level = log.InfoLevel
	} else if logQuiet {
		level = log.ErrorLevel
	}

	log.SetLevel(level)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName("keyvault-certdeploy")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()

		if err == nil {
			viper.AddConfigPath(path.Join(home, ".config"))
		}

		viper.AddConfigPath("/etc/keyvault-certdeploy")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()

	if err != nil {
		log.Errorf("Could not open config file: %s", err)
	}
}
