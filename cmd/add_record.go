package cmd

import (
	"a-dns/internal/config"
	"a-dns/internal/ovhclient"

	"github.com/spf13/cobra"
)

func NewAddRecordCmd() *cobra.Command {
	var zone string
	var subDomain string
	var recordType string
	var target string
	var ttl int

	cmd := &cobra.Command{
		Use:   "add-record [zone] [type] [target]",
		Short: "Ajoute un enregistrement DNS",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			zone = args[0]
			recordType = args[1]
			target = args[2]

			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			if zone == "" {
				zone = cfg.DefaultZone
			}

			client, err := cfg.GetOVHClient()
			if err != nil {
				return err
			}

			record, err := ovhclient.AddRecord(client, zone, subDomain, recordType, target, ttl)
			if err != nil {
				return err
			}

			return outputResult(record)
		},
	}

	cmd.Flags().StringVarP(&zone, "zone", "z", "", "zone DNS")
	cmd.Flags().StringVarP(&subDomain, "subdomain", "s", "", "sous-domaine")
	cmd.Flags().IntVarP(&ttl, "ttl", "t", 0, "TTL (par défaut: automatique)")

	return cmd
}
