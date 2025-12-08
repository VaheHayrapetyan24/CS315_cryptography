package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"turbobloom/node/nodeUtils"
	"turbobloom/shared"
)

type MessagePayload struct {
	Gcol    []uint64 `json:"g_col"`
	Message string   `json:"message"`
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
	log.Printf("Encrypted message content: %s", payload.Message)

	k := nodeUtils.GetKey(payload.Gcol, nodeConfig.Acol, nodeConfig.Q)
	log.Printf("Calculated key: %v", k)

	plaintext, err := nodeUtils.DecryptStringAES128CBC(k, payload.Message)
	if err != nil {
		log.Printf("Failed to decrypt message: %v", err)
		http.Error(w, "Failed to decrypt message", http.StatusBadRequest)
		return
	}

	log.Printf("Decrypted message content: %s", plaintext)

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
	filename := flag.String("config", "./node/parameters/parameters.json", "Relative path to parameters file")
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
