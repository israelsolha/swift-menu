package config

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"swift-menu-session/utils"
)

func SetupDocker(env string) {
	composeFile := findComposePath(env)
	cmd := exec.Command("docker-compose", "-f", composeFile, "up", "-d")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to run docker compose %s", err)
	}
}

func TearDown(env string, dockerConfig Docker) {
	composeFile := findComposePath(env)
	args := []string{"-f", composeFile, "down"}
	args = append(args, dockerConfig.TearDownTags...)
	cmd := exec.Command("docker-compose", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to run docker compose %s", err)
	}
}

func findComposePath(env string) string {
	rootPath := utils.FindProjectRoot()
	fileName := fmt.Sprintf("docker-compose-%s.yml", env)
	return path.Join(rootPath, "resources", "dockerfiles", fileName)
}
