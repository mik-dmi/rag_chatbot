package env

import (
	"fmt"
	"os"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		fmt.Print("Problem i  lockup")
		return fallback
	}

	return val
}
