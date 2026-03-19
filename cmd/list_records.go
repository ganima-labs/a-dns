package cmd

import (
	"fmt"

	"a-dns/internal/config"
	"a-dns/internal/ovhclient"

	"github.com/spf13/cobra"
)

func NewListRecordsCmd() *cobra.Command {
	var zone string
	var recordType string
	var subDomain string

	cmd := &cobra.Command{
		Use:   "list-records [zone]",
		Short: "Liste les enregistrements DNS d'une zone",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				zone = args[0]
			}

			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			if zone == "" {
				zone = cfg.DefaultZone
			}
			if zone == "" {
				return fmt.Errorf("zone non spécifiée")
			}

			client, err := cfg.GetOVHClient()
			if err != nil {
				return err
			}

			records, err := ovhclient.ListRecords(client, zone, recordType, subDomain)
			if err != nil {
				return err
			}

			return outputResult(records)
		},
	}

	cmd.Flags().StringVarP(&zone, "zone", "z", "", "zone DNS (ex: ganima.xyz)")
	cmd.Flags().StringVarP(&recordType, "type", "t", "", "type d'enregistrement (A, AAAA, CNAME, MX, TXT, etc.)")
	cmd.Flags().StringVarP(&subDomain, "subdomain", "s", "", "sous-domaine à filtrer")

	return cmd
}
