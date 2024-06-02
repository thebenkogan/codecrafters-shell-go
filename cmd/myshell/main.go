package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

		if err := handleCommand(command); err != nil {
			return fmt.Errorf("error handling command: %w", err)
		}
	}
}

var BUILTINS = []string{"exit", "echo", "type"}

func handleCommand(command string) error {
	parts := strings.Split(command, " ")
	switch parts[0] {
	case "exit":
		os.Exit(0)
	case "echo":
		fmt.Println(strings.Join(parts[1:], " "))
	case "type":
		if slices.Contains(BUILTINS, parts[1]) {
			fmt.Printf("%s is a shell builtin\n", parts[1])
			return nil
		}

		path, _ := os.LookupEnv("PATH")
		commandPath, err := locateCommand(parts[1], strings.Split(path, ":"))
		if err != nil {
			return fmt.Errorf("error locating command: %w", err)
		}

		if commandPath != "" {
			fmt.Printf("%s is %s\n", parts[1], commandPath)
		} else {
			fmt.Printf("%s: command not found\n", parts[1])
		}
	default:
		fmt.Printf("%s: command not found\n", command)
	}
	return nil
}

func locateCommand(command string, path []string) (string, error) {
	for _, dir := range path {
		var found string
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if !d.IsDir() && d.Name() == command {
				found = path
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			return "", fmt.Errorf("error walking directory: %w", err)
		}
		if found != "" {
			return found, nil
		}
	}
	return "", nil
}
