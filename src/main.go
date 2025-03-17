package main

import (
	"os"
	"zeroctl/src/config"
	"zeroctl/src/daemon"
	"zeroctl/src/database"
	"zeroctl/src/handlers"
	"zeroctl/src/types"
	"zeroctl/src/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		logrus.Fatalf("Error: %v", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "daemon" {
		utils.SetStartTime()

		if err := database.InitBoltDB(); err != nil {
			logrus.Fatalf("Error: %v", err)
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
				daemon.RunClient(cmdName, args)
			},
		}
		rootCmd.AddCommand(localCmd)
	}

	rootCmd.AddCommand(daemonCmd)

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Error: %v", err)
	}
}

func startupJobs() {
	err := handlers.CacheWeatherData()
	if err != nil {
		logrus.Errorf("Error: %v", err)
	}
}
