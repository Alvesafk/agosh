package config

import (
	"fmt"
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
	config_absolute_path = getUserHomeDir() + "/.agoshconfig"
)

type User_Config struct {
	First_Color, Second_Color string
	Base_Wd_Depth             int
}

func Test() {
	parseConfig()
}

func GetUserConfig() User_Config {
	return User_Config{
		First_Color:   CYAN,
		Second_Color:  BLUE,
		Base_Wd_Depth: 2,
	}
}

func getDataConfig() []byte {
	data_config, err := os.ReadFile(config_absolute_path)
	if err != nil {
		log.Print("Could not find from the config file, creating the file.")

		err := os.WriteFile(config_absolute_path, []byte("// agosh config file"), 0644)
		if err != nil {
			log.Print("Error during creation, trying again.")
		}
	}

	return data_config
}

func parseConfig() {
	data_config := getDataConfig()

	fmt.Println(string(data_config))
	fmt.Printf("%T\n", data_config)
}

func getUserHomeDir() string {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return home_dir
}
