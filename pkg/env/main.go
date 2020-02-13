package env

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Init imports ENV vars
func Init() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := os.Stat(fmt.Sprintf("%s/%s", dir, ".env")); err == nil {
		fmt.Println(".env file found, loading")

		err := godotenv.Load()
		if err != nil {
			return err
		}
	} else {
		fmt.Println("No .env file found")
	}

	return nil
}

// MustGetString ...
func MustGetString(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Panicf("%s environment variable not set.", key)
	}

	return v
}

// MayGetString ...
func MayGetString(key string) string {
	v := os.Getenv(key)
	if v == "" {
		return ""
	}

	return v
}

// MustGetInt ...
func MustGetInt(key string) int {
	v := os.Getenv(key)
	if v == "" {
		log.Panicf("%s environment variable not set.", key)
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		log.Panicf("Unable to parse %s, err: %s", key, err.Error())
	}

	return i
}

// MustPresent ...
func MustPresent(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		return false
	}

	return true
}
