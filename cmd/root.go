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
	"os"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string


var (

	OutputFileTypes  = []string{ "yaml", "json"}
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "frconfig",
	Short: "frconfig - import and export configurations",
	Long: `Import and Export Configurations from the ForgeRock stack.

	Example:

	frconfig get policy -o yaml


	.`,
// Uncomment the following line if your bare application
// has an action associated with it:
//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)


	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/frconfig.yaml.yaml)")

	RootCmd.PersistentFlags().String("realm","", "realm - set OpenAM Realm to act on")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("frconfig") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AddConfigPath(".")  // adding current dir

	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Viper err ", err)
	}

	//fmt.Printf("Viper url= %v\n", viper.GetString("default.openam.url"))
}
