package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Port        string
	AdminKey    string
	DataRoot    string
	PostsDir    string
	UploadsDir  string
	DataDir     string
	DBFile      string
}

func New() *Config {
	dataRoot := getEnv("DATA_ROOT", ".")

	return &Config{
		Port:       getEnv("PORT", "3000"),
		AdminKey:   getEnv("ADMIN_KEY", "changeme"),
		DataRoot:   dataRoot,
		PostsDir:   filepath.Join(dataRoot, "posts"),
		UploadsDir: filepath.Join(dataRoot, "uploads"),
		DataDir:    filepath.Join(dataRoot, "data"),
		DBFile:     filepath.Join(dataRoot, "data", "posts.json"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
