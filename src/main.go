package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

var (
	last_working_directory = getUserHomeDir()
	history_absolute_path = getUserHomeDir() + "/.gosh_history"
	command_index = 0
) 

func main() {
	reader := bufio.NewReader(os.Stdin)

	history_file, err := os.OpenFile(history_absolute_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Could not open gosh history file.")
		return
	}

	defer history_file.Close()

	u, err := user.Current()
	if err != nil {
		log.Fatal("User not found!")
		return
	}

	for {
		wd := getWorkingDirectory()
		fmt.Printf("%s on %s\n$ ", u.Username, wd)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if err = handleInput(input, history_file); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

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

func handleInput(input string, history_file *os.File) error {
	input = strings.TrimSpace(input)
        
	_, err := history_file.WriteString(input + "\n")
	if err != nil {
		log.Fatal("Could not write on history file.")
		return nil
	}

	args := strings.Split(input, " ")

	if len(args) >= 2 {
		if strings.Contains(args[1], "~") {
			args[1] = strings.Replace(args[1], "~", getUserHomeDir(), 1)
		}
	}

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
