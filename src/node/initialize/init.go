package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"net/http"
	distribute_models "turbobloom/distribute/models"
	shared_types "turbobloom/shared/types"
)

func getParameters(serverUrl string) shared_types.DistributeResponse {
	response, err := http.Get(serverUrl)
	if err != nil {
		log.Fatalf("Error getting parameters: %v", err)
	}

	defer response.Body.Close()

	var distributeResponse shared_types.DistributeResponse
	err = json.NewDecoder(response.Body).Decode(&distributeResponse)
	if err != nil {
		log.Fatalf("Error parsing parameters: %v", err)
	}

	return distributeResponse
}


func saveParameters(parameters shared_types.DistributeResponse, filepath string) {
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

	distributeUrl := readStringWithDefault("Enter distribute server URL", "http://localhost:8080/distribute")
	filePath := readStringWithDefault("Enter relative path to store parameters file", "./parameters/parameters.json")

	if _, err := os.Stat(filepath); err == nil {
		if !readYesNo("File already exists. Overwrite?") {
			log.Fatal("File not saved.")
		}
	}

	distributeResponse := getParameters(distributeUrl)
	saveParameters(distributeResponse, filePath)

	fmt.Printf("Parameters saved to %s\n", filepath)
}
