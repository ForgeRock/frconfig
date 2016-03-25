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

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get an object type",
	Long: `Get a configuration object.
	Example    get policy`,
	// Commented out - since a subcommand is needed
	//Run: func(cmd *cobra.Command, args []string) {
	//	// TODO: Work your own magic here
	//	fmt.Println("get called")
	//},
}

func init() {
	RootCmd.AddCommand(getCmd)

	// add output option

	getCmd.PersistentFlags().StringP("output", "o", "json", "Output format: json or yaml")
	//getCmd.PersistentFlags().Lookup("output").NoOptDefVal = "json"


}
