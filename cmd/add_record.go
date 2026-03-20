package cmd

import (
	"a-dns/internal/config"
	"a-dns/internal/i18n"
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
		Short: i18n.T("dns.add_record.short"),
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

	cmd.Flags().StringVarP(&zone, "zone", "z", "", i18n.T("dns.flag.zone"))
	cmd.Flags().StringVarP(&subDomain, "subdomain", "s", "", i18n.T("flag.subdomain"))
	cmd.Flags().IntVarP(&ttl, "ttl", "t", 0, i18n.T("flag.ttl"))

	return cmd
}
