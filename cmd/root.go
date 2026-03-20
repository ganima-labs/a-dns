package cmd

import (
	"fmt"
	"os"

	"a-dns/internal/i18n"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var outputFormat string
var language string

func init() {
	if err := i18n.Init("fr"); err != nil {
		fmt.Fprintf(os.Stderr, "i18n init error: %v\n", err)
	}
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "a-dns",
		Short: i18n.T("app.short"),
		Long:  i18n.T("app.long"),
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", i18n.T("flag.config"))
	cmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", i18n.T("flag.output"))
	cmd.PersistentFlags().StringVarP(&language, "lang", "l", "fr", "langue: fr, en, es")

	cmd.AddCommand(
		NewListZonesCmd(),
		NewListRecordsCmd(),
		NewAddRecordCmd(),
		NewDeleteRecordCmd(),
		NewUpdateRecordCmd(),
		NewSetupCmd(),
		NewCertCmd(),
		NewSkillCmd(),
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
		if lang := viper.GetString("language"); lang != "" && language == "fr" {
			language = lang
			i18n.SetLanguage(language)
		}
	}
}
