package cmd

import (
	"a-dns/internal/certbot"
	"a-dns/internal/config"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var certDir string
var certEmail string
var certStaging bool
var certRenewDays int
var certSANs []string

func getCertbot(cfg *config.Config) *certbot.Certbot {
	if cfg.OAuth2ClientID != "" && cfg.OAuth2ClientSecret != "" {
		return certbot.NewOAuth2(cfg.Endpoint, cfg.OAuth2ClientID, cfg.OAuth2ClientSecret)
	}
	return certbot.New(cfg.Endpoint, cfg.AppKey, cfg.AppSecret, cfg.ConsumerKey)
}

func NewCertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cert",
		Short: "Gestion des certificats Let's Encrypt",
		Long:  `Commandes pour gérer les certificats SSL/TLS via Let's Encrypt avec challenge DNS OVH.`,
	}

	cmd.PersistentFlags().StringVar(&certDir, "cert-dir", "", "répertoire de stockage des certificats (défaut: ~/.a-dns/certs)")
	cmd.PersistentFlags().StringVar(&certEmail, "email", "", "email pour Let's Encrypt")
	cmd.PersistentFlags().BoolVar(&certStaging, "staging", false, "utiliser l'environnement de test Let's Encrypt")
	cmd.PersistentFlags().IntVar(&certRenewDays, "renew-days", 30, "renouveler si moins de X jours restants")

	cmd.AddCommand(
		NewCertRequestCmd(),
		NewCertRenewCmd(),
		NewCertListCmd(),
		NewCertStatusCmd(),
	)

	return cmd
}

func NewCertRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request <domain>",
		Short: "Demander un nouveau certificat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]

			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			email := certEmail
			if email == "" {
				email = cfg.Email
			}
			if email == "" {
				return fmt.Errorf("email requis (via --email ou configuration)")
			}

			dir := certDir
			if dir == "" {
				home, _ := os.UserHomeDir()
				dir = home + "/.a-dns/certs"
			}

			cb := getCertbot(cfg)

			info, err := cb.RequestCertificate(certbot.CertRequest{
				Domain:    domain,
				SANs:      certSANs,
				Email:     email,
				CertDir:   dir,
				Staging:   certStaging,
				RenewDays: certRenewDays,
			})
			if err != nil {
				return err
			}

			return outputResult(info)
		},
	}

	cmd.Flags().StringSliceVar(&certSANs, "sans", []string{}, "Subject Alternative Names (séparés par des virgules)")

	return cmd
}

func NewCertRenewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "renew <domain>",
		Short: "Renouveler un certificat existant",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]

			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			email := certEmail
			if email == "" {
				email = cfg.Email
			}
			if email == "" {
				return fmt.Errorf("email requis (via --email ou configuration)")
			}

			dir := certDir
			if dir == "" {
				home, _ := os.UserHomeDir()
				dir = home + "/.a-dns/certs"
			}

			cb := getCertbot(cfg)

			info, err := cb.RenewCertificate(certbot.CertRequest{
				Domain:    domain,
				Email:     email,
				CertDir:   dir,
				Staging:   certStaging,
				RenewDays: certRenewDays,
			})
			if err != nil {
				return err
			}

			return outputResult(info)
		},
	}
}

func NewCertListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lister les certificats locaux",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			dir := certDir
			if dir == "" {
				home, _ := os.UserHomeDir()
				dir = home + "/.a-dns/certs"
			}

			cb := getCertbot(cfg)

			certs, err := cb.ListCertificates(dir)
			if err != nil {
				return err
			}

			return outputResult(certs)
		},
	}
}

func NewCertStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status <domain>",
		Short: "Afficher le statut d'un certificat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]

			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			dir := certDir
			if dir == "" {
				home, _ := os.UserHomeDir()
				dir = home + "/.a-dns/certs"
			}

			cb := getCertbot(cfg)

			info, err := cb.GetCertInfo(domain, dir, certRenewDays)
			if err != nil {
				if strings.Contains(err.Error(), "lecture certificat") {
					return fmt.Errorf("aucun certificat trouvé pour %s", domain)
				}
				return err
			}

			return outputResult(info)
		},
	}
}
