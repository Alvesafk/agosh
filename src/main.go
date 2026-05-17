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

const (
	BASE_WD_DEPTH = 2
)

var (
	last_working_directory = getUserHomeDir()
	history_absolute_path  = getUserHomeDir() + "/.gosh_history"
	command_index          = 0
	buf_position           = 0
)

func main() {

	history_file, err := os.OpenFile(history_absolute_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Could not open gosh history file.")
		return
	}

	defer history_file.Close()

	u := getUser()

	for {
		wd := getFormattedWorkingDirectory(BASE_WD_DEPTH)
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

func getUser() *user.User {
	u, err := user.Current()
	if err != nil {
		log.Fatal("User not found!")
		return nil
	}

	return u
}

func getFormattedWorkingDirectory(depth int) string {
	wd_array := reverseStringArray([]string{getWorkingDirectory()})
	wd_formatted := strings.Split(wd_array[0], "/")

	if depth < 1 || len(wd_formatted) <= 2 {
		return getWorkingDirectory()
	}

	var result string

	for i := depth; i > 0; i-- {
		result += wd_formatted[len(wd_formatted)-i] + "/"
	}

	return result
}

func clearStdin() {
	fmt.Print("\r\033[K$ ")
}

func reverseStringArray(input []string) []string {
	result := input
	for i, s := 0, len(input)-1; i < s; i, s = i+1, s-1 {
		result[i], result[s] = result[s], result[i]
	}

	return result
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

	if len(args) > 1 {
		for i, _ := range args {
			args[i] = strings.Replace(args[i], "~", getUserHomeDir(), 1)
		}
	}

	switch args[0] {
	case "cd":

		if len(args) < 2 {
			last_working_directory = getWorkingDirectory()
			return os.Chdir(getUserHomeDir())
		}

		if args[1] == "-" {
			fmt.Println(last_working_directory)
			wd_to_go := last_working_directory
			last_working_directory = getWorkingDirectory()
			return os.Chdir(wd_to_go)
		}

		last_working_directory = getWorkingDirectory()

		return os.Chdir(args[1])

	case "exit":
		os.Exit(0)

	case "getHost", "gethost":
		fmt.Printf("%s\n", getHostName())
		return nil
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func readInput() (string, error) {
	var (
		buf      []rune
		result   string
		inputErr error
	)

	data_raw, err := os.ReadFile(history_absolute_path)
	if err != nil {
		log.Fatal("Could not read from the history file")
		return "", err
	}

	data := reverseStringArray(strings.Split(string(data_raw), "\n"))

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
			if len(buf) > 0 || buf_position > 0 {
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
			clearStdin()
			command_index++
			buf = []rune(data[command_index])
			fmt.Print(data[command_index])

		case keys.Down:
			clearStdin()
			if command_index > 0 {
				command_index--
				buf = []rune(data[command_index])
				fmt.Print(data[command_index])
			}

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
