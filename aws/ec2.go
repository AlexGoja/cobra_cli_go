package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var getResourcesCmd *cobra.Command
var getTaggedEC2Cmd *cobra.Command

type EC2Resource struct {
	cmd *cobra.Command
	sess *session.Session
}

func CreateEC2(cmd *cobra.Command, sess *session.Session) *EC2Resource {
	ec2R := &EC2Resource{cmd, sess}
	return ec2R
}

func (e *EC2Resource) describeKeyPair(keyNames []*string) ([]*ec2.KeyPairInfo, error) {

	ec2Client := ec2.New(e.sess, e.sess.Config)

	req, resp := ec2Client.DescribeKeyPairsRequest(&ec2.DescribeKeyPairsInput{
		DryRun:   aws.Bool(false),
		Filters:  nil,
		KeyNames: keyNames,
	})

	err := req.Send()

	return resp.KeyPairs, err
}

func (e *EC2Resource) createKeyPair(keyPairName string) (*ec2.CreateKeyPairOutput, error) {

	ec2Client := ec2.New(e.sess, e.sess.Config)

	req, resp := ec2Client.CreateKeyPairRequest(&ec2.CreateKeyPairInput{
		DryRun:  aws.Bool(false),
		KeyName: aws.String(KEY_NAME),
	})

	err := req.Send()

	return resp, err
}

func (e *EC2Resource) deleteKeyPair(keyPairName string) (*ec2.DeleteKeyPairOutput, error) {

	ec2Client := ec2.New(e.sess, e.sess.Config)

	req, resp := ec2Client.DeleteKeyPairRequest(&ec2.DeleteKeyPairInput{
		DryRun:  aws.Bool(false),
		KeyName: aws.String(KEY_NAME),
	})

	err := req.Send()

	return resp, err
}






