package commands

import (
	"errors"
	"zeroctl/src/database"
	"zeroctl/src/types"
)

func init() {
	types.RegisterCommand(types.Command{
		Name:        "toggle",
		Description: "Enable or disable a task",
		Handler: func(args []string) string {
			if len(args) < 1 {
				return "Please provide a task name\n"
			}

			taskName := args[0]
			status, err := database.GetValue(taskName)

			if err != nil && !errors.Is(err, database.ErrKeyNotFound) {
				return "Failed to get task status: " + err.Error() + "\n"
			}

			if status == "" {
				err = database.StoreValue(taskName, "disabled")
				if err != nil {
					return "Failed to disable task: " + err.Error() + "\n"
				}
				return "Task '" + taskName + "' is now disabled\n"
			}

			newStatus := "enabled"
			if status == "enabled" {
				newStatus = "disabled"
			}

			err = database.StoreValue(taskName, newStatus)
			if err != nil {
				return "Failed to toggle task: " + err.Error() + "\n"
			}

			return "Task '" + taskName + "' is now " + newStatus + "\n"
		},
	})
}
