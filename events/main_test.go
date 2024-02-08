package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func startDockerCompose() {
	cmd := exec.Command("./setup.sh")

	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error starting docker-compose")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Docker-compose started successfully")
}

func stopDockerCompose() {
	cmd := exec.Command("./teardown.sh")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error stopping docker-compose")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Docker-compose stopped successfully")
}

func TestSubscribe(t *testing.T) {
	startDockerCompose()
	defer stopDockerCompose()

	val, err := json.Marshal(fiber.Map{"host": "http://localhost:8081", "event_type": "comment"})
	if err != nil {
		t.Error(err)
	}

	// Test logic here
	req, err := http.NewRequest("POST", "http://localhost:8080/api/subscribe", bytes.NewBuffer(val))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(body))
	t.Log("Subscribed successfully")
}
