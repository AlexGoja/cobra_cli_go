package aws

import (
	cmd2 "devops-tool/cmd"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
)

var versionCmd *cobra.Command

type DevToolVersion struct {
	cmd *cobra.Command
	sess *session.Session
}

func CreateDevToolVersion(cmd *cobra.Command, sess *session.Session) *DevToolVersion {
	version := &DevToolVersion{cmd, sess}
	version.getVersion()
	cmd.AddCommand(versionCmd)
	return version
}

func (i *DevToolVersion) getVersion() {
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print current version",
		Long:  `Version of devops-tool`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Devops-tool %s \n", cmd2.VERSION)
		},
	}
}
//
//func summary(roles []*iam.Role) {
//	aiiMap := assembleRolesByAII(roles)
//
//	summary := make(map[string]int)
//
//	for aii, roles := range aiiMap {
//		summary[aii] = len(roles)
//	}
//
//	n := map[int][]string{}
//	var a []int
//	for k, v := range summary {
//		n[v] = append(n[v], k)
//	}
//	for k := range n {
//		a = append(a, k)
//	}
//	sort.Sort(sort.Reverse(sort.IntSlice(a)))
//	var sum int
//	for _, k := range a {
//		for _, s := range n[k] {
//			sum = sum + k
//			fmt.Printf("%s: %d\n", s, k)
//		}
//	}
//	fmt.Println(fmt.Printf("Total: %d", sum))
//
//}
//
//func printRolesByAII(roles []*iam.Role) {
//
//	aiiMap := assembleRolesByAII(roles)
//
//	for key, value := range aiiMap {
//		fmt.Println()
//		fmt.Println()
//		fmt.Println(fmt.Printf("AII = %s ====================================", key))
//		for _, val := range value {
//			fmt.Println(val)
//		}
//		fmt.Println("=======================================")
//		fmt.Println()
//		fmt.Println()
//	}
//}
//
//func assembleRolesByAII(roles []*iam.Role) map[string][]string {
//	aiiMap := make(map[string][]string)
//	for _, role := range roles {
//		r, _ := regexp.Compile("[0-9]{6}")
//		aii := r.FindString(*role.RoleName)
//		if _, ok := aiiMap[aii]; ok {
//			arr := aiiMap[aii]
//			arr = append(arr, *role.RoleName)
//			aiiMap[aii] = arr
//		} else {
//			array := make([]string, 0)
//			array = append(array, *role.RoleName)
//			aiiMap[aii] = array
//		}
//	}
//	return aiiMap
//}
//
//var getRolesRaw = &cobra.Command{
//	Use:   "raw",
//	Short: "Print roles",
//	Long:  `Print roles for specific profile`,
//	Run: func(cmd *cobra.Command, args []string) {
//
//		sess, err := session.NewSessionWithOptions(session.Options{
//			Profile: profile,
//			Config: aws.Config{
//				Region: aws.String(region),
//			},
//		})
//
//		if err != nil {
//			fmt.Println(fmt.Errorf("error creating AWS session: %s", err))
//		}
//
//		sess = session.Must(sess, err)
//
//		iamSvc := iam.New(sess)
//
//		var firstTime = true
//		var output *iam.ListRolesOutput
//		var allRoles []*iam.Role
//		fmt.Println("Loading roles...")
//		for {
//
//			if firstTime == true {
//				output, err = iamSvc.ListRoles(&iam.ListRolesInput{})
//				if err != nil {
//					fmt.Println(fmt.Errorf("request to list roles failed: %s", err))
//				} else {
//					allRoles = append(allRoles, output.Roles...)
//				}
//				firstTime = false
//			}
//
//			if output.IsTruncated != nil && *output.IsTruncated == true {
//				output, err = iamSvc.ListRoles(&iam.ListRolesInput{
//					Marker: output.Marker,
//				})
//				if err != nil {
//					fmt.Println(fmt.Errorf("request to list roles failed: %s", err))
//				} else {
//					allRoles = append(allRoles, output.Roles...)
//				}
//			} else {
//				break
//			}
//		}
//		fmt.Println(allRoles)
//	},
