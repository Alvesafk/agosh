package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func getUserHomeDir() string {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return home_dir
}

func getHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}

	return hostname
}

func handleInput(input string) error {
	input = strings.TrimSuffix(input, "\n")

	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		if args[1] == "~" {
			return os.Chdir(getUserHomeDir())
		}

		if len(args) < 2 {
			return errors.New("Path is required!")
		}

		return os.Chdir(args[1])

	case "exit":
		os.Exit(0)
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	u, err := user.Current()
	if err != nil {
		fmt.Println("Error: user not found!")
		return
	}

	for {
		fmt.Printf("%s on ionknow\n$ ", u.Username)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if err = handleInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
