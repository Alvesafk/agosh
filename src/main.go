package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func handleInput(input string) error {
	input = strings.TrimSuffix(input, "\n")

	args := strings.Split(input, " ")

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func main() {
	fmt.Println("Hello, World!")

	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if err = handleInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
