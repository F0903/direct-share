package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func send(addr, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}
	if fileInfo.Size() <= 0 {
		return fmt.Errorf("file size must be greater than 0")
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("error connecting to %s: %w", addr, err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", addr)

	// Calculate CRC32 hash
	fmt.Println("Calculating file hash...")
	hasher := crc32.NewIEEE()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("error calculating hash: %w", err)
	}
	fileHash := hex.EncodeToString(hasher.Sum(nil))
	fmt.Printf("File hash: %s\n", fileHash)

	// Reset file pointer to beginning
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("error seeking file: %w", err)
	}
	// Prepare metadata
	meta := FileMetadata{
		Name: filepath.Base(filePath),
		Size: fileInfo.Size(),
		Hash: fileHash,
	}

	// Send metadata
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(meta); err != nil {
		return fmt.Errorf("error sending metadata: %w", err)
	}

	fmt.Printf("Sending file: %s (%d bytes)...\n", meta.Name, meta.Size)

	bar := progressbar.DefaultBytes(
		meta.Size,
		"sending",
	)

	// Send file content
	n, err := io.Copy(io.MultiWriter(conn, bar), file)
	if err != nil {
		return fmt.Errorf("error sending file data: %w", err)
	}

	fmt.Printf("\nSent %d bytes successfully.\n", n)
	return nil
}
