package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Maxtors/surisoc"
)

// Global variables for the Suricata Socket Application
var (
	socketPath  string
	interactive bool
	session     *surisoc.SuricataSocket
	signals     chan os.Signal
	done        chan bool
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

	// Create channels
	signals = make(chan os.Signal, 1)
	done = make(chan bool, 1)

	// Set up the signals to listen to, SIGQUIT is used internaly to stop
	// the current process
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
}

func main() {
	defer session.Close()

	// Go-Routine to handle reciving of signals
	go func() {
		signal := <-signals
		log.Printf("Recived signal: %s\n", signal)
		done <- true
	}()

	// Go-Routine to handle the actuall logic of this program
	go func() {
		// If we are not starting an interactive session just send the command
		if !interactive {
			sendCommandLine(flag.Args())
			done <- true
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
					log.Printf("Error: %s\n", err.Error())
					signals <- syscall.SIGQUIT
				}
				sendCommandLine(strings.Split(strings.TrimSuffix(text, "\n"), " "))
			}
		}
	}()

	// Wait for a done signal
	<-done
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
