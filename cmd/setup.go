package cmd

import (
	"fmt"
	"os"

	"a-dns/internal/config"
	"a-dns/internal/i18n"

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
	var email string

	cmd := &cobra.Command{
		Use:   "setup",
		Short: i18n.T("setup.title"),
		Long: i18n.T("setup.title") + `

` + i18n.T("setup.auth_method") + `

1. ` + i18n.T("setup.method.oauth2") + `
   - ` + i18n.T("setup.method.oauth2_desc") + `
   - ` + i18n.T("setup.method.oauth2_iam") + `

2. ` + i18n.T("setup.method.app_keys") + `
   - ` + i18n.T("setup.method.app_keys_desc") + `
   - ` + i18n.T("setup.method.app_keys_ck") + `

Examples:
  # OAuth2 interactive
  ./a-dns setup -m 2

  # OAuth2 one-liner
  ./a-dns setup -m 2 -i EU.xxxx -x secret -z domain.com -e email@example.com

  # Application Keys
  ./a-dns setup -k APP_KEY -s APP_SECRET -c CONSUMER_KEY -z domain.com -e email@example.com
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("╔══════════════════════════════════════════════════════════════╗")
			fmt.Println("║         " + i18n.T("setup.title") + "                          ║")
			fmt.Println("╚══════════════════════════════════════════════════════════════╝")
			fmt.Println("")

			if authMethod == "" {
				fmt.Println(i18n.T("setup.auth_method"))
				fmt.Println("")
				fmt.Println("  [1] " + i18n.T("setup.method.app_keys"))
				fmt.Println("      - " + i18n.T("setup.method.app_keys_desc"))
				fmt.Println("      - " + i18n.T("setup.method.app_keys_ck"))
				fmt.Println("")
				fmt.Println("  [2] " + i18n.T("setup.method.oauth2"))
				fmt.Println("      - " + i18n.T("setup.method.oauth2_desc"))
				fmt.Println("      - " + i18n.T("setup.method.oauth2_iam"))
				fmt.Println("")
				fmt.Print(i18n.T("setup.choice_default") + ": ")
				fmt.Scanln(&authMethod)
				if authMethod == "" {
					authMethod = "2"
				}
			}

			if endpoint == "" {
				fmt.Println("")
				fmt.Println(i18n.T("setup.endpoint.title"))
				fmt.Println("  - ovh-eu: " + i18n.T("setup.endpoint.eu"))
				fmt.Println("  - ovh-ca: " + i18n.T("setup.endpoint.ca"))
				fmt.Println("  - ovh-us: " + i18n.T("setup.endpoint.us"))
				fmt.Print(i18n.T("setup.endpoint.prompt") + ": ")
				fmt.Scanln(&endpoint)
				if endpoint == "" {
					endpoint = "ovh-eu"
				}
			}

			cfg := &config.Config{
				Endpoint:    endpoint,
				DefaultZone: defaultZone,
				Email:       email,
			}

			if authMethod == "2" {
				fmt.Println("")
				fmt.Println("╔══════════════════════════════════════════════════════════════╗")
				fmt.Println("║              " + i18n.T("setup.oauth2.title") + "                          ║")
				fmt.Println("╚══════════════════════════════════════════════════════════════╝")
				fmt.Println("")
				fmt.Println(i18n.T("setup.oauth2.step1"))
				fmt.Println("  " + i18n.T("setup.oauth2.step1_1"))
				fmt.Println("  " + i18n.T("setup.oauth2.step1_2"))
				fmt.Println("  " + i18n.T("setup.oauth2.step1_3"))
				fmt.Println("  " + i18n.T("setup.oauth2.step1_4"))
				fmt.Println("")
				fmt.Println(i18n.T("setup.oauth2.step2"))
				fmt.Println("  " + i18n.T("setup.oauth2.step2_desc"))
				fmt.Println(`  {
    "rules": [{
      "effect": "allow",
      "action": "*",
      "resource": "urn:v1:eu:resource:dnsZone:*"
    }]
  }`)
				fmt.Println("")
				fmt.Println(i18n.T("setup.oauth2.step3"))
				fmt.Println("  ⚠️  " + i18n.T("setup.oauth2.step3_warning"))
				fmt.Println("")

				if oauth2ClientID == "" {
					fmt.Print(i18n.T("setup.oauth2.client_id") + ": ")
					fmt.Scanln(&oauth2ClientID)
				}

				if oauth2ClientSecret == "" {
					fmt.Print(i18n.T("setup.oauth2.client_secret") + ": ")
					fmt.Scanln(&oauth2ClientSecret)
				}

				cfg.OAuth2ClientID = oauth2ClientID
				cfg.OAuth2ClientSecret = oauth2ClientSecret
			} else {
				fmt.Println("")
				fmt.Println("╔══════════════════════════════════════════════════════════════╗")
				fmt.Println("║              " + i18n.T("setup.app_keys.title") + "                          ║")
				fmt.Println("╚══════════════════════════════════════════════════════════════╝")
				fmt.Println("")
				fmt.Println(i18n.T("setup.app_keys.step1"))
				fmt.Println("  " + i18n.T("setup.app_keys.step1_desc"))
				fmt.Println("  " + i18n.T("setup.app_keys.step1_note"))
				fmt.Println("")
				fmt.Println(i18n.T("setup.app_keys.step2"))
				fmt.Println("  " + i18n.T("setup.app_keys.step2_desc"))
				fmt.Println("")
				fmt.Printf("  curl -X POST https://eu.api.ovh.com/1.0/auth/credential \\\n")
				fmt.Printf("    -H 'Content-Type: application/json' \\\n")
				fmt.Printf("    -H 'X-Ovh-Application: YOUR_APP_KEY' \\\n")
				fmt.Printf("    -d '{\"accessRules\":[")
				fmt.Printf("{\"method\":\"GET\",\"path\":\"/domain/zone/*\"},")
				fmt.Printf("{\"method\":\"POST\",\"path\":\"/domain/zone/*\"},")
				fmt.Printf("{\"method\":\"PUT\",\"path\":\"/domain/zone/*\"},")
				fmt.Printf("{\"method\":\"DELETE\",\"path\":\"/domain/zone/*\"}")
				fmt.Printf("]}'\n")
				fmt.Println("")
				fmt.Println(i18n.T("setup.app_keys.step3"))
				fmt.Println("  " + i18n.T("setup.app_keys.step3_desc"))
				fmt.Println("  " + i18n.T("setup.app_keys.step3_login"))
				fmt.Println("")

				if appKey == "" {
					fmt.Print(i18n.T("setup.app_keys.app_key") + ": ")
					fmt.Scanln(&appKey)
				}

				if appSecret == "" {
					fmt.Print(i18n.T("setup.app_keys.app_secret") + ": ")
					fmt.Scanln(&appSecret)
				}

				if consumerKey == "" {
					fmt.Print(i18n.T("setup.app_keys.consumer_key") + ": ")
					fmt.Scanln(&consumerKey)
				}

				cfg.AppKey = appKey
				cfg.AppSecret = appSecret
				cfg.ConsumerKey = consumerKey
			}

			if email == "" {
				fmt.Println("")
				fmt.Println(i18n.T("setup.email.title"))
				fmt.Print(i18n.T("setup.email.prompt") + ": ")
				fmt.Scanln(&email)
				cfg.Email = email
			}

			if defaultZone == "" {
				fmt.Println("")
				fmt.Println(i18n.T("setup.zone.title"))
				fmt.Print(i18n.T("setup.zone.prompt") + ": ")
				fmt.Scanln(&defaultZone)
				if defaultZone != "" {
					cfg.DefaultZone = defaultZone
				}
			}

			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			configPath := home + "/.a-dns.yaml"
			viper.SetConfigFile(configPath)
			viper.Set("endpoint", cfg.Endpoint)
			if cfg.DefaultZone != "" {
				viper.Set("default_zone", cfg.DefaultZone)
			}
			if cfg.Email != "" {
				viper.Set("email", cfg.Email)
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
				if os.IsNotExist(err) {
					if err := viper.SafeWriteConfig(); err != nil {
						return err
					}
				} else {
					return err
				}
			}

			fmt.Println("")
			fmt.Println("╔══════════════════════════════════════════════════════════════╗")
			fmt.Println("║              " + i18n.T("setup.complete.title") + "                        ║")
			fmt.Println("╚══════════════════════════════════════════════════════════════╝")
			fmt.Println("")
			fmt.Println(i18n.T("setup.complete.file"))
			fmt.Println("")
			fmt.Println(i18n.T("setup.complete.test"))
			fmt.Println("  " + i18n.T("setup.complete.test_cmd"))
			fmt.Println("")
			fmt.Println(i18n.T("setup.complete.available"))
			fmt.Println("  " + i18n.T("setup.complete.cmd_records"))
			fmt.Println("  " + i18n.T("setup.complete.cmd_add"))
			fmt.Println("  " + i18n.T("setup.complete.cmd_cert"))
			fmt.Println("")

			return nil
		},
	}

	cmd.Flags().StringVarP(&endpoint, "endpoint", "e", "ovh-eu", "endpoint OVH (ovh-eu, ovh-ca, etc.)")
	cmd.Flags().StringVarP(&authMethod, "method", "m", "", "méthode: 1=App Keys, 2=OAuth2")
	cmd.Flags().StringVarP(&appKey, "app-key", "k", "", i18n.T("setup.app_keys.app_key"))
	cmd.Flags().StringVarP(&appSecret, "app-secret", "s", "", i18n.T("setup.app_keys.app_secret"))
	cmd.Flags().StringVarP(&consumerKey, "consumer-key", "c", "", i18n.T("setup.app_keys.consumer_key"))
	cmd.Flags().StringVarP(&oauth2ClientID, "oauth2-client-id", "i", "", i18n.T("setup.oauth2.client_id"))
	cmd.Flags().StringVarP(&oauth2ClientSecret, "oauth2-client-secret", "x", "", i18n.T("setup.oauth2.client_secret"))
	cmd.Flags().StringVarP(&defaultZone, "default-zone", "z", "", i18n.T("setup.zone.title"))
	cmd.Flags().StringVarP(&email, "email", "", "", i18n.T("cert.flag.email"))

	return cmd
}
