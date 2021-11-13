package main

import (
	"flag"

	"github.com/golang-test-task/container"
)

var (
	image            *string // Docker Image name
	bashCmd          *string // Cmd for the container
	cloudwatchGroup  *string // Name of the Cloudwatch Log Group
	cloudwatchStream *string // Name of the Cloudwatch Log Stream
	awsAccessKey     *string // AWS Access Key ID
	awsSecretKey     *string // AWS Secret Access Key
	awsRegion        *string // AWS Region
)

func init() {
	image = flag.String("docker-image", "python:3.7", "Name of the docker image")
	bashCmd = flag.String("bash-command", `$'pip install pip -U && pip
	install tqdm && python -c \"import time\ncounter = 0\nwhile True:\n\tprint(counter)\n\tcounter
	= counter + 1\n\ttime.sleep(0.1)"' `, "Bash Command to run in the container")
	cloudwatchGroup = flag.String("cloudwatch-group", "golang-test-task-group-1", "Name of the cloudwatch Group")
	cloudwatchStream = flag.String("cloudwatch-stream", "golang-test-task-group-2", "Name of the Cloudwatch Stream")
	awsAccessKey = flag.String("aws-access-key-id", "XXXXXXXXXXXXXXXXXX", "AWS Access Key ID")
	awsSecretKey = flag.String("aws-secret-access-key", "XXXXXxxXXXXxxxXXX", "AWS Secret Access Key")
	awsRegion = flag.String("aws-region", "ap-southeast-1", "AWS Region")
	flag.Parse()
}

func main() {
	container.CreateContainer(*image, *bashCmd, *cloudwatchGroup, *cloudwatchStream, *awsAccessKey, *awsSecretKey, *awsRegion)
}
