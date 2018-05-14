package cmd

import (
	"log"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "keyvault-certdeploy",
	Short: "TBD",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default locations are $HOME/.config/keyvault-certdeploy.yml, /etc/keyvault-certdeploy.yml, $PWD/keyvault-certdeploy.yml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName("keyvault-certdeploy")

	// set defaults for redis
	//viper.SetDefault("xxx", "xxx")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()

		if err != nil {
			viper.AddConfigPath(path.Join(home, ".config"))
		}

		viper.AddConfigPath("/etc")
		viper.AddConfigPath(".")
	}

	viper.SetEnvPrefix("kvcd")
	viper.AutomaticEnv()

	// if a config file is found, read it in.
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Could not open config file.")
	}
}
