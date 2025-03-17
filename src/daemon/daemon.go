package daemon

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	_ "zeroctl/src/commands"
	"zeroctl/src/tasks"
	"zeroctl/src/types"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

type CommandMessage struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		logrus.Error("Error reading from socket: ", err)
		return
	}

	var cmdMsg CommandMessage
	if err := json.Unmarshal(buffer[:n], &cmdMsg); err != nil {
		logrus.Errorf("Error unmarshaling command message: %v", err)
		conn.Write([]byte("Error parsing command message\n"))
		return
	}

	response := processCommand(cmdMsg.Command, cmdMsg.Args)
	conn.Write([]byte(response))
}

func processCommand(cmdStr string, args []string) string {
	cmdStr = strings.TrimSpace(cmdStr)
	if cmd, exists := types.GetCommand(cmdStr); exists {
		return cmd.Handler(args)
	}
	return "Unknown command\n"
}

func RunClient(command string, args []string) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Println("Error connecting to zeroctl daemon:", err)
		return
	}
	defer conn.Close()

	cmdMsg := CommandMessage{
		Command: command,
		Args:    args,
	}

	payload, err := json.Marshal(cmdMsg)
	if err != nil {
		fmt.Println("Error encoding command:", err)
		return
	}

	_, err = conn.Write(payload)
	if err != nil {
		fmt.Println("Error sending command:", err)
		return
	}

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	response := string(buffer[:n])
	if len(response) > 0 && !strings.HasSuffix(response, "\n") {
		response += "\n"
	}

	fmt.Print(response)
}
