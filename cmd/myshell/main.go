package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
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

var BUILTINS = []string{"exit", "echo", "type", "pwd"}

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
		commandPath := locateCommand(parts[1], strings.Split(path, ":"))

		if commandPath != "" {
			fmt.Printf("%s is %s\n", parts[1], commandPath)
		} else {
			fmt.Printf("%s not found\n", parts[1])
		}
	case "pwd":
		pwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting working directory: %w", err)
		}
		fmt.Println(pwd)
	case "cd":
		if strings.HasPrefix(parts[1], "~") {
			usr, _ := user.Current()
			parts[1] = strings.Replace(parts[1], "~", usr.HomeDir, 1)
		}
		if err := os.Chdir(parts[1]); err != nil {
			fmt.Printf("%s: No such file or directory\n", parts[1])
		}
	default:
		path, _ := os.LookupEnv("PATH")
		commandPath := locateCommand(parts[0], strings.Split(path, ":"))
		if commandPath == "" {
			fmt.Printf("%s: command not found\n", command)
			return nil
		}

		cmd := exec.Command(commandPath, parts[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
		var exitErr *exec.ExitError
		if err := cmd.Run(); err != nil && !errors.As(err, &exitErr) {
			return fmt.Errorf("error executing external command: %w", err)
		}
	}
	return nil
}

func locateCommand(command string, path []string) string {
	// otherwise, try each directory in the path
	for _, dir := range path {
		fullpath := filepath.Join(dir, command)
		if f, err := os.Stat(fullpath); err == nil && f.Mode()&0111 != 0 {
			return fullpath
		}
	}

	// command might be a path to an executable, let's check
	if f, err := os.Stat(command); err == nil && f.Mode()&0111 != 0 {
		return command
	}

	return ""
}
