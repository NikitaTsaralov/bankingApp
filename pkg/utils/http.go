package utils

import (
	"fmt"
	"math/rand"
)

// Get config path for local or docker
func GetConfigPath(filename string, postfix string) string {
	if postfix == "docker" {
		return fmt.Sprintf("./config/%s-docker", filename)
	}
	return fmt.Sprintf("./config/%s-local", filename)
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
