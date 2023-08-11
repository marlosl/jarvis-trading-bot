package command

import (
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/utils"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "A CLI utility to manage the Jarvis Trading Bot API",
}

func Execute() error {
	cobra.OnInitialize(utils.InitConfig)

	os.Setenv(consts.AwsRegion, "us-east-1")
	os.Setenv(consts.ConfigTableName, "trading-bot-config")
	os.Setenv(consts.TransactionsTableName, "trading-bot-transactions")

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd.Execute()
}
