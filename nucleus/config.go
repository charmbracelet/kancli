package nucleus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	NucleusDirName  = ".kancli-nucleus" // dotfile style for home dir
	NucleusFileName = "kancli.json"
)

var (
	ConfigDir  string
	ConfigPath string
)

func init() {
	initConfigPaths()
	Init()
}

func initConfigPaths() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("❌ Error: Unable to determine user home directory:", err)
		os.Exit(1)
	}

	ConfigDir = filepath.Join(homeDir, NucleusDirName)
	ConfigPath = filepath.Join(ConfigDir, NucleusFileName)
}

func Init() {
	if !NucleusExists() {
		err := createDefaultConfig()
		if err != nil {
			fmt.Println("❌ Failed to create default config:", err)
			os.Exit(1)
		}
	} else {
		loadExistingConfig()
	}
}

func NucleusExists() bool {
	_, err := os.Stat(ConfigPath)
	return !os.IsNotExist(err)
}

func GetConfigDir() string {
	return ConfigDir
}

// --- Default config setup ---

type KancliConfig struct {
	Version string `json:"version"`
	User    string `json:"user"`
	// Add more config fields as needed
}

func createDefaultConfig() error {
	if err := os.MkdirAll(ConfigDir, 0755); err != nil {
		return fmt.Errorf("unable to create config directory: %w", err)
	}

	defaultConfig := KancliConfig{
		Version: "1.0.0",
		User:    "default",
	}

	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	err = os.WriteFile(ConfigPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Println("✅ Nucleus created at", ConfigPath)
	return nil
}

func loadExistingConfig() {
	// Placeholder: implement loading logic here if needed
	fmt.Println("🔁 Loading existing config from", ConfigPath)
}
