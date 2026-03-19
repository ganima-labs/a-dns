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
	var email string

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Configure les identifiants API OVH",
		Long: `Configuration interactive du CLI a-dns.

Deux méthodes d'authentification sont disponibles:

1. OAuth2 Service Account (RECOMMANDÉ)
   - Plus sécurisé et moderne
   - Pas d'expiration des tokens
   - Gestion fine des permissions via IAM

2. Application Keys (traditionnel)
   - Méthode historique
   - Consumer Key avec durée limitée
   - Plus complexe à configurer

Exemples:
  # Configuration OAuth2 interactive
  ./a-dns setup -m 2

  # Configuration OAuth2 en une ligne
  ./a-dns setup -m 2 -i EU.xxxx -x secret -z domaine.com -e email@example.com

  # Configuration Application Keys
  ./a-dns setup -k APP_KEY -s APP_SECRET -c CONSUMER_KEY -z domaine.com -e email@example.com
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("╔══════════════════════════════════════════════════════════════╗")
			fmt.Println("║           Configuration a-dns - CLI DNS OVH                   ║")
			fmt.Println("╚══════════════════════════════════════════════════════════════╝")
			fmt.Println("")

			if authMethod == "" {
				fmt.Println("Méthode d'authentification:")
				fmt.Println("")
				fmt.Println("  [1] Application Keys (traditionnel)")
				fmt.Println("      - Créer une app sur https://eu.api.ovh.com/createApp")
				fmt.Println("      - Générer un Consumer Key avec droits DNS")
				fmt.Println("")
				fmt.Println("  [2] OAuth2 Service Account (RECOMMANDÉ)")
				fmt.Println("      - Manager OVH → Gestion des comptes → Service Accounts")
				fmt.Println("      - Permissions IAM automatiques")
				fmt.Println("")
				fmt.Print("Choix [2]: ")
				fmt.Scanln(&authMethod)
				if authMethod == "" {
					authMethod = "2"
				}
			}

			if endpoint == "" {
				fmt.Println("")
				fmt.Println("Endpoint OVH:")
				fmt.Println("  - ovh-eu: Europe (défaut)")
				fmt.Println("  - ovh-ca: Canada")
				fmt.Println("  - ovh-us: USA")
				fmt.Print("Endpoint [ovh-eu]: ")
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
				fmt.Println("║              OAuth2 Service Account                          ║")
				fmt.Println("╚══════════════════════════════════════════════════════════════╝")
				fmt.Println("")
				fmt.Println("ÉTAPE 1: Créer un Service Account")
				fmt.Println("  1. Allez sur https://www.ovh.com/manager/dedicated")
				fmt.Println("  2. Cliquez sur votre nom → Gestion des comptes")
				fmt.Println("  3. Onglet 'Comptes de service' → Ajouter")
				fmt.Println("  4. Nom: a-dns, Description: CLI DNS")
				fmt.Println("")
				fmt.Println("ÉTAPE 2: Configurer les permissions IAM")
				fmt.Println("  Dans le compte de service → Politique IAM → Ajouter:")
				fmt.Println(`  {
    "rules": [{
      "effect": "allow",
      "action": "*",
      "resource": "urn:v1:eu:resource:dnsZone:*"
    }]
  }`)
				fmt.Println("")
				fmt.Println("ÉTAPE 3: Récupérer les credentials")
				fmt.Println("  ⚠️  Notez le Client ID et Client Secret (affichés une seule fois)")
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
				fmt.Println("╔══════════════════════════════════════════════════════════════╗")
				fmt.Println("║              Application Keys (traditionnel)                 ║")
				fmt.Println("╚══════════════════════════════════════════════════════════════╝")
				fmt.Println("")
				fmt.Println("ÉTAPE 1: Créer une application")
				fmt.Println("  Allez sur https://eu.api.ovh.com/createApp")
				fmt.Println("  Notez l'Application Key et Application Secret")
				fmt.Println("")
				fmt.Println("ÉTAPE 2: Générer un Consumer Key")
				fmt.Println("  Exécutez cette commande:")
				fmt.Println("")
				fmt.Printf("  curl -X POST https://eu.api.ovh.com/1.0/auth/credential \\\n")
				fmt.Printf("    -H 'Content-Type: application/json' \\\n")
				fmt.Printf("    -H 'X-Ovh-Application: VOTRE_APP_KEY' \\\n")
				fmt.Printf("    -d '{\"accessRules\":[")
				fmt.Printf("{\"method\":\"GET\",\"path\":\"/domain/zone/*\"},")
				fmt.Printf("{\"method\":\"POST\",\"path\":\"/domain/zone/*\"},")
				fmt.Printf("{\"method\":\"PUT\",\"path\":\"/domain/zone/*\"},")
				fmt.Printf("{\"method\":\"DELETE\",\"path\":\"/domain/zone/*\"}")
				fmt.Printf("]}'\n")
				fmt.Println("")
				fmt.Println("ÉTAPE 3: Valider le Consumer Key")
				fmt.Println("  Ouvrez l'URL 'validationUrl' retournée")
				fmt.Println("  Connectez-vous et validez les droits")
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

			if email == "" {
				fmt.Println("")
				fmt.Println("Email pour Let's Encrypt (notifications expiration):")
				fmt.Print("Email: ")
				fmt.Scanln(&email)
				cfg.Email = email
			}

			if defaultZone == "" {
				fmt.Println("")
				fmt.Println("Zone DNS par défaut (optionnel):")
				fmt.Print("Zone [votre-domaine.com]: ")
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
			fmt.Println("║              Configuration terminée !                        ║")
			fmt.Println("╚══════════════════════════════════════════════════════════════╝")
			fmt.Println("")
			fmt.Println("Fichier créé: ~/.a-dns.yaml")
			fmt.Println("")
			fmt.Println("Testez la connexion:")
			fmt.Println("  ./a-dns list-zones")
			fmt.Println("")
			fmt.Println("Commandes disponibles:")
			fmt.Println("  ./a-dns list-records <zone>")
			fmt.Println("  ./a-dns add-record <zone> <type> <valeur>")
			fmt.Println("  ./a-dns cert request <domain>")
			fmt.Println("")

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
	cmd.Flags().StringVarP(&defaultZone, "default-zone", "z", "", "Zone par défaut")
	cmd.Flags().StringVarP(&email, "email", "", "", "Email pour Let's Encrypt")

	return cmd
}
