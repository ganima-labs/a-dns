package cmd

import (
	"a-dns/internal/config"
	"a-dns/internal/ovhclient"

	"github.com/spf13/cobra"
)

func NewListZonesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-zones",
		Short: "Liste toutes les zones DNS disponibles",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			client, err := cfg.GetOVHClient()
			if err != nil {
				return err
			}

			zones, err := ovhclient.ListZones(client)
			if err != nil {
				return err
			}

			return outputResult(zones)
		},
	}
	return cmd
}
