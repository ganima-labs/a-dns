package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var outputFormat string

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "a-dns",
		Short: "CLI DNS OVH pour agent IA",
		Long:  `CLI d'administration des zones DNS OVH, optimisé pour une utilisation par agent IA avec sortie JSON structurée.`,
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "fichier de configuration (par défaut ~/.a-dns.yaml)")
	cmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "format de sortie: json, yaml, table")

	cmd.AddCommand(
		NewListZonesCmd(),
		NewListRecordsCmd(),
		NewAddRecordCmd(),
		NewDeleteRecordCmd(),
		NewUpdateRecordCmd(),
		NewSetupCmd(),
		NewCertCmd(),
	)

	return cmd
}

var rootCmd = NewRootCmd()

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".a-dns")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// Config loaded successfully
	}
}
