package tests

import (
	"fmt"
	"os"
	"os/exec"
)

func startDockerCompose() {
	cmd := exec.Command("./setup.sh")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error starting docker-compose")
		fmt.Println(string(output))
		os.Exit(1)
	}

	fmt.Println("Docker-compose started successfully")
}

func stopDockerCompose() {
	cmd := exec.Command("./teardown.sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error stopping docker-compose")
		fmt.Println(string(output))
		os.Exit(1)
	}

	fmt.Println("Docker-compose stopped successfully")
}
