package main

import (
	"fmt"
	"os"
)

func writeDataToFile(data []byte, filename string) error {
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
