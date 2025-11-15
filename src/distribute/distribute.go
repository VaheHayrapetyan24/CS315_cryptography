package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
	"turbobloom/models"
)

var nodeCounter uint32 = 0
var keyMutex sync.Mutex
var paramsMutex sync.Mutex

const q uint64 = 19
const lambda uint32 = 2
const n uint32 = 5

var G_t = [n][lambda + 1]uint64{ // this is transpose of G
	{1, 2, 4},
	{1, 4, 16},
	{1, 8, 7},
	{1, 16, 9},
	{1, 13, 17},
}
var D = [lambda + 1][lambda + 1]uint64{
	{6, 15, 14},
	{15, 1, 8},
	{14, 8, 3},
}

var A [n][lambda + 1]uint64

func initialize(filename string) {
	// TODO: this should read the file and set q, n, lambda, d
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Failed to read parameters file")
	}

	var config Parameters
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Failed to parse parameters file")
	}

	log.Printf("%s", config)

	for i := uint32(0); i <= lambda; i++ {
		for j := uint32(0); j < n; j++ {
			A[j][i] = 0
			for k := uint32(0); k <= lambda; k++ {
				A[j][i] += (D[i][k] * G_t[j][k]) % q
			}
			A[j][i] %= q
		}
	}
}

var ErrKeysExhausted = errors.New("keys are exhausted")

func getNextKey() (uint32, []uint64, []uint64, error) {
	keyMutex.Lock()

	if nodeCounter == n {
		return 0, []uint64{}, []uint64{}, ErrKeysExhausted
	}

	id := nodeCounter
	nodeCounter++
	keyMutex.Unlock()
	return id, G_t[id][:], A[id][:], nil
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

	response := models.DistributeResponse{
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

// func parametersHandler(w http.ResponseWriter, r *http.Request) {
// 	log.Printf("Received parameters request")

// }

func main() {
	filename := flag.String("config", "./parameters/parameters.json", "Relative path to parameters file")
	flag.Parse()

	initialize(*filename)

	http.HandleFunc("/distribute", distributeHandler)

	port := ":8080"
	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Could not start server: ", err)
	}
}
