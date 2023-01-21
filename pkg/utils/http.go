package utils

import "fmt"

// Get config path for local or docker
func GetConfigPath(filename string, postfix string) string {
	if postfix == "docker" {
		return fmt.Sprintf("./config/%s-docker", filename)
	}
	return fmt.Sprintf("./config/%s-local", filename)
}
