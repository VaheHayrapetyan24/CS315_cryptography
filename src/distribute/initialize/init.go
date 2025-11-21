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
	distribute_models "turbobloom/distribute/models"
)

func readUint64WithDefault(prompt string, defaultVal uint64) uint64 {
	fmt.Printf("%s (default: %d): ", prompt, defaultVal)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}

	val, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		fmt.Printf("Invalid input, using default: %d\n", defaultVal)
		return defaultVal
	}

	return val
}

func readUint32WithDefault(prompt string, defaultVal uint32) uint32 {
	fmt.Printf("%s (default: %d): ", prompt, defaultVal)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}

	val, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		fmt.Printf("Invalid input, using default: %d\n", defaultVal)
		return defaultVal
	}

	return uint32(val)
}

func readStringWithDefault(prompt string, defaultVal string) string {
	fmt.Printf("%s (default: %s): ", prompt, defaultVal)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}

	return input
}

func readYesNo(prompt string) bool {
	fmt.Printf("%s (y/n): ", prompt)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	return strings.ToLower(input) == "y" || strings.ToLower(input) == "yes"
}

func generateG(n uint32, lambda uint32) [][]uint64 {
	rand.Seed(time.Now().UnixNano())
	g := make([][]uint64, n)
	for i := 0; i < int(n); i++ {
		g[i] = make([]uint64, lambda+1)
		for j := 0; j < int(lambda+1); j++ {
			g[i][j] = uint64(rand.Intn(2))
		}
	}
	return g
}

func generateD(lambda uint32, q uint64) [][]uint64 {
	rand.Seed(time.Now().UnixNano())
	d := make([][]uint64, lambda+1)
	len := lambda + 1
	for i := 0; i < int(len); i++ {
		d[i] = make([]uint64, lambda+1)
		for j := 0; j < int(i+1); j++ {
			d[i][j] = uint64(rand.Intn(int(q-1))) + 1
			d[j][i] = d[i][j]
		}
	}

	return d
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Enter cryptographic parameters:")

	// Q with default 19
	q := readUint64WithDefault("Enter Q", 19)

	// N with default 5
	n := readUint32WithDefault("Enter N", 5)

	// Lambda with default 2
	lambda := readUint32WithDefault("Enter Lambda", 2)

	fmt.Printf("\nParameters set:\n")
	fmt.Printf("Q: %d\n", q)
	fmt.Printf("N: %d\n", n)
	fmt.Printf("Lambda: %d\n", lambda)

	// Ask for file path to store parameters
	filepath := readStringWithDefault("\nEnter relative path to store parameters file", "./parameters/parameters.json")

	// Check if file exists
	if _, err := os.Stat(filepath); err == nil {
		// File exists, ask to overwrite
		if !readYesNo("File already exists. Overwrite?") {
			fmt.Println("File not saved.")
			return
		}
	}

	// Create parameters struct
	params := distribute_models.Parameters{
		Q:      q,
		N:      n,
		Lambda: lambda,
		G:      generateG(n, lambda),
		D:      generateD(lambda, q),
		Count:  0,
	}

	// Marshal to JSON
	data, err := json.Marshal(params)
	if err != nil {
		fmt.Printf("Error marshaling parameters: %v\n", err)
		return
	}

	// Write to file
	err = os.WriteFile(filepath, data, 0600)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Parameters saved to %s\n", filepath)
}
