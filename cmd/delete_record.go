package cmd

import (
	"a-dns/internal/config"
	"a-dns/internal/i18n"
	"a-dns/internal/ovhclient"

	"github.com/spf13/cobra"
)

func NewDeleteRecordCmd() *cobra.Command {
	var zone string

	cmd := &cobra.Command{
		Use:   "delete-record [zone] [record-id]",
		Short: i18n.T("dns.delete_record.short"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
		zone = args[0]
		recordID := args[1]

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

		return ovhclient.DeleteRecord(client, zone, recordID)
	},
}

cmd.Flags().StringVarP(&zone, "zone", "z", "", i18n.T("dns.flag.zone"))

return cmd
}
