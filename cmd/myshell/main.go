package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

var BUILTINS = []string{"exit", "echo", "type"}

func handleCommand(command string) {
	parts := strings.Split(command, " ")
	switch parts[0] {
	case "exit":
		os.Exit(0)
	case "echo":
		fmt.Println(strings.Join(parts[1:], " "))
	case "type":
		if slices.Contains(BUILTINS, parts[1]) {
			fmt.Printf("%s is a shell builtin\n", parts[1])
		} else {
			fmt.Printf("%s not found\n", parts[1])
		}
	default:
		fmt.Printf("%s: command not found\n", command)
	}
}
