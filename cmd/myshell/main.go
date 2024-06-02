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

func handleCommand(command string) error {
	parts := strings.Split(command, " ")
	switch parts[0] {
	case "exit":
		os.Exit(0)
	case "echo":
		fmt.Println(strings.Join(parts[1:], " "))
	case "type":
		handleType(parts)
	case "pwd":
		return handlePwd()
	case "cd":
		handleCd(parts)
	default:
		return handleExternal(parts)
	}
	return nil
}

func handleExternal(command []string) error {
	path := os.Getenv("PATH")
	commandPath := locateCommand(command[0], strings.Split(path, ":"))
	if commandPath == "" {
		fmt.Printf("%s: command not found\n", command[0])
		return nil
	}

	cmd := exec.Command(commandPath, command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	var exitErr *exec.ExitError
	if err := cmd.Run(); err != nil && !errors.As(err, &exitErr) {
		return fmt.Errorf("error executing external command: %w", err)
	}
	return nil
}

var BUILTINS = []string{"exit", "echo", "type", "pwd", "cd"}

func handleType(command []string) {
	if slices.Contains(BUILTINS, command[1]) {
		fmt.Printf("%s is a shell builtin\n", command[1])
		return
	}

	path, _ := os.LookupEnv("PATH")
	commandPath := locateCommand(command[1], strings.Split(path, ":"))

	if commandPath != "" {
		fmt.Printf("%s is %s\n", command[1], commandPath)
	} else {
		fmt.Printf("%s not found\n", command[1])
	}
}

func handlePwd() error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %w", err)
	}
	fmt.Println(pwd)
	return nil
}

func handleCd(command []string) {
	if strings.HasPrefix(command[1], "~") {
		home := os.Getenv("HOME")
		command[1] = strings.Replace(command[1], "~", home, 1)
	}
	if err := os.Chdir(command[1]); err != nil {
		fmt.Printf("%s: No such file or directory\n", command[1])
	}
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
