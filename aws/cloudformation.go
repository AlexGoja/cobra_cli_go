package aws

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var createStack *cobra.Command
var deleteStack *cobra.Command

const (
	KEY_NAME = "efskeypair"
	INSTANCE_TYPE = "t2.micro"
	FILE = "efskeypair.pem"
)

type Cloudformation struct {
	cmd *cobra.Command
	sess *session.Session
	ec2Res *EC2Resource
}

func CreateCloudformation(cmd *cobra.Command, sess *session.Session) *Cloudformation {
	ec2res := &EC2Resource{cmd, sess}
	cf := &Cloudformation{cmd, sess, ec2res}
	cf.createStack()
	cf.deleteEFS()
	createStack.Flags().StringP("onfailure", "f", "ROLLBACK", "flag used for on stack creation failure")
	createStack.Flags().StringArrayP("capabilities", "c", nil, "flag for injecting capabilities")
	createStack.Flags().StringP("stackname", "s", "", "flag used to pass stack name")
	createStack.Flags().StringP("instance-type", "i", INSTANCE_TYPE, "instance type flag")
	createStack.Flags().StringP("key-pair-name", "k", KEY_NAME, "key pair flag type flag")
	createStack.Flags().StringArrayP("subnets", "n", nil, "list of subnets")
	createStack.Flags().StringP("vpc", "v", "", "flag for vpc id")
	createStack.Flags().StringP("image", "m", "", "flag for image id")

	_ = createStack.MarkFlagRequired("stackname")
	_ = createStack.MarkFlagRequired("subnets")
	_ = createStack.MarkFlagRequired("vpc")
	_ = createStack.MarkFlagRequired("image")

	deleteStack.Flags().StringP("stackname", "s", "", "flag used to pass stack name")
	deleteStack.Flags().StringP("key-pair-name", "k", KEY_NAME, "key pair flag type flag")
	_ = deleteStack.MarkFlagRequired("stackname")
	_ = deleteStack.MarkFlagRequired("key-pair-name")

	cmd.AddCommand(deleteStack)
	cmd.AddCommand(createStack)
	return cf
}

/**
	Creates efs file system using cloudfromation stack
 */
func (cf *Cloudformation) createStack() {

	createStack = &cobra.Command {
		Use:   "create-efs",
		Short: "Create stack",
		Long:  `Create stack based on template`,
		Run: func(cmd *cobra.Command, args []string) {

			parameters := make([]*cloudformation.Parameter, 0, 5)

			onFailure, err := cmd.Flags().GetString("onfailure")
			verifyOnFailureArg(err, onFailure)

			capabilities, err := cmd.Flags().GetStringArray("capabilities")
			verifyCapability(err, capabilities)

			stackName, err := cmd.Flags().GetString("stackname")
			verify(err, "stackname")

			instanceType, err := cmd.Flags().GetString("instance-type")
			verify(err, "instance-type")
			addParameter(&parameters, "InstanceType", instanceType)

			keyPair, err := cmd.Flags().GetString("key-pair-name")
			verify(err, "key-pair-name")
			addParameter(&parameters, "KeyName", keyPair)

			subnets, err := cmd.Flags().GetStringArray("subnets")
			verify(err, "subnets")
			addParameters(&parameters, "Subnets", subnets)

			vpc, err := cmd.Flags().GetString("vpc")
			verify(err, "vpc")
			addParameter(&parameters, "VPC" ,vpc)

			image, err := cmd.Flags().GetString("image")
			verify(err, "image")
			addParameter(&parameters, "ImageId", image)

			for _, param := range parameters {
				fmt.Printf("----Param : %s ---- \n", param)
			}

			cfClient := cloudformation.New(cf.sess, cf.sess.Config)

			resp, err := cf.ec2Res.createKeyPair(keyPair)

			if err != nil {
				fmt.Printf("Create key pair err: %s \n", err)
				os.Exit(1)
			}

			createPemFile(resp.KeyMaterial)

			file := readFile("templates/distributed_file_system.yaml")

			if _, err := cfClient.CreateStack(&cloudformation.CreateStackInput{
				Capabilities:                aws.StringSlice(capabilities),
				ClientRequestToken:          nil,
				DisableRollback:             nil,
				EnableTerminationProtection: nil,
				NotificationARNs:            nil,
				OnFailure:                   &onFailure,
				Parameters:                  parameters,
				ResourceTypes:               nil,
				RoleARN:                     nil,
				RollbackConfiguration:       nil,
				StackName:                   &stackName,
				StackPolicyBody:             nil,
				StackPolicyURL:              nil,
				Tags:                        nil,
				TemplateBody:                &file,
				TemplateURL:                 nil,
				TimeoutInMinutes:            nil,
			}); err != nil {
				deleteKeyPair(cf, keyPair, err)
			}

		},
	}
}

/**
	Creates pem file with the rsa key to ssh to ec2 instance
*/
func createPemFile(rsa *string) {
	f, err := os.Create(FILE)
	defer f.Close()

	if err != nil {
		fmt.Printf("Could not create file %s \n", FILE)
	}

	w := bufio.NewWriter(f)
	_, err = w.WriteString(*rsa)

	if err != nil {
		fmt.Printf("Could not write to file %s \n", FILE)
	}

	err = os.Chmod(FILE, 0400)
	if err != nil {
		fmt.Printf("Could not set permission on file %s \n", err)
	}

	w.Flush()

}

/**
	Deletes pem file with the rsa key to ssh to ec2 instance
*/
func deletePemFile() {
	err := os.Remove(FILE)

	if err != nil {
		fmt.Println(err)
		return
	}
}

/**
	Deletes all resources created for efs and general clean-up
*/
func (cf *Cloudformation) deleteEFS() {

	deleteStack = &cobra.Command {
		Use:   "delete-efs",
		Short: "Delete stack",
		Long:  `Delete efs stack resources`,
		Run: func(cmd *cobra.Command, args []string) {

			stackName, err := cmd.Flags().GetString("stackname")
			verify(err, "stackname")

			keyPair, err := cmd.Flags().GetString("key-pair-name")
			verify(err, "key-pair-name")

			cfClient := cloudformation.New(cf.sess, cf.sess.Config)

			if _, err := cfClient.DeleteStack(&cloudformation.DeleteStackInput{
				ClientRequestToken: nil,
				RetainResources:    nil,
				RoleARN:            nil,
				StackName:          &stackName,
			}); err != nil {
				fmt.Printf("Stack failed to delete with err: %s \n", err)
			}

			deleteKeyPair(cf, keyPair, err)

			deletePemFile()

		},
	}
}

/**
	Clean-up operation in case something critical fails to be created.
 */
func deleteKeyPair(cf *Cloudformation, keyPair string, err error) {
	if err != nil {
		fmt.Printf("Key pair deleted %s because stack failed with err: \n %s \n", keyPair, err)
		os.Exit(1)
	}
	_, keyErr := cf.ec2Res.deleteKeyPair(keyPair)
	if keyErr != nil {
		fmt.Printf("Delete key pair err: %s \n", keyErr)
	}
}

/**
	Read the cloudformation template and return as string.
	Used to inject in stack creation
*/
func readFile (path string) string {
	bytesStrFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}
	return string(bytesStrFile)
}

/**
	Create parameter to inject to cloudformation stack in teamplates dir
*/
func addParameter(p *[]*cloudformation.Parameter, key string, value string) {

	*p = append(*p, &cloudformation.Parameter{
		ParameterKey:     &key,
		ParameterValue:   &value,
		ResolvedValue:    nil,
		UsePreviousValue: nil,
	})
}

/**
	Create list of parameters to inject to cloudformation parameter arg in teamplates dir
*/
func addParameters(p *[]*cloudformation.Parameter, key string, values []string) {

	for _, values := range values {
		*p = append(*p, &cloudformation.Parameter{
			ParameterKey:     &key,
			ParameterValue:   &values,
			ResolvedValue:    nil,
			UsePreviousValue: nil,
		})
	}
}

/**
	Checks if flag was retrieved ok.
*/
func verify(err error, str interface{}) {
	if err != nil {
		fmt.Printf("Error retrieving flag for %s \n", str)
		os.Exit(1)
	}
}

/**
	Checks capabilities flag to be one of the accepted value.
	Stop if value is wrong as stack will fail to create anyway.
*/
func verifyCapability(err error, capabilities []string) {
	if capabilities == nil {
		return
	}

	if err != nil {
		fmt.Printf("problem with retrieving the capability flag with err: %s \n", err)
	}

	for _, capability := range capabilities {
		if capability != cloudformation.CapabilityCapabilityIam  &&
			capability !=  cloudformation.CapabilityCapabilityNamedIam {
			fmt.Printf("capability provided with invalid flag, valid are \"CAPABILITY_IAM\" or \"CAPABILITY_NAMED_IAM\" \n")
			os.Exit(1)
		}
	}
}

/**
	Checks onFailure flag to be one of the accepted value.
	Stop if value is wrong as stack will fail to create anyway.
 */
func verifyOnFailureArg(err error, onFailure string) {
	if err != nil && (
		onFailure == cloudformation.OnFailureDelete ||
			onFailure == cloudformation.OnFailureRollback ||
			onFailure == cloudformation.OnFailureDoNothing) {
		fmt.Printf("onFailure provided with invalid flag with err: %s \n", err)
		os.Exit(1)
	}
}
