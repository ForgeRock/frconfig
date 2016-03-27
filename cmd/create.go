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
	"path/filepath"
)

type CreateOptions struct {
	Filenames []string
}

var (
	overwrite, continueOnError  bool
)


func init() {
	options := &CreateOptions{}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a configuration object",
		Long: `frconfig create -f object.json

	Create a configuration object`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(options.Filenames) == 0 {
				cmd.Help()
				return
			}
			overwrite, _ = cmd.Flags().GetBool("overwrite")
			continueOnError, _ = cmd.Flags().GetBool("continue")
			//fmt.Println("create called file = ", options.Filenames)

			for _,v := range options.Filenames {
				err := createObject(v)
				if err != nil {
					fatal(fmt.Sprintf("Can't create object, err %v", err))
				}
			}
		},
	}

	RootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().Bool("overwrite",true,"If true, overwrite any existing value")


	usage := "Filename, directory, or URL to file to use to create the resource"

	AddJsonFilenameFlag(createCmd, &options.Filenames, usage)

	// there could be sub commands - so -f is not required
	// createCmd.MarkFlagRequired("filename")

}

// Create the object described by fileName
// If fileName is - read from stdin
// If fileName is a directory, recurse and read all files (*.json, *.yaml) in that diretory
func createObject(fileName string) (err error) {
	if fileName == "-" {
		return createStream(os.Stdin)
	}
	err = filepath.Walk(fileName, visit)
	return err
}

func visit(path string, f os.FileInfo, err error) error {
	//fmt.Printf("Visited: %s\n", path)
	if  ! f.IsDir() {
		ext := filepath.Ext(path)
		// skip any file that is not a json or yaml format
		if ext != "json" && ext != "yml" && ext != "yaml" {
			return nil
		}
		file,err := os.Open(path)
		defer file.Close()
		if err != nil {
			return err
		}
		return createStream( file )
	}
	return nil
}


func createStream(f io.Reader) (err error) {
	err = amconfig.CreateFRObjects(f, overwrite,continueOnError)
	return err
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

