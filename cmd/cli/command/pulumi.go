package command

import (
	"jarvis-trading-bot/cmd/cli/awsdeploy"

	"github.com/spf13/cobra"
)

var (
	awsCmd = &cobra.Command{
		Use:   "aws",
		Short: "Manage AWS infrastructure",
	}

	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy the Jarvis Trading Bot infrastructure.",
		Run: func(cmd *cobra.Command, args []string) {
			awsdeploy.ExecuteCommand(awsdeploy.AWS_DEPLOY_COMMAND)
		},
	}

	destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the Jarvis Trading Bot infrastructure.",
		Run: func(cmd *cobra.Command, args []string) {
			awsdeploy.ExecuteCommand(awsdeploy.AWS_DESTROY_COMMAND)
		},
	}
)

func init() {
	awsCmd.AddCommand(deployCmd)
	awsCmd.AddCommand(destroyCmd)

	rootCmd.AddCommand(awsCmd)
}
