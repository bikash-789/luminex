package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"luminex-service/internal/interfaces/entity"
	"luminex-service/utils"
	"os"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	"github.com/go-kratos/kratos/v2/log"
)

func LoadEnvConfig(path *string) (*Bootstrap, log.Logger) {
	flag.Parse()
	var bootstrap Bootstrap
	fileSource := *path + utils.GetConfig()
	c := config.New(config.WithSource(file.NewSource(fileSource)))
	defer c.Close()

	// Load configuration
	if err := c.Load(); err != nil {
		panic(fmt.Errorf("failed to load config: %v", err))
	}

	// Scan configuration to bootstrap
	if err := c.Scan(&bootstrap); err != nil {
		panic(fmt.Errorf("failed to scan config: %v", err))
	}

	// Create logger
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)

	// Set log level from config
	logLevel := "info" // Default level
	if bootstrap.Logger != nil && bootstrap.Logger.Level != "" {
		logLevel = bootstrap.Logger.Level
	}

	logger = log.NewFilter(logger, log.FilterLevel(log.ParseLevel(logLevel)))
	log.SetLogger(logger)

	return &bootstrap, logger
}

func GetSecret(secretFilePath string, result interface{}) error {
	byteValue, err := os.ReadFile(secretFilePath)
	if err != nil {
		return fmt.Errorf("error reading secret file: %w", err)
	}

	err = json.Unmarshal(byteValue, result)
	if err != nil {
		return fmt.Errorf("error unmarshalling secret JSON: %w", err)
	}

	return nil
}

func GetGithubConfig(bootstrap *Bootstrap) entity.GithubConfig {
	fileLocation := bootstrap.Server.GetGithubSecretFileLocation()
	var githubConfig entity.GithubConfig
	if err := GetSecret(fileLocation, &githubConfig); err != nil {
		log.Fatalf("Error reading github secret file: %v", err)
		panic(err)
	}
	return githubConfig
}
