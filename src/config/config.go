package config

const (
	RED = "\033[31m"
	GREEN = "\033[32m"
	YELLOW = "\033[33m"
	BLUE = "\033[34m"
	MAGENTA = "\033[35m"
	CYAN = "\033[36m"
	GRAY = "\033[37m"
	WHITE = "\033[97m"
)

type User_Config struct {
	First_Color string
	Second_Color string
}

func GetUserConfig() User_Config {
	return User_Config{
		First_Color: CYAN,
		Second_Color: BLUE,
	}
}
