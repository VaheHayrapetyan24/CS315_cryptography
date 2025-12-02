package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
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

func generateGColumn(config dm.Parameters) []uint64 {
	rand.Seed(time.Now().UnixNano())
	g := make([]uint64, config.Lambda+1)
	len := config.Lambda + 1
	for i := 0; i < int(len); i++ {
		g[i] = uint64(rand.Intn(2))
	}

	return g
}

func calculateAColumn(g []uint64, config dm.Parameters) []uint64 {
	a_col := make([]uint64, config.Lambda+1)
	for i := uint32(0); i <= config.Lambda; i++ {
		a_col[i] = 0
		for k := uint32(0); k <= config.Lambda; k++ {
			a_col[i] += (config.D[i][k] * g[k]) % config.Q
		}
		a_col[i] %= config.Q
	}
	return a_col
}

// var ErrKeysExhausted = errors.New("keys are exhausted")

func getNextKey() (uint32, []uint64, []uint64, uint64) {
	keyMutex.Lock()
	defer keyMutex.Unlock()

	var config dm.Parameters = readParams()

	id := config.Count
	config.Count += 1
	// I want this to be executed in the current lock
	updateCount(config)
	g_col := generateGColumn(config)
	a_col := calculateAColumn(g_col, config)

	return id, g_col, a_col, config.Q
}

func distributeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received key distribution request")

	nodeId, g_col, a_col, q := getNextKey()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	response := sm.DistributeResponse{
		Id:   nodeId,
		Q:    q,
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
