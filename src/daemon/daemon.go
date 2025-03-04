package daemon

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"zeroctl/src/tasks"
	"zeroctl/src/types"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	_ "zeroctl/src/commands"
)

const socketPath = "/tmp/zeroctl.sock"

func RunDaemon(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	scheduler := &tasks.Scheduler{}
	scheduler.InitScheduler()
	defer scheduler.StopScheduler()

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		logrus.Fatalf("Error starting UNIX socket listener: %v", err)
	}
	defer listener.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("Shutting down daemon...")
		os.Remove(socketPath)
		os.Exit(0)
	}()

	logrus.Infof("Daemon is running and listening on %s", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Error("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		logrus.Error("Error reading from socket: ", err)
		return
	}

	command := string(buffer[:n])
	response := processCommand(command)
	conn.Write([]byte(response))
}

func processCommand(cmdStr string) string {
	cmdStr = strings.TrimSpace(cmdStr)

	if cmd, exists := types.GetCommand(cmdStr); exists {
		return cmd.Handler()
	}

	return "Unknown command\n"
}

func RunClient(command string) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Println("Error connecting to zeroctl daemon:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(command))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Print(string(buffer[:n]))
}
