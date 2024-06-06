package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/siemens-mobile-hacks/pattern89/pkg/pattern89"
)

var (
	fullFlash  = flag.String("fullflash", "", "Path to the fullflash file.")
	patternStr = flag.String("pattern", "", "Pattern to look for.")
)

func main() {
	flag.Parse()

	// Check if the file path is provided.
	if *fullFlash == "" || *patternStr == "" {
		fmt.Println("Usage: pattern89 -fullflash <file-path> -pattern <pattern>")
		return
	}

	// Convert the pattern to a byte slice with wildcards represented as 0xFF.
	pattern, err := pattern89.ParsePattern(*patternStr)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Open the file.
	file, err := os.Open(*fullFlash)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a buffer to read chunks of the file
	const chunkSize = 4096
	buffer := make([]byte, chunkSize)
	offset := int64(0)

	for {
		// Read a chunk of the file.
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading file:", err)
			return
		}

		// Search for the pattern in the chunk.
		if pos := pattern89.FindPattern(buffer[:bytesRead], pattern); pos != -1 {
			fmt.Printf("Pattern found at offset: %08X\n", offset+int64(pos))
			return
		}

		// Move the offset to the beginning of the next chunk
		// (keeping some overlap to account for patterns that might be split across chunks).
		offset += int64(bytesRead) - int64(pattern.Length())
		if bytesRead < chunkSize {
			break
		}
		_, err = file.Seek(offset, io.SeekStart)
		if err != nil {
			fmt.Println("Error seeking file:", err)
			return
		}
	}

	fmt.Println("Pattern not found")
}
