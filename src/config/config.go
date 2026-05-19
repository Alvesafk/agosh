package config

import (
	"log"
	"os"
)

const (
	RED     = "\033[31m"
	GREEN   = "\033[32m"
	YELLOW  = "\033[33m"
	BLUE    = "\033[34m"
	MAGENTA = "\033[35m"
	CYAN    = "\033[36m"
	GRAY    = "\033[37m"
	WHITE   = "\033[97m"
)

var (
	config_absolute_path = getUserHomeDir() + ".agoshconfig"
)

type User_Config struct {
	First_Color, Second_Color string
	Base_Wd_Depth             int
}

func GetUserConfig() User_Config {
	return User_Config{
		First_Color:   CYAN,
		Second_Color:  BLUE,
		Base_Wd_Depth: 2,
	}
}

func getConfigFile() *os.File{
	config_file, err := os.OpenFile(config_absolute_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Could not open gosh history file.")
		return nil
	}

	return config_file
}

func parseConfig() {
	config_file := getConfigFile()

	defer config_file.Close()
}

func getUserHomeDir() string {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return home_dir
}
