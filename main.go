package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Maxtors/surisoc"
)

// Global variables for the Suricata Socket Application
var (
	socketPath  string
	interactive bool
	session     *surisoc.SuricataSocket
)

// Initializer function to set up commandline arguments and socket session
func init() {
	var err error

	// Parse commandline arguments
	flag.StringVar(&socketPath, "socket", "/var/run/suricata/suricata-command.socket", "Full path to the suricata unix socket")
	flag.BoolVar(&interactive, "interactive", false, "Opens an interactive session to send commands to the socket")
	flag.Parse()

	// If the user wants to start an interactive session but has created / sent
	// command arguments at call
	if interactive && len(flag.Args()) > 0 {
		log.Fatalf("When running in interactive mode, do not supply command arguments: %+v\n", flag.Args())
	}

	// If the user does not want to start an interactive session but has not
	// created / sent any commandline arguments
	if !interactive && len(flag.Args()) < 1 {
		log.Fatalf("When running in interactive mode, supply atleast one command argument")
	}

	// Create a new Suricata Socket session
	session, err = surisoc.NewSuricataSocket(socketPath)
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
}

func main() {
	defer session.Close()

	// If we are not starting an interactive session just send the command
	if !interactive {
		sendCommandLine(flag.Args())
	} else {

		// Tell the user that an interactive session has started, and
		// display all the available commands
		fmt.Println(">> Entering Interactive Mode <<")
		fmt.Println(">> Valid Commands:")
		sendCommandLine([]string{"command-list"})

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print(">> ")
			text, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error: %s\n", err.Error())
			}
			sendCommandLine(strings.Split(strings.TrimSuffix(text, "\n"), " "))
		}
	}
}

// sendCommandLine will take a slice of commands and send them to the socket
func sendCommandLine(commands []string) {
	var err error
	var response *surisoc.SocketResponse
	var command string
	var arguments []string

	// Check if we should parse and send arguments or not
	if len(commands) > 1 {
		command = commands[0]
		arguments = commands[1:]
		response, err = session.Send(command, arguments...)
	} else {
		command = commands[0]
		response, err = session.Send(command)
	}

	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}

	// Get the string representation of the response message
	res, err := response.ToString()
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
	fmt.Println(res)
}
