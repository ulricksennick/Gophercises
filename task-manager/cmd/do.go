/*
Copyright © 2020 Nicholas Ulricksen <n.ulricksen@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ulricksennick/gophercises/task-manager/db"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task complete.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var tasks []db.Task
		var completedTask string
		var taskDb db.DB
		var completedDb db.DB

		// Tasks DB connection
		err := taskDb.Open("tasks")
		if err != nil {
			log.Fatal(err)
		}

		// Check task number is valid
		tasks = taskDb.List()
		completedTaskNum, err := strconv.Atoi(args[0])
		if err != nil || !inRange(completedTaskNum, 1, len(tasks)) {
			fmt.Printf("Invalid task number: \"%v\"\n\n", args[0])
			fmt.Println("Run 'task list' to view task numbers.")
			return
		}

		// Delete task from "tasks" bucket
		taskToDelete := tasks[completedTaskNum-1]
		completedTask = string(taskToDelete.Task)
		taskDb.Delete(taskToDelete.Key)
		taskDb.Close()

		// Completed DB connection
		err = completedDb.Open("completed")
		if err != nil {
			log.Fatal(err)
		}

		// Add task to "completed" bucket
		completedDb.Insert(completedTask)
		completedDb.Close()
		fmt.Printf("Task \"%v\" completed.\n", completedTask)
	},
}

func inRange(x int, lower int, upper int) bool {
	return x >= lower && x <= upper
}

func init() {
	rootCmd.AddCommand(doCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
