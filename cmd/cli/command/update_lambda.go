package command

import (
	"fmt"
	"strings"

	"jarvis-trading-bot/cmd/cli/awsdeploy/lambda"
	"jarvis-trading-bot/cmd/cli/helpers"

	"github.com/spf13/cobra"
)

var (
	handlers = helpers.GetFunctionNames()

	cmdLambda = &cobra.Command{
		Use:       "update-lambda",
		Short:     "Update Lambda function code in the API",
		Long:      "Update Lambda function code in the API.\nValid options are: " + strings.Join(handlers, ", ") + ". Use no options will update all functions.",
		ValidArgs: handlers,
		Args:      cobra.OnlyValidArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				for _, handler := range handlers {
					zipFile, functionName := helpers.GetZipFilenameAndFullFunctionName(handler)
					lambda.UploadFunction(&zipFile, &functionName)
				}
				return
			}

			zipFile, functionName := helpers.GetZipFilenameAndFullFunctionName(args[0])
			fmt.Printf("zipFile: %s, functionName: %s\n", zipFile, functionName)
			lambda.UploadFunction(&zipFile, &functionName)
		},
	}
)

func init() {
	rootCmd.AddCommand(cmdLambda)
}
