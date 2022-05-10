package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// "github.com/aws/aws-lambda-go/lambda"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pulumi/pulumi-linode/sdk/v3/go/linode"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CreateVMOptions struct {
	RequesterUsername string      `json:"requesterUsername"`
	VmOptions         VmOptions   `json:"vmOptions"`
	PulumiDetails     PulumiModel `json:"pulumiDetails"`
}

type PulumiModel struct {
	Username     string `json:"username"`
	ProjectName  string `json:"projectName"`
	InstanceName string `json:"instanceName"`
}

type VmOptions struct {
	ImageName       string `json:"imageName"`
	OperatingSystem string `json:"operatingSystem"`
	LabelName       string `json:"labelName"`
	PrivateIp       bool   `json:"privateIp"`
	RegionName      string `json:"regionName"`
	Password        string `json:"password"`
	VmType          string `json:"type"`
	SwapSize        int16  `json:"swapSize"`
}

func HandleRequest() (map[string]interface{}, error) {
	ctx := context.Background()

	// These are the data points that would be required as input for the lambda to work
	vmImageName := "linode/ubuntu18.04"
	vmOperatingSystem := "ubuntu18.04"
	vmLabelName := "simple_instance"
	vmPrivateIp := true
	vmRegionName := "us-central"
	vmPassword := "terr4form-test"
	vmSwapSize := 256
	vmType := "g6-standard-1"

	// projectname and username would be neccessary
	pulumiUserName := "ganesh"
	pulumiStackName := fmt.Sprintf("%s-%s-%s", pulumiUserName, vmOperatingSystem, time.Now().Format("01-02-2006"))
	pulumiProjectName := "slack-bot"

	requesterUsername := "carlos"
	instanceName := "web"

	deployFunc := func(ctx *pulumi.Context) error {
		_, err := linode.NewInstance(ctx, instanceName, &linode.InstanceArgs{
			Image:     pulumi.String(vmImageName),
			Label:     pulumi.String(vmLabelName),
			PrivateIp: pulumi.Bool(vmPrivateIp),
			Region:    pulumi.String(vmRegionName),
			RootPass:  pulumi.String(vmPassword),
			SwapSize:  pulumi.Int(vmSwapSize),
			Tags: pulumi.StringArray{
				pulumi.String(pulumiProjectName),
				pulumi.String(requesterUsername),
			},
			Type: pulumi.String(vmType),
		})

		if err != nil {
			return err
		}

		return nil
	}

	is_destroy := len(os.Args) > 1 && os.Args[1] == "destroy"

	if is_destroy {

		// create or select a stack matching the specified name and project.
		// this will set up a workspace with everything necessary to run our inline program (deployFunc)
		s, err := auto.SelectStackInlineSource(ctx, pulumiStackName, pulumiProjectName, deployFunc)

		if err != nil {
			log.Println(err)
		}

		log.Printf("Delete instruction")
		res, err := s.Destroy(ctx)

		if err != nil {
			log.Println("Error Destroying the VM instance")
		} else {
			log.Println("Destroying Project Successfully done")
			log.Println(res)
		}

	} else {

		// create or select a stack matching the specified name and project.
		// this will set up a workspace with everything necessary to run our inline program (deployFunc)
		s, err := auto.NewStackInlineSource(ctx, pulumiStackName, pulumiProjectName, deployFunc)

		if err != nil {
			log.Println(err)
		}

		log.Printf("Created/Selected stack %q\n", pulumiStackName)

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

		log.Printf("Update succeeded! with result %s", res.Summary.Message)
	}

	return map[string]interface{}{
		"message": "something else",
	}, nil
}

func main() {
	app_env, app_env_exists := os.LookupEnv("APP_ENV")
	if app_env_exists && app_env == "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}
	// lambda.Start(HandleRequest)
	// HandleRequest()

	r := gin.New()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "home path"})
	})

	r.POST("/create-vm", func(ctx *gin.Context) {
		log.Println(os.Getenv("PATH"))
		// HandleRequest()
		ctx.JSON(200, gin.H{"message": "create-vm"})
	})

	r.POST("/delete-vm", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "delete-vm"})
	})

	r.Run()
}
