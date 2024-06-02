package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	pwd, _ := os.Getwd()
	s := &shell{pwd}

	if err := s.run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type shell struct {
	pwd string // absolute
}

func (s *shell) run() error {
	for {
		fmt.Print("$ ")

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}
		command = strings.TrimSpace(command)

		if err := s.handleCommand(command); err != nil {
			return fmt.Errorf("error handling command: %w", err)
		}
	}
}

var BUILTINS = []string{"exit", "echo", "type", "pwd"}

func (s *shell) handleCommand(command string) error {
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
		fmt.Println(s.pwd)
	case "cd":
		if f, err := os.Stat(parts[1]); err == nil && f.IsDir() {
			s.pwd = parts[1]
		} else {
			fmt.Printf("cd: %s: No such file or directory\n", parts[1])
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
		cmd.Dir = s.pwd
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
