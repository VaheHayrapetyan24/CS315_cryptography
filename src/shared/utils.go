package shared

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadUint64WithDefault(prompt string, defaultVal uint64) uint64 {
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

func ReadUint32WithDefault(prompt string, defaultVal uint32) uint32 {
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

func ReadStringWithDefault(prompt string, defaultVal string) string {
	fmt.Printf("%s (default: %s): ", prompt, defaultVal)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}

	return input
}

func ReadYesNo(prompt string) bool {
	fmt.Printf("%s (y/n): ", prompt)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	return strings.ToLower(input) == "y" || strings.ToLower(input) == "yes"
}
