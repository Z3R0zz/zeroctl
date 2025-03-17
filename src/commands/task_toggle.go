package commands

import (
	"fmt"
	"strings"
	"zeroctl/src/database"
	"zeroctl/src/tasks"
	"zeroctl/src/types"
)

var scheduler *tasks.Scheduler

func SetScheduler(s *tasks.Scheduler) {
	scheduler = s
}

func init() {
	types.RegisterCommand(types.Command{
		Name:        "toggle",
		Description: "Enable or disable a task",
		Handler: func(args []string) string {
			if len(args) < 1 {
				return "Please provide a task name\n"
			}

			taskName := args[0]

			availableTasks := scheduler.GetAvailableTasks()
			if _, exists := availableTasks[taskName]; !exists {
				return fmt.Sprintf("Task '%s' does not exist. Available tasks:\n%s\n",
					taskName, getAvailableTasksList())
			}

			status, err := scheduler.GetTaskStatus(taskName)
			if err != nil && err != database.ErrKeyNotFound {
				return "Failed to get task status: " + err.Error() + "\n"
			}

			newStatus := "enabled"
			if status == "enabled" {
				newStatus = "disabled"
			}

			err = database.StoreValue(taskName, newStatus)
			if err != nil {
				return "Failed to toggle task: " + err.Error() + "\n"
			}

			return fmt.Sprintf("Task '%s' is now %s\n", taskName, newStatus)
		},
	})

	types.RegisterCommand(types.Command{
		Name:        "tasks",
		Description: "List all available tasks and their status",
		Handler: func(args []string) string {
			if scheduler == nil {
				return "Scheduler not initialized\n"
			}

			availableTasks := scheduler.GetAvailableTasks()
			if len(availableTasks) == 0 {
				return "No tasks available\n"
			}

			var result strings.Builder
			result.WriteString("Available tasks:\n")

			for name, task := range availableTasks {
				status, err := scheduler.GetTaskStatus(name)
				if err != nil {
					status = "unknown"
				}

				result.WriteString(fmt.Sprintf("- %s: %s [%s]\n", name, task.Description, status))
			}

			return result.String()
		},
	})
}

func getAvailableTasksList() string {
	if scheduler == nil {
		return "Scheduler not initialized"
	}

	availableTasks := scheduler.GetAvailableTasks()
	if len(availableTasks) == 0 {
		return "No tasks available"
	}

	var names []string
	for name := range availableTasks {
		names = append(names, name)
	}

	return "  - " + strings.Join(names, "\n  - ")
}
