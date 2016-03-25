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
	"strings"
	"os"
	"io"
	"github.com/forgerock/frconfig/amconfig"
)

type CreateOptions struct {
	Filenames []string
}


func init() {
	options := &CreateOptions{}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a configuration object",
		Long: `frconfig create -f object.json

	Create a configuration object`,

		Run: func(cmd *cobra.Command, args []string) {

			var err error
			if len(options.Filenames) == 0 {
				cmd.Help()
				return
			}

			//fmt.Println("create called file = ", options.Filenames)

			for _,v := range options.Filenames {
				//fmt.Printf("File name = %s", v)
				var  f io.Reader
				if v == "-" {
					f = os.Stdin
				} else {
					f,err = os.Open(v)
				}

				if err != nil {
					msg := fmt.Sprintf("Can't open %s", v)
					fatal(msg)

				}
				err = amconfig.CreateFRObjects(f)
				if err != nil {
					fatal(fmt.Sprintf("Error reading config %v ",err))
				}

			}
		},
	}

	RootCmd.AddCommand(createCmd)


	usage := "Filename, directory, or URL to file to use to create the resource"

	AddJsonFilenameFlag(createCmd, &options.Filenames, usage)

	// there could be sub commands - so -f is not required
	// createCmd.MarkFlagRequired("filename")

}

var FileExtensions = []string{".json", ".yaml", ".yml"}
//var InputExtensions = append(FileExtensions, "stdin")

func AddJsonFilenameFlag(cmd *cobra.Command, value *[]string, usage string) {
	cmd.Flags().StringSliceVarP(value, "filename", "f", *value, usage)
	annotations := []string{}

	// todo: where to set these
	for _, ext := range FileExtensions {
		annotations = append(annotations, strings.TrimLeft(ext, "."))
	}
	cmd.Flags().SetAnnotation("filename", cobra.BashCompFilenameExt, annotations)
}

