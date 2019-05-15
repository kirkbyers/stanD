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
	"strconv"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a day of tasks",
	Long: `Brings up all tasks for a given day offset from the previous day. 
If an offset is not given, then the previous daywill be given For example:

stanD get # Gives yesterdays tasks
stanD get -1 # Gives tasks recorded 2 days ago`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			offset int
			err    error
		)
		if len(args) <= 0 {
			offset = -1
		} else {
			offset, err = strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
		}
		day, err := state.GetDay(offset)
		if err != nil {
			fmt.Printf("There was nothing recorded for %v\n", day.Timestamp.Format("Jan 2 2006"))
			return
		}
		fmt.Printf("For %v, you recorded the following:\n", day.Timestamp.Format("Jan 2 2006"))
		for _, t := range day.Tasks {
			fmt.Printf("\t- %v\n", t.Body)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
