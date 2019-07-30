package main

import (
	cmd "cobra_cli_go/cmd"
	"cobra_cli_go/aws"
)

func main() {

	// Creates the devops-tool root command
	root := cmd.CreateRootCommand()

	// Creates aws sessions
	sess := aws.CreateAWSSession()

	//Instantiate clients
	aws.CreateDevToolVersion(root, sess)
	aws.CreateEC2(root, sess)
	aws.CreateCloudformation(root, sess)
	root.Execute()
}
