package env

import (
	"os"
	"strconv"
)

func GetInt(name string, def int) int {
	var env_str string = os.Getenv(name)
	if env_str == "" {
		return def
	}
	rv, err := strconv.Atoi(env_str)
	if err != nil {
		return def
	}
	return rv
}

func GetStr(name string, def string) string {
	env_str, exists := os.LookupEnv(name)
	if exists {
		return env_str
	}
	return def
}
