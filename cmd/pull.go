/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull in vendored dependences and validate",
	Long: `Vendor each of the depedencies in the project ( Go, Git, Python, etc ), and validate against the highest integrity source`,
	Example: `  vendorme pull foo

	# vendor dependencies for everything in the current project
	vendorme pull . 	
	
	# vendor dependencies for go
	vendorme pull go .
	
	# vendor dependencies for git
	vendorme pull git . `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pull called")
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
