package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"turbobloom/shared"
)

type MessagePayload struct {
	Gcol    []uint64 `json:"g_col"`
	Message string   `json:"message"`
}

func getKey(Gcol []uint64, Acol []uint64, q uint64) uint64 {
	var k uint64 = 0
	for i := 0; i < len(Gcol) && i < len(Acol); i++ {
		k += (Acol[i] * Gcol[i]) % q
	}
	return k % q
}

func messageHandler(w http.ResponseWriter, r *http.Request, nodeConfig shared.DistributeResponse) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload MessagePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		log.Printf("Failed to decode message payload: %v", err)
		return
	}

	if len(payload.Gcol) == 0 {
		http.Error(w, "Gcol is required", http.StatusBadRequest)
		return
	}

	log.Printf("Received message with Gcol of length %d: %v", len(payload.Gcol), payload.Gcol)
	log.Printf("Message content: %s", payload.Message)

	// Calculate the key using the received Gcol and the node's Acol
	k := getKey(payload.Gcol, nodeConfig.Acol, nodeConfig.Q)
	log.Printf("Calculated key: %v", k)

	// TODO: Process the message using the calculated key
	// For now, just acknowledge receipt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "received"})
}

func gcolHandler(w http.ResponseWriter, r *http.Request, nodeConfig shared.DistributeResponse) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"g_col": nodeConfig.Gcol,
	})
}

func handlerGenerator(handler func(http.ResponseWriter, *http.Request, shared.DistributeResponse), nodeConfig shared.DistributeResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, nodeConfig)
	}
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

func main() {
	portPtr := flag.Int("port", 8080, "Port to run the server on")
	filename := flag.String("config", "./parameters/parameters.json", "Relative path to parameters file")
	flag.Parse()

	nodeConfig := readParams(*filename)
	http.HandleFunc("/message", handlerGenerator(messageHandler, nodeConfig))
	http.HandleFunc("/gcol", handlerGenerator(gcolHandler, nodeConfig))

	port := fmt.Sprintf(":%d", *portPtr)
	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Could not start server: ", err)
	}
}
