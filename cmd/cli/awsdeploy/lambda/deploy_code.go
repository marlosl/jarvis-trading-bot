package lambda

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var sess *session.Session

func UploadFunction(zipFile *string, functionName *string) error {
	var err error

	if sess == nil {
		sess, err = session.NewSession()
		if err != nil {
			fmt.Printf("Got an error creating a new session: %v\n", err)
			return err
		}
	}

	svc := lambda.New(sess)

	contents, err := os.ReadFile(*zipFile)
	if err != nil {
		fmt.Printf("Got error trying to read "+*zipFile+"\n%v\n", err)
		return err
	}

	updateArgs := &lambda.UpdateFunctionCodeInput{
		FunctionName: functionName,
		Publish:      aws.Bool(true),
		ZipFile:      contents,
	}

	_, err = svc.UpdateFunctionCode(updateArgs)
	if err != nil {
		fmt.Printf("Cannot update function: %v\n", err)
		return err
	}

	fmt.Printf("Successfully updated function: %s\n", *functionName)
	return nil
}
