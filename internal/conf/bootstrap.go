package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/go-kratos/kratos/v2/log"
	"gopkg.in/yaml.v3"
)

func LoadConfig(configPath string, logger log.Logger) (*Bootstrap, error) {
	helper := log.NewHelper(logger)
	
	
	var confFile string
	stat, err := os.Stat(configPath)
	if err == nil && stat.IsDir() {
		confFile = filepath.Join(configPath, "config.yaml")
	} else {
		confFile = configPath
	}
	
	helper.Infof("Loading configuration from: %s", confFile)
	
	
	data, err := os.ReadFile(confFile)
	if err != nil {
		helper.Errorf("Failed to read configuration file: %v", err)
		return nil, err
	}
	
	
	var bootstrap Bootstrap
	if err := yaml.Unmarshal(data, &bootstrap); err != nil {
		helper.Errorf("Failed to parse configuration: %v", err)
		return nil, err
	}
	
	return &bootstrap, nil
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

func GetGithubConfig(configRoot string, logger log.Logger) (*Github, error) {
	helper := log.NewHelper(logger)
	secretPath := filepath.Join(configRoot, "secrets", "github.json")
	
	var githubConfig Github
	err := GetSecret(secretPath, &githubConfig)
	if err != nil {
		helper.Warnf("Failed to read github secret file: %v", err)
		return &Github{}, err
	}
	
	return &githubConfig, nil
}

func GetDatabaseConfig(configRoot string, logger log.Logger) (*Data_Database, error) {
	helper := log.NewHelper(logger)
	secretPath := filepath.Join(configRoot, "secrets", "database.json")
	
	var dbConfig Data_Database
	err := GetSecret(secretPath, &dbConfig)
	if err != nil {
		helper.Warnf("Failed to read database secret file: %v", err)
		return &Data_Database{}, err
	}
	
	return &dbConfig, nil
}

func GetRedisConfig(configRoot string, logger log.Logger) (*Data_Redis, error) {
	helper := log.NewHelper(logger)
	secretPath := filepath.Join(configRoot, "secrets", "redis.json")
	
	var redisConfig Data_Redis
	err := GetSecret(secretPath, &redisConfig)
	if err != nil {
		helper.Warnf("Failed to read redis secret file: %v", err)
		return &Data_Redis{}, err
	}
	
	return &redisConfig, nil
}


func LoadSecrets(bootstrap *Bootstrap, configRoot string, logger log.Logger) (*Bootstrap, error) {
	helper := log.NewHelper(logger)
	
	
	githubConfig, err := GetGithubConfig(configRoot, logger)
	if err == nil && githubConfig.Token != "" {
		bootstrap.Github = githubConfig
		helper.Info("Loaded GitHub configuration from secrets")
	}
	
	
	dbConfig, err := GetDatabaseConfig(configRoot, logger)
	if err == nil && dbConfig.Driver != "" {
		if bootstrap.Data == nil {
			bootstrap.Data = &Data{}
		}
		bootstrap.Data.Database = dbConfig
		helper.Info("Loaded Database configuration from secrets")
	}
	
	
	redisConfig, err := GetRedisConfig(configRoot, logger)
	if err == nil && redisConfig.Addr != "" {
		if bootstrap.Data == nil {
			bootstrap.Data = &Data{}
		}
		bootstrap.Data.Redis = redisConfig
		helper.Info("Loaded Redis configuration from secrets")
	}
	
	return bootstrap, nil
} 