/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"path/filepath"

	"os"

	"strings"

	"io/fs"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// jenvCmd represents the jenv command
var jenvCmd = &cobra.Command{
	Use:   "jenv",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("jenv called")

		jdkInstallationDirectory := os.Getenv("JAVA_HOME")

		// get the flag value, its default value is false
		flist, _ := cmd.Flags().GetBool("list")

		if flist {
			listAllJdkInstallations(jdkInstallationDirectory)
		} else {
			fmt.Println("JAVA_HOME: " + jdkInstallationDirectory)
			fmt.Println("To see all available options, run `go-env --help`")
		}

	},
}

func listAllJdkInstallations(jdkInstallationDirectory string) {

	javaDirectory := filepath.Dir(jdkInstallationDirectory)
	fmt.Println("Parent directory: " + javaDirectory)

	directoryList, err := getFilesList(javaDirectory)
	if err != nil {
		panic(err)
	}

	fmt.Println("All available java installations -")

	for index, jdkDirectory := range directoryList {
		fmt.Println(index, ": ", jdkDirectory.Name())
	}
}

func getFilesList(root string) ([]fs.FileInfo, error) {

	var files []fs.FileInfo
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if strings.Contains(file.Name(), "jdk") {
			files = append(files, file)
		}
	}
	return files, nil
}

func init() {
	rootCmd.AddCommand(jenvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jenvCmd.PersistentFlags().String("foo", "", "A help for foo")

	jenvCmd.Flags().BoolP("list", "l", false, "List available jdk installations")
}
