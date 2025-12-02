package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	shared "turbobloom/shared"
)

func getParameters(serverUrl string) shared.DistributeResponse {
	response, err := http.Get(serverUrl)
	if err != nil {
		log.Fatalf("Error getting parameters: %v", err)
	}

	defer response.Body.Close()

	var distributeResponse shared.DistributeResponse
	err = json.NewDecoder(response.Body).Decode(&distributeResponse)
	if err != nil {
		log.Fatalf("Error parsing parameters: %v", err)
	}

	return distributeResponse
}

func saveParameters(parameters shared.DistributeResponse, filepath string) {
	data, err := json.Marshal(parameters)
	if err != nil {
		log.Fatalf("Error marshalling parameters: %v", err)
	}

	err = os.WriteFile(filepath, data, 0600)
	if err != nil {
		log.Fatalf("Error writing parameters file: %v", err)
	}
}

func main() {
	fmt.Println("Enter parameters:")

	distributeUrl := shared.ReadStringWithDefault("Enter distribute server URL", "http://localhost:8080/distribute")
	filePath := shared.ReadStringWithDefault("Enter relative path to store parameters file", "./parameters/parameters.json")

	if _, err := os.Stat(filePath); err == nil {
		if !shared.ReadYesNo("File already exists. Overwrite?") {
			log.Fatal("File not saved.")
		}
	}

	distributeResponse := getParameters(distributeUrl)
	saveParameters(distributeResponse, filePath)

	fmt.Printf("Parameters saved to %s\n", filePath)
}
