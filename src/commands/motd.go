package commands

import (
	"zeroctl/src/types"
)

func init() {
	types.RegisterCommand(types.Command{
		Name:        "motd",
		Description: "Display Message of the Day",
		Handler: func() string {
			return "Message of the Day: Stay awesome, zero!\n"
		},
	})
}
