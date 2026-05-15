package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

var last_working_directory = getUserHomeDir()

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

func getWorkingDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return wd
}

func handleInput(input string) error {
	input = strings.TrimSpace(input)

	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":

		if len(args) < 2 {
			last_working_directory = getWorkingDirectory()
			return os.Chdir(getUserHomeDir())
		}

		if args[1] == "~" {
			last_working_directory = getWorkingDirectory()
			return os.Chdir(getUserHomeDir())
		}

		if args[1] == "-" {
			fmt.Println(last_working_directory)
			return os.Chdir(last_working_directory)
		}

		last_working_directory = getWorkingDirectory()

		return os.Chdir(args[1])

	case "exit":
		os.Exit(0)

	case "":
		return nil
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
		wd := getWorkingDirectory()
		fmt.Printf("%s on %s\n$ ", u.Username, wd)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if err = handleInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
