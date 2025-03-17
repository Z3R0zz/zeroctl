package commands

import (
	"fmt"
	"runtime"
	"strings"
	"zeroctl/src/types"
	"zeroctl/src/utils"
)

func init() {
	types.RegisterCommand(types.Command{
		Name:        "stats",
		Description: "Get stats about zeroctl and its usage",
		Handler: func(args []string) string {
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			memAlloc := fmt.Sprintf("%v MB", utils.BToMb(memStats.Alloc))
			memTotalAlloc := fmt.Sprintf("%v MB", utils.BToMb(memStats.TotalAlloc))
			memSys := fmt.Sprintf("%v MB", utils.BToMb(memStats.Sys))

			numCPU := runtime.NumCPU()
			numGoroutines := runtime.NumGoroutine()

			goVersion := runtime.Version()

			osInfo := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

			var sb strings.Builder

			sb.WriteString("┌─────────────── ZEROCTL STATS ────────────────┐\n")

			sb.WriteString(fmt.Sprintf("│ Go Version:   %-30s │\n", goVersion))
			sb.WriteString(fmt.Sprintf("│ OS/Arch:      %-30s │\n", osInfo))
			sb.WriteString("├──────────────────────────────────────────────┤\n")

			sb.WriteString(fmt.Sprintf("│ CPU Cores:    %-30d │\n", numCPU))
			sb.WriteString(fmt.Sprintf("│ Goroutines:   %-30d │\n", numGoroutines))
			sb.WriteString("├──────────────────────────────────────────────┤\n")

			sb.WriteString(fmt.Sprintf("│ Mem Alloc:    %-30s │\n", memAlloc))
			sb.WriteString(fmt.Sprintf("│ Total Alloc:  %-30s │\n", memTotalAlloc))
			sb.WriteString(fmt.Sprintf("│ Sys Memory:   %-30s │\n", memSys))
			sb.WriteString(fmt.Sprintf("│ GC Cycles:    %-30d │\n", memStats.NumGC))
			sb.WriteString("└──────────────────────────────────────────────┘")

			return sb.String()
		},
	})
}
