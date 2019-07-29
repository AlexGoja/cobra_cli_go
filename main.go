package main

import (
	cmd "devops-tool/cmd"
	"devops-tool/aws"
)

func main() {

	root := cmd.CreateRootCommand()
	sess := aws.CreateAWSSession()

	aws.CreateDevToolVersion(root, sess)
	aws.CreateEC2(root, sess)
	aws.CreateCloudformation(root, sess)
	root.Execute()
}
