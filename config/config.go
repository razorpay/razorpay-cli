package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configDir  = ".razorpay"
	configFile = "config"
	configType = "yaml"
)

func ConfigFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, configDir, configFile+"."+configType)
}

func Init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error finding home directory: %v\n", err)
		os.Exit(1)
	}

	viper.SetConfigName(configFile)
	viper.SetConfigType(configType)
	viper.AddConfigPath(filepath.Join(home, configDir))
	viper.AutomaticEnv()

	// env overrides
	viper.SetEnvPrefix("RAZORPAY")
	_ = viper.BindEnv("key_id", "RAZORPAY_KEY_ID")
	_ = viper.BindEnv("key_secret", "RAZORPAY_KEY_SECRET")

	_ = viper.ReadInConfig()
}

func KeyID() string {
	return viper.GetString("key_id")
}

func KeySecret() string {
	return viper.GetString("key_secret")
}

func Save(keyID, keySecret string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, configDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	viper.Set("key_id", keyID)
	viper.Set("key_secret", keySecret)
	return viper.WriteConfigAs(filepath.Join(dir, configFile+"."+configType))
}
