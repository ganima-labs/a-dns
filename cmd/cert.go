package cmd

import (
	"a-dns/internal/certbot"
	"a-dns/internal/config"
	"a-dns/internal/i18n"
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
		Short: i18n.T("cert.title"),
		Long:  i18n.T("cert.title_long"),
	}

	cmd.PersistentFlags().StringVar(&certDir, "cert-dir", "", i18n.T("cert.flag.cert_dir"))
	cmd.PersistentFlags().StringVar(&certEmail, "email", "", i18n.T("cert.flag.email"))
	cmd.PersistentFlags().BoolVar(&certStaging, "staging", false, i18n.T("cert.flag.staging"))
	cmd.PersistentFlags().IntVar(&certRenewDays, "renew-days", 30, i18n.T("cert.flag.renew_days"))

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
		Short: i18n.T("cert.request.short"),
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
				return fmt.Errorf("%s", i18n.T("cert.error.email_required"))
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

	cmd.Flags().StringSliceVar(&certSANs, "sans", []string{}, i18n.T("cert.flag.sans"))

	return cmd
}

func NewCertRenewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "renew <domain>",
		Short: i18n.T("cert.renew.short"),
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
				return fmt.Errorf("%s", i18n.T("cert.error.email_required"))
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
		Short: i18n.T("cert.list.short"),
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
		Short: i18n.T("cert.status.short"),
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
				if strings.Contains(err.Error(), "lecture certificat") || strings.Contains(err.Error(), "reading certificate") {
					return fmt.Errorf("%s", i18n.TWithData("cert.error.not_found", map[string]interface{}{"Domain": domain}))
				}
				return err
			}

			return outputResult(info)
		},
	}
}
