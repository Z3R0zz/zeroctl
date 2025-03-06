package main

import (
	"fmt"
	"os"
	"zeroctl/src/config"
	"zeroctl/src/daemon"
	"zeroctl/src/database"
	"zeroctl/src/handlers"
	"zeroctl/src/types"
	"zeroctl/src/utils"

	"github.com/spf13/cobra"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 && os.Args[1] == "daemon" {
		utils.SetStartTime()

		if err := database.InitBoltDB(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer database.CloseBoltDB()
	}

	rootCmd := &cobra.Command{
		Use:   "zeroctl",
		Short: "zeroctl is a custom CLI tool",
	}

	daemonCmd := &cobra.Command{
		Use:   "daemon",
		Short: "Start the zeroctl daemon (runs tasks and listens for commands)",
		PreRun: func(cmd *cobra.Command, args []string) {
			go startupJobs()
		},
		Run: daemon.RunDaemon,
	}

	for cmdName, cmd := range types.CommandRegistry {
		localCmd := &cobra.Command{
			Use:   cmdName,
			Short: cmd.Description,
			Run: func(cobraCmd *cobra.Command, args []string) {
				daemon.RunClient(cmdName)
			},
		}
		rootCmd.AddCommand(localCmd)
	}

	rootCmd.AddCommand(daemonCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func startupJobs() {
	err := handlers.CacheWeatherData()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
