package commands

import (
	"fmt"
	"zeroctl/src/types"
	"zeroctl/src/utils"
)

func init() {
	types.RegisterCommand(types.Command{
		Name:        "uptime",
		Description: "Display how long the zeroctl daemon has been running",
		Handler: func() string {
			uptime := utils.GetUptime()

			hours := int(uptime.Hours()) % 24
			minutes := int(uptime.Minutes()) % 60
			seconds := int(uptime.Seconds()) % 60

			return fmt.Sprintf("Daemon uptime: %d hours, %d minutes, %d seconds",
				hours, minutes, seconds)
		},
	})
}
