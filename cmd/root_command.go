package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

const VERSION = "v0.3"

type RCommand struct {
	command *cobra.Command
}

func CreateRootCommand() *cobra.Command {
	return &cobra.Command {
			Use:   "devops-tool",
			Short: "Devops tool to manage different processes",
			Long:  `Tool for automation with cobra and go`,
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Printf("Devops-tool %s \n", VERSION)
			},
	}
}

func (r *RCommand) Execute() {
	if err := r.command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
