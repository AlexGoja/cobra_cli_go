package main

import (
	cmd "devops-tool/cmd"
	"devops-tool/aws"
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
