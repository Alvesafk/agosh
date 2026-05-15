package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

var (
	last_working_directory = getUserHomeDir()
	history_absolute_path = getUserHomeDir() + "/.gosh_history"
	command_index = 0
	buf_position = 0
) 

func main() {

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

		input, err := readInput()
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
        
	if input == "" {
		return nil
	}

	_, err := history_file.WriteString(input + "\n")
	if err != nil {
		log.Fatal("Could not write on history file.")
		return nil
	}

	args := strings.Split(input, " ")

	if len(args) >= 2 {
		args[1] = strings.Replace(args[1], "~", getUserHomeDir(), 1)
	}

	switch args[0] {
	case "cd":

		if len(args) < 2 {
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
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func readInput() (string, error) {
	var (
		buf []rune
		result string
		inputErr error
	)

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
			case keys.CtrlC, keys.CtrlD:
				inputErr = io.EOF
				buf_position = 0
				return true, nil

			case keys.Enter:
				fmt.Print("\r\n")
				result = string(buf) + "\n"
				buf_position = 0
				return true, nil

			case keys.Backspace:
				if len(buf) > 0  || buf_position > 0 {
					buf = buf[:len(buf)-1]
					fmt.Printf("\b \b")
					buf_position--
				}

			case keys.Right:
				if buf_position <= len(buf) {
					fmt.Printf("\033[%dC", 1)
					buf_position++
				}

			case keys.Left:
				if len(buf) > 0 || buf_position > 0 {
					fmt.Printf("\033[%dD", 1)
					buf_position--
				}

			case keys.Up:
				data, err := os.ReadFile(history_absolute_path)
				if err != nil {
					log.Fatal("Could not read from the history file")		
					return  true, nil
				}

				fmt.Println(string(data[:]))

			default:
				if key.Code == keys.RuneKey || key.Code == keys.Space {
					buf = append(buf, key.Runes[0])
					fmt.Print(string(key.Runes[0]))
					buf_position++
				}
		}

		return false, nil
	})

	return result, inputErr
}
