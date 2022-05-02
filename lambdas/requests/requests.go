package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/pulumi/pulumi-linode/sdk/v3/go/linode"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CreateVMOptions struct {
	requesterUsername string      `json:"requesterUsername"`
	vmOptions         VmOptions   `json:"vmOptions"`
	pulumiDetails     PulumiModel `json:pulumiDetails`
}

type PulumiModel struct {
	username     string `json:"username"`
	projectName  string `json:"projectName"`
	instanceName string `json:"instanceName"`
}

type VmOptions struct {
	imageName       string `json:"imageName"`
	operatingSystem string `json:"operatingSystem"`
	labelName       string `json:"labelName"`
	privateIp       bool   `json:"privateIp"`
	regionName      string `json:"regionName"`
	password        string `json:"password"`
	vmType          string `json:"type"`
	swapSize        int16  `json:"swapSize"`
}

func HandleRequest(ctx context.Context, request CreateVMOptions) (map[string]interface{}, error) {

	pulumiStackName := fmt.Sprintf("%s-%s-%s", request.pulumiDetails.username, request.vmOptions.operatingSystem, time.Now().Format("01-02-2006"))

	deployFunc := func(ctx *pulumi.Context) error {
		_, err := linode.NewInstance(ctx, request.pulumiDetails.instanceName, &linode.InstanceArgs{
			Image:     pulumi.String(request.vmOptions.imageName),
			Label:     pulumi.String(request.vmOptions.labelName),
			PrivateIp: pulumi.Bool(request.vmOptions.privateIp),
			Region:    pulumi.String(request.vmOptions.regionName),
			RootPass:  pulumi.String(request.vmOptions.password),
			SwapSize:  pulumi.Int(request.vmOptions.swapSize),
			Tags: pulumi.StringArray{
				pulumi.String("slack-bot"),
				pulumi.String(request.pulumiDetails.username),
			},
			Type: pulumi.String(request.vmOptions.vmType),
		})

		if err != nil {
			return err
		}

		return nil
	}

	s, err := auto.NewStackInlineSource(ctx, pulumiStackName, request.pulumiDetails.projectName, deployFunc)

	if err != nil {
		log.Println(err)
	}

	log.Printf("Created/Selected stack %q\n", s)

	w := s.Workspace()

	log.Println("Installing the Linode plugin")

	// for inline source programs, we must manage plugins ourselves
	err = w.InstallPlugin(ctx, "linode", "v3.7.1")
	if err != nil {
		log.Printf("Failed to install program plugins: %v\n", err)
		os.Exit(1)
	}

	log.Println("Successfully installed Linode plugin")

	// set stack configuration specifying the AWS region to deploy
	linode_token := os.Getenv("LINODE_TOKEN")
	s.SetConfig(ctx, "linode:token", auto.ConfigValue{Value: linode_token, Secret: true})

	log.Println("Successfully set config")
	log.Println("Starting refresh")

	_, err = s.Refresh(ctx)
	if err != nil {
		log.Printf("Failed to refresh stack: %v\n", err)
		os.Exit(1)
	}

	log.Println("Refresh succeeded!")

	log.Println("Starting update")

	// wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// run the update to deploy our Linode instance
	res, err := s.Up(ctx, stdoutStreamer)

	if err != nil {
		log.Printf("Failed to update stack: %v\n\n", err)
		os.Exit(1)
	}

	log.Println("Update succeeded! with result ")
	log.Println(res)

	return map[string]interface{}{
			"statusCode": 200,
			"headers":    map[string]string{"Content-Type": "application/json"},
			"body":       request,
		},
		nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	lambda.Start(HandleRequest)
}
