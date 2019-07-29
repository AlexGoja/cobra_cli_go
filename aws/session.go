package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type conf struct {
	Profile string `yaml:"profile"`
	Region  string `yaml:"region"`
}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &conf {
		Profile: c.Profile,
		Region:  c.Region,
	}
}

/*
	Create an AWS session to enable AWS API interaction
 */
func CreateAWSSession() *session.Session {

	var c conf
	c.getConf()

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: c.Profile,
		Config: aws.Config{
			Region: aws.String(c.Region),
		},
	})

	if err != nil {
		fmt.Println(fmt.Errorf("error creating AWS session: %s", err))
		panic(err)
	}

	sess = session.Must(sess, err)

	return sess
}


