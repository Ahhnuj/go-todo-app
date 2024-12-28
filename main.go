package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type Todo struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var TodoList []Todo
var taskID int
var filename = "todo.json"

func loadTodoList() error {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&TodoList); err != nil && err.Error() != "EOF" {
		return err
	}

	if len(TodoList) > 0 {
		taskID = TodoList[len(TodoList)-1].ID
	}
	return nil
}

func saveTodoList() error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(TodoList); err != nil {
		return err
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A simple CLI Todo app",
}

var addCmd = &cobra.Command{
	Use:   "add [task]",
	Short: "Add a new task to the Todo list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		task := args[0]
		taskID++
		newTask := Todo{
			ID:        taskID,
			Task:      task,
			Completed: false,
			CreatedAt: time.Now(),
		}
		TodoList = append(TodoList, newTask)
		if err := saveTodoList(); err != nil {
			fmt.Println("Error saving todo list:", err)
		} else {
			fmt.Printf("Task '%s' added at %s\n", task, newTask.CreatedAt.Format(time.RFC1123))
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		if len(TodoList) == 0 {
			fmt.Println("No tasks found.")
			return
		}
		fmt.Println("ID | Task                    | Status      | Created At / Completed At")
		fmt.Println("-----------------------------------------------------------")
		for _, task := range TodoList {
			status := "Incomplete"
			if task.Completed {
				status = "Completed"
			}
			fmt.Printf("%-3d | %-23s | %-12s | %s\n", task.ID, task.Task, status, task.CreatedAt.Format(time.RFC1123))
		}
	},
}

var toggleCmd = &cobra.Command{
	Use:   "toggle [taskID]",
	Short: "Toggle the status of a task between Complete and Incomplete",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		var found bool
		for i, task := range TodoList {
			if fmt.Sprintf("%d", task.ID) == id {

				TodoList[i].Completed = !TodoList[i].Completed
				if err := saveTodoList(); err != nil {
					fmt.Println("Error saving todo list:", err)
				} else {
					status := "Incomplete"
					if TodoList[i].Completed {
						status = "Completed"
					}
					fmt.Printf("Task '%s' marked as %s.\n", task.Task, status)
				}
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Task ID not found.")
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove [taskID]",
	Short: "Remove a task by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		var index int
		var found bool
		for i, task := range TodoList {
			if fmt.Sprintf("%d", task.ID) == id {
				index = i
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Task ID not found.")
			return
		}

		TodoList = append(TodoList[:index], TodoList[index+1:]...)
		if err := saveTodoList(); err != nil {
			fmt.Println("Error saving todo list:", err)
		} else {
			fmt.Printf("Task ID %s removed.\n", id)
		}
	},
}

func main() {

	if err := loadTodoList(); err != nil {
		fmt.Println("Error loading todo list:", err)
	}

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(toggleCmd)
	rootCmd.AddCommand(removeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
