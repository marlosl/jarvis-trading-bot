package command

import (
	"jarvis-trading-bot/cmd/cli/helpers"

	"github.com/spf13/cobra"
)

var (
	updateDnsCmd = &cobra.Command{
		Use:   "update-dns",
		Short: "Update DNS record in Cloudflare.",
		Run: func(cmd *cobra.Command, args []string) {
			helpers.UpdateDNS()
		},
	}
)

func init() {
	rootCmd.AddCommand(updateDnsCmd)
}
