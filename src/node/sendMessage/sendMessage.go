package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"turbobloom/node/nodeUtils"
	"turbobloom/shared"
)

type GcolResponse struct {
	Gcol []uint64 `json:"g_col"`
}

type MessagePayload struct {
	Gcol    []uint64 `json:"g_col"`
	Message string   `json:"message"`
}

func readParams(filepath string) shared.DistributeResponse {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read parameters file: %v", err)
	}

	var distributeResponse shared.DistributeResponse
	err = json.Unmarshal(data, &distributeResponse)
	if err != nil {
		log.Fatalf("Failed to parse parameters file: %v", err)
	}

	return distributeResponse
}

func fetchRemoteGcol(baseURL string) ([]uint64, error) {
	gcolURL := strings.TrimSuffix(baseURL, "/") + "/gcol"

	response, err := http.Get(gcolURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Gcol: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("failed to fetch Gcol: status %d, body: %s", response.StatusCode, string(body))
	}

	var gcolResp GcolResponse
	err = json.NewDecoder(response.Body).Decode(&gcolResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gcol response: %v", err)
	}

	return gcolResp.Gcol, nil
}

func sendMessage(baseURL string, payload MessagePayload) error {
	messageURL := strings.TrimSuffix(baseURL, "/") + "/message"

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	response, err := http.Post(messageURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("failed to send message: status %d, body: %s", response.StatusCode, string(body))
	}

	return nil
}

func main() {
	fmt.Println("Enter parameters:")

	configPath := shared.ReadStringWithDefault("Enter config file path", "./node/parameters/parameters.json")
	targetURL := shared.ReadStringWithDefault("Enter target node URL", "http://localhost:8080")
	message := shared.ReadStringWithDefault("Enter message", "Hello!")

	// Validate target URL
	if _, err := url.Parse(targetURL); err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}

	// Make configPath absolute if it's relative
	if !filepath.IsAbs(configPath) {
		configPath, _ = filepath.Abs(configPath)
	}

	// Read local parameters
	fmt.Printf("\nReading local config from %s\n", configPath)
	myConfig := readParams(configPath)
	fmt.Printf("Local node ID: %d\n", myConfig.Id)

	// Fetch remote node's Gcol
	fmt.Printf("\nFetching Gcol from %s\n", targetURL)
	remoteGcol, err := fetchRemoteGcol(targetURL)
	if err != nil {
		log.Fatalf("Error fetching remote Gcol: %v", err)
	}
	fmt.Printf("Received remote Gcol of length %d: %v\n", len(remoteGcol), remoteGcol)

	// Calculate key using local Acol and remote Gcol
	k := nodeUtils.GetKey(remoteGcol, myConfig.Acol, myConfig.Q)
	fmt.Printf("\nCalculated key: %v\n", k)

	// Send message with local Gcol and plaintext message
	fmt.Printf("\nSending message to %s/message\n", targetURL)
	payload := MessagePayload{
		Gcol:    myConfig.Gcol,
		Message: message,
	}

	err = sendMessage(targetURL, payload)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	fmt.Println("Message sent successfully!")
}
