// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/forgerock/frconfig/amconfig"
)

// policyCmd represents the policy command
var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Get policy objects from OpenAM",
	Long: `Get Policies from openam`,
	Run: func(cmd *cobra.Command, args []string) {
		am, err := amconfig.GetOpenAMConnection()
		if err != nil {
			msg := fmt.Sprintf("Can't create OpenAM connection: %v", err)
			fatal(msg)
		}

		// get output format
		o,err := cmd.Flags().GetString("output")

		str,err := am.ExportPolicies(o)
		fmt.Println(str)

	},
}

func init() {
	getCmd.AddCommand(policyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// policyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// policyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
