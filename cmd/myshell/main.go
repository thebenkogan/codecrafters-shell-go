package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Print("$ ")

	// Wait for user input
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
}
