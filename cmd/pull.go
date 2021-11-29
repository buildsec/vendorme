//
// Copyright 2021 Tim Miller.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (

	"github.com/spf13/cobra"
	"github.com/trmiller/vendorme/cmd/cli"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull in vendored dependences and validate",
	Long: `Vendor each of the depedencies in the project ( Go, Git, Python, etc ), and validate against the highest integrity source`,
	Example: `  vendorme pull foo

	# vendor dependencies for everything in the current project
	vendorme pull
	
	# vendor dependencies specifying config
	vendorme pull ./myvendorconfig.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := cli.PullCommand{
			VendorMeConfig:	"vendor.yaml",
		}
		return p.Exec(cmd.Context(), args)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
