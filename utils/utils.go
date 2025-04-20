package utils

import (
	"luminex-service/constants"
	"os"
)

func GetEnvironment() string {
	env := os.Getenv(constants.Env)
	if env == "" {
		return env
	}
	return constants.LocalEnv
}

func GetConfig() string {
	env := GetEnvironment()
	return constants.ConfigFileMap[env]
}
