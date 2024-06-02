package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	for {
		fmt.Print("$ ")

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}
		command = strings.TrimSpace(command)

		handleCommand(command)
	}
}

func handleCommand(command string) {
	parts := strings.Split(command, " ")
	switch parts[0] {
	case "exit":
		os.Exit(0)
	default:
		fmt.Printf("%s: command not found\n", command)
	}
}
