package cmd

import (
	"a-dns/internal/config"
	"a-dns/internal/i18n"
	"a-dns/internal/ovhclient"

	"github.com/spf13/cobra"
)

func NewUpdateRecordCmd() *cobra.Command {
	var zone string
	var subDomain string
	var target string
	var ttl int

	cmd := &cobra.Command{
		Use:   "update-record [zone] [record-id] [type] [target]",
		Short: i18n.T("dns.update_record.short"),
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			zone = args[0]
			recordID := args[1]
			recordType := args[2]
			target = args[3]

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

			record, err := ovhclient.UpdateRecord(client, zone, recordID, subDomain, recordType, target, ttl)
			if err != nil {
				return err
			}

			return outputResult(record)
		},
	}

	cmd.Flags().StringVarP(&zone, "zone", "z", "", i18n.T("dns.flag.zone"))
	cmd.Flags().StringVarP(&subDomain, "subdomain", "s", "", i18n.T("flag.subdomain"))

	return cmd
}
