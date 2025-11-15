package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// var nodeCounter uint32 = 0
// var mutex sync.Mutex

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

var MyId uint32 = 2
var U = [lambda + 1]uint64{5, 14, 18}

func getKey(id uint32) uint64 {
	var k uint64 = 0
	for i := uint32(0); i <= lambda; i++ {
		k += (U[i] * G_t[id][i]) % q
	}
	return k % q
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	nodeIdHeader := r.Header.Get("x-node-id")
	if nodeIdHeader == "" {
		http.Error(w, "x-node-id header is required", http.StatusBadRequest)
		return
	}

	nodeId64, err := strconv.ParseUint(nodeIdHeader, 10, 32)
	if err != nil {
		http.Error(w, "x-node-id header is not a valid unsigned int", http.StatusBadRequest)
		return
	}

	log.Printf("Received encoded message %v", nodeId64)
	nodeId32 := uint32(nodeId64)

	k := getKey(nodeId32)

	log.Printf("Resulting key %v", k)

	w.WriteHeader(http.StatusOK)
}

func main() {
	portPtr := flag.Int("port", 8080, "Port to run the server on")
	flag.Parse()

	http.HandleFunc("/message", messageHandler)

	port := fmt.Sprintf(":%d", *portPtr)
	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Could not start server: ", err)
	}
}
