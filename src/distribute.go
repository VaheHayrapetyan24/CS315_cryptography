package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
)

var nodeCounter uint32 = 0
var mutex sync.Mutex

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

func initialize() {
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

var ErrKeysExhausted = errors.New("Keys are exhausted. Cannot support more nodes.")

type DistributeResponse struct {
	Id   uint32             `json:"id"`
	Gcol [lambda + 1]uint64 `json:"g_col"`
	Acol [lambda + 1]uint64 `json:"a_col"`
}

func getNextKey() (uint32, [lambda + 1]uint64, [lambda + 1]uint64, error) {
	mutex.Lock()

	if nodeCounter == n {
		return 0, [lambda + 1]uint64{}, [lambda + 1]uint64{}, ErrKeysExhausted
	}

	id := nodeCounter
	nodeCounter++
	mutex.Unlock()
	return id, G_t[id], A[id], nil
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

	response := DistributeResponse{
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
	initialize()

	http.HandleFunc("/distribute", distributeHandler)

	port := ":8080"
	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Could not start server: ", err)
	}
}
