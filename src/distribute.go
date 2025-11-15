package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	// "sync/atomic"
)

var nodeCounter uint32 = 0
var mutex sync.Mutex

const q uint64 = 19
const lambda uint32 = 2
const n uint32 = 5

var G_t = [n][lambda + 1]uint64{ // this is transpose of G
	// {1, 1, 1, 1, 1},
	// {2, 4, 8, 16, 13},
	// {4, 16, 7, 9, 17},
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

	for i := uint32(0); i < n; i++ {
		fmt.Println(A[i])
	}
}

var ErrKeysExhausted = errors.New("Keys are exhausted. Cannot support more nodes.")

type DistributeResponse struct {
	id    uint32             `json:"id"`
	g_col [lambda + 1]uint64 `json:"g_col"`
	a_col [lambda + 1]uint64 `json:"a_col"`
}

func getNextKey() (uint32, [lambda + 1]uint64, [lambda + 1]uint64, error) {
	mutex.Lock()

	if nodeCounter == n-1 {
		return 0, [lambda + 1]uint64{}, [lambda + 1]uint64{}, ErrKeysExhausted
	}

	id := nodeCounter
	nodeCounter++
	mutex.Unlock()
	return id, G_t[id], A[id], nil
}

// Handler for the root path ("/")
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the TurboBlom Server Node!")
}

// Handler for the key distribution path ("/distribute")
func distributeHandler(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, you would perform the key distribution logic here.
	// For now, let's just confirm the route is working.

	// nodeID := r.URL.Query().Get("node")
	// if nodeID == "" {
	// 	nodeID = "Unknown Node"
	// }
	log.Printf("Received key distribution request")

	nodeId, g_col, a_col, err := getNextKey()
	if errors.Is(err, ErrKeysExhausted) {
		http.Error(w, err.Error(), http.StatusTeapot) // whatever, can do 404 also
		return
	}

	w.WriteHeader(http.StatusOK) // Set the HTTP status code to 200 OK
	w.Header().Set("Content-Type", "application/json")

	response := DistributeResponse{
		id:    nodeId,
		g_col: g_col,
		a_col: a_col,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If encoding fails for some reason, log the error and handle it (e.g., send a 500 error).
		log.Printf("Error encoding response: %v", err)

		// Use http.Error for simple, non-JSON error responses if the primary response failed
		// Note: This response will not be JSON, but plain text, as the header was already set.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Key distribution request received for Node: %s", nodeId)
}

func main() {
	initialize()
	// 1. Define Routes (Handlers)
	// http.HandleFunc registers a handler function for a specific path.
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/distribute", distributeHandler)

	// 2. Start the Server
	// The http.ListenAndServe function starts the server on a specific address and port.
	// The second argument (nil) means it uses the default handler (the routes defined above).

	port := ":8080"
	log.Printf("Starting server on port %s", port)

	// log.Fatal prints the error and then exits the program if the server fails to start
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Could not start server: ", err)
	}
}
