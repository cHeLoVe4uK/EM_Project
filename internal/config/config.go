package config

import "github.com/joho/godotenv"

func Load() error {

	_ = godotenv.Load()

	return nil
}
