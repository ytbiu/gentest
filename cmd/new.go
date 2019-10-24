// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"gentest/pkg/parser"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called")
		run()
	},
}

var (
	fileName        string
	mockSrcFileName string
	methodName      string
	mockPath        string
	doGoGet         bool
)

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.PersistentFlags().StringVarP(&fileName, "file", "f", "", "文件名")
	newCmd.PersistentFlags().StringVarP(&mockSrcFileName, "mock_src", "s", "", "mock源文件名")
	newCmd.PersistentFlags().StringVarP(&methodName, "method", "m", "", "函数名")
	newCmd.PersistentFlags().StringVarP(&mockPath, "mock", "o", "mock", "mock目标文件路径")
	newCmd.PersistentFlags().BoolVarP(&doGoGet, "goget", "g", false, "是否执行go get下载单元测试依赖(default false)")

}

func run() {
	engine := parser.NewEngine(
		parser.FileName(fileName),
		parser.MockSrcFileName(mockSrcFileName),
		parser.MethodName(methodName),
		parser.MockPath(mockPath),
		parser.DoGoGet(doGoGet),
	)
	engine.Run()
}
