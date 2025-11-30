package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	dm "turbobloom/distribute/models" // distribute models
	sm "turbobloom/shared"            // shared models
)

var keyMutex sync.Mutex

var filename *string

func readParams() dm.Parameters {
	data, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatal("Failed to read parameters file")
	}

	var config dm.Parameters
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Failed to parse parameters file")
	}

	return config
}

func updateCount(params dm.Parameters) {
	data, err := json.Marshal(params)
	if err != nil {
		fmt.Printf("Error marshaling parameters: %v\n", err)
		return
	}

	err = os.WriteFile(*filename, data, 0600)
	if err != nil {
		// maybe delete the file and fatally exit the server
		fmt.Printf("Error writing file: %v\n", err)
		return
	}
}

func calculateAColumn(ind uint32, config dm.Parameters) []uint64 {
	a_col := make([]uint64, config.Lambda+1)
	for i := uint32(0); i <= config.Lambda; i++ {
		a_col[i] = 0
		for k := uint32(0); k <= config.Lambda; k++ {
			a_col[i] += (config.D[i][k] * config.G[ind][k]) % config.Q
		}
		a_col[i] %= config.Q
	}
	return a_col
}

var ErrKeysExhausted = errors.New("keys are exhausted")

func getNextKey() (uint32, []uint64, []uint64, error) {
	keyMutex.Lock()
	defer keyMutex.Unlock()

	var config dm.Parameters = readParams()

	if config.Count == config.N {
		// we might want to delete the file and exit the server
		return 0, []uint64{}, []uint64{}, ErrKeysExhausted
	}

	id := config.Count
	config.Count += 1
	// I want this to be executed in the current lock
	updateCount(config)
	g_col := config.G[id]
	a_col := calculateAColumn(id, config)

	return id, g_col, a_col, nil
}

func distributeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received key distribution request")

	nodeId, g_col, a_col, err := getNextKey()
	if errors.Is(err, ErrKeysExhausted) {
		http.Error(w, err.Error(), http.StatusTeapot) // whatever, can do 404 also
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	response := sm.DistributeResponse{
		Id:   nodeId,
		Gcol: g_col,
		Acol: a_col,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	filename = flag.String("config", "./parameters/parameters.json", "Relative path to parameters file")
	flag.Parse()

	readParams()

	http.HandleFunc("/distribute", distributeHandler)

	port := ":8080"
	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Could not start server: ", err)
	}
}
