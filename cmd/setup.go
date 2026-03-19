package cmd

import (
	"fmt"
	"os"

	"a-dns/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewSetupCmd() *cobra.Command {
	var endpoint string
	var authMethod string
	var appKey string
	var appSecret string
	var consumerKey string
	var oauth2ClientID string
	var oauth2ClientSecret string
	var defaultZone string

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Configure les identifiants API OVH",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("=== Configuration OVH API ===")
			fmt.Println("")
			fmt.Println("Méthodes d'authentification:")
			fmt.Println("  1) Application Keys (traditionnel)")
			fmt.Println("  2) OAuth2 Service Account (moderne, recommandé)")
			fmt.Println("")

			if authMethod == "" {
				fmt.Print("Méthode [1]: ")
				fmt.Scanln(&authMethod)
				if authMethod == "" {
					authMethod = "1"
				}
			}

			if endpoint == "" {
				fmt.Print("Endpoint [ovh-eu]: ")
				fmt.Scanln(&endpoint)
				if endpoint == "" {
					endpoint = "ovh-eu"
				}
			}

			cfg := &config.Config{
				Endpoint:    endpoint,
				DefaultZone: defaultZone,
			}

			if authMethod == "2" {
				fmt.Println("")
				fmt.Println("OAuth2 Service Account:")
				fmt.Println("  Dans OVH Manager → Gestion des comptes → Service Accounts")
				fmt.Println("  Créer un compte service avec permissions DNS (IAM)")
				fmt.Println("")

				if oauth2ClientID == "" {
					fmt.Print("Client ID: ")
					fmt.Scanln(&oauth2ClientID)
				}

				if oauth2ClientSecret == "" {
					fmt.Print("Client Secret: ")
					fmt.Scanln(&oauth2ClientSecret)
				}

				cfg.OAuth2ClientID = oauth2ClientID
				cfg.OAuth2ClientSecret = oauth2ClientSecret
			} else {
				fmt.Println("")
				fmt.Println("Application Keys:")
				fmt.Println("  1. Créez une app sur https://eu.api.ovh.com/createApp")
				fmt.Println("  2. Notez l'Application Key et le Application Secret")
				fmt.Println("  3. Créez un Consumer Key avec les droits DNS")
				fmt.Println("")

				if appKey == "" {
					fmt.Print("Application Key: ")
					fmt.Scanln(&appKey)
				}

				if appSecret == "" {
					fmt.Print("Application Secret: ")
					fmt.Scanln(&appSecret)
				}

				if consumerKey == "" {
					fmt.Print("Consumer Key: ")
					fmt.Scanln(&consumerKey)
				}

				cfg.AppKey = appKey
				cfg.AppSecret = appSecret
				cfg.ConsumerKey = consumerKey
			}

			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			viper.SetConfigFile(home + "/.a-dns.yaml")
			viper.Set("endpoint", cfg.Endpoint)
			if cfg.DefaultZone != "" {
				viper.Set("default_zone", cfg.DefaultZone)
			}

			if cfg.OAuth2ClientID != "" {
				viper.Set("oauth2_client_id", cfg.OAuth2ClientID)
				viper.Set("oauth2_client_secret", cfg.OAuth2ClientSecret)
			} else {
				viper.Set("app_key", cfg.AppKey)
				viper.Set("app_secret", cfg.AppSecret)
				viper.Set("consumer_key", cfg.ConsumerKey)
			}

			if err := viper.WriteConfig(); err != nil {
				return err
			}

			fmt.Println("")
			fmt.Println("Configuration sauvegardée dans ~/.a-dns.yaml")

			return nil
		},
	}

	cmd.Flags().StringVarP(&endpoint, "endpoint", "e", "ovh-eu", "endpoint OVH (ovh-eu, ovh-ca, etc.)")
	cmd.Flags().StringVarP(&authMethod, "method", "m", "", "méthode: 1=App Keys, 2=OAuth2")
	cmd.Flags().StringVarP(&appKey, "app-key", "k", "", "Application Key")
	cmd.Flags().StringVarP(&appSecret, "app-secret", "s", "", "Application Secret")
	cmd.Flags().StringVarP(&consumerKey, "consumer-key", "c", "", "Consumer Key")
	cmd.Flags().StringVarP(&oauth2ClientID, "oauth2-client-id", "i", "", "OAuth2 Client ID")
	cmd.Flags().StringVarP(&oauth2ClientSecret, "oauth2-client-secret", "x", "", "OAuth2 Client Secret")
	cmd.Flags().StringVarP(&defaultZone, "default-zone", "z", "ganima.xyz", "Zone par défaut")

	return cmd
}
