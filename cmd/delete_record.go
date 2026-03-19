package cmd

import (
	"a-dns/internal/config"
	"a-dns/internal/ovhclient"

	"github.com/spf13/cobra"
)

func NewDeleteRecordCmd() *cobra.Command {
	var zone string

	cmd := &cobra.Command{
		Use:   "delete-record [zone] [record-id]",
		Short: "Supprime un enregistrement DNS",
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

			err = ovhclient.DeleteRecord(client, zone, recordID)
			if err != nil {
				return err
			}

			return outputResult(map[string]string{
				"status":   "deleted",
				"zone":     zone,
				"recordID": recordID,
			})
		},
	}

	cmd.Flags().StringVarP(&zone, "zone", "z", "", "zone DNS")

	return cmd
}
