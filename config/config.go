package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func LoadEnv() {
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "development")
	}

	if err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("GO_ENV"))); err != nil {
		panic(err)
	}
}

func GetEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}

	return v
}
