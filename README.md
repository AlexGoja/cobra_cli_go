# Cobra CLI
Example of cobra go cli interacting with AWS API for automation tasks

# Structure

|____cmd <br>
| &nbsp;&nbsp;&nbsp;&nbsp; |____root_command.go -> root command definition <br>
|____LICENSE -> licence <br>
|____Makefile -> installation Makefile <br>
|____config.yaml -> defining aws config profile and region <br>
|____templates <br>
| &nbsp;&nbsp;&nbsp;&nbsp; |____distributed_file_system.yaml -> keeps cloudformation stacks <br>
|____aws <br>
| &nbsp;&nbsp;&nbsp;&nbsp; |____version.go -> command to determine version of tool<br>
| &nbsp;&nbsp;&nbsp;&nbsp; |____session.go -> creates an aws session <br>
| &nbsp;&nbsp;&nbsp;&nbsp; |____cloudformation.go -> cloudformation aws client <br>
| &nbsp;&nbsp;&nbsp;&nbsp; |____ec2.go -> ec2 aws client <br>
|____main.go <br>

# Install Devops-tool

1. Prerequisites

Go installation https://golang.org/doc/install

2. Go to project root

`make install `

3. You will need a user with relevant permissions to interact with the tool. Setup aws with the user's credentials 

`aws configure --profile <your_profile>`

4. Add <your_profile> to config.yaml

5. To interact with devops-tool examples

`devops-tool help`<br>
`devops-tool create-efs help`<br>
`devops-tool delete-efs help`

Example creating efs: <br>
`devops-tool create-efs -k efskeypair -s efsstack -v vpc-703***09 -m ami-035b3c7efe6d061d5 -n=subnet-49b**02,subnet-55d***79 -c CAPABILITY_NAMED_IAM`

Example deleting efs:<br>
`devops-tool delete-efs -k efskeypair -s efsstack`


# Notes

Bug:

pem file created will need to have chmod modified
chmod 400 efskeypair.pem before connecting to ec2 instance


