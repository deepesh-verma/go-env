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
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

type FullFileInfo struct {
	directory string
	fileInfo  fs.FileInfo
}

func (fullFileInfo FullFileInfo) String() string {
	return filepath.Join(fullFileInfo.directory, fullFileInfo.fileInfo.Name())
}

const JAVA_HOME_ENV = "JAVA_HOME"

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

		jdkInstallationDirectory := os.Getenv(JAVA_HOME_ENV)

		// get the flag value, its default value is false
		flist, _ := cmd.Flags().GetBool("list")
		fsetVersion, _ := cmd.Flags().GetBool("set-version")

		if flist {
			listAllJdkInstallations(jdkInstallationDirectory)
		} else if fsetVersion {
			setJavaHomeWithJdkInstallation(jdkInstallationDirectory, 0)
		} else {
			fmt.Println(JAVA_HOME_ENV, ": "+jdkInstallationDirectory)
			fmt.Println("To see all available options, run `go-env --help`")
		}
	},
}

func setJavaHomeWithJdkInstallation(jdkInstallationDirectory string, index int) {

	javaDirectory := filepath.Dir(jdkInstallationDirectory)
	fmt.Println("Parent directory: " + javaDirectory)

	fullDirectoryInfoList, err := getFilesList(javaDirectory)
	if err != nil {
		panic(err)
	}

	selectedJdkFullFileInfo := fullDirectoryInfoList[index]
	fmt.Println("Going to use the java installation: ", selectedJdkFullFileInfo)

	osSpecificCommand := getOsSpecificSetCommand(selectedJdkFullFileInfo)

	fmt.Println("Please use the below commmand to set jdk version")
	fmt.Println(osSpecificCommand)
}

func getOsSpecificSetCommand(selectedJdkFullFileInfo FullFileInfo) string {
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprint("SETX ", JAVA_HOME_ENV, " \""+selectedJdkFullFileInfo.String()+"\"")
	case "linux":
		return fmt.Sprint("export ", JAVA_HOME_ENV, " \""+selectedJdkFullFileInfo.String()+"\"")
	default:
		panic("Not yet supported" + runtime.GOOS)
	}
}

func listAllJdkInstallations(jdkInstallationDirectory string) {

	javaDirectory := filepath.Dir(jdkInstallationDirectory)
	fmt.Println("Parent directory: " + javaDirectory)

	fullDirectoryInfoList, err := getFilesList(javaDirectory)
	if err != nil {
		panic(err)
	}

	fmt.Println("All available java installations -")

	for index, fullJdkDirectoryInfo := range fullDirectoryInfoList {
		fmt.Println(index, ": ", fullJdkDirectoryInfo)
	}
}

func getFilesList(root string) ([]FullFileInfo, error) {

	var files []FullFileInfo
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if strings.Contains(file.Name(), "jdk") {
			files = append(files, FullFileInfo{directory: root, fileInfo: file})
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
	jenvCmd.Flags().BoolP("set-version", "s", false, "Set jdk installations as JAVA_HOME")
}
