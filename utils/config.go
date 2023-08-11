package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"jarvis-trading-bot/consts"

	"github.com/spf13/viper"
)

var configVars = []string{
	consts.AwsRegion,
	consts.AwsAccessKeyId,
	consts.AwsSecretAccessKey,
	consts.PulumiAccessToken,
	consts.ProjectDir,
	consts.ProjectOutputDir,
	consts.CloudfareApiToken,
	consts.DnsZone,
	consts.DnsRecord,
	consts.CloudfareApiKey,
	consts.CloudfareApiEmail,
}

func InitConfig() {
	viper.AddConfigPath(filepath.Join(GetExecutablePath(), "../"))
	viper.SetConfigType("env")
	viper.SetConfigName(".config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Can't config file: %v", err)
		os.Exit(1)
	}

	InitEnvVars()
}

func InitEnvVars() {
	for _, envVar := range configVars {
		if os.Getenv(envVar) == "" {
			os.Setenv(envVar, viper.GetString(envVar))
		}
	}
}

func GetConfigValue(key string) string {
	return viper.GetString(key)
}

func GetExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
