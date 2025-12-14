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

func listen(port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("error listening on %s: %w\n", port, err)
	}
	defer listener.Close()

	fmt.Printf("Listening on %s...\n", port)

	conn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("error accepting connection: %w\n", err)
	}

	return handleConnection(conn)
}

func handleConnection(conn net.Conn) error {
	defer conn.Close()
	fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())

	decoder := json.NewDecoder(conn)
	var meta FileMetadata
	if err := decoder.Decode(&meta); err != nil {
		return fmt.Errorf("error decoding metadata: %w", err)
	}

	filename := filepath.Base(meta.Name)
	fmt.Printf("Receiving file: %s (%d bytes)\n", filename, meta.Size)

	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filename, err)
	}
	defer outFile.Close()

	// Use MultiReader to combine buffered data from decoder and the rest of the connection
	source := io.MultiReader(decoder.Buffered(), conn)

	// Consume the newline character added by json.Encoder
	var oneByte [1]byte
	if _, err := source.Read(oneByte[:]); err != nil {
		return fmt.Errorf("error reading separator: %w", err)
	}

	bar := progressbar.DefaultBytes(
		meta.Size,
		"receiving",
	)

	hasher := crc32.NewIEEE()
	writer := io.MultiWriter(outFile, bar, hasher)

	n, err := io.CopyN(writer, source, meta.Size)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error receiving file content: %w", err)
	}

	receivedHash := hex.EncodeToString(hasher.Sum(nil))
	if receivedHash != meta.Hash {
		return fmt.Errorf("hash mismatch! (expected: %s | received: %s)", meta.Hash, receivedHash)
	}

	fmt.Printf("\nSuccessfully received %s (%d bytes)\nVerified Hash (CRC32): %s\n", filename, n, receivedHash)
	return nil
}
