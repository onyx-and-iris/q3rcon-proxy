package main

import (
	"os"
	"strconv"
)

func getEnvInt(key string) (int, error) {
	s := os.Getenv(key)
	if s == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}
