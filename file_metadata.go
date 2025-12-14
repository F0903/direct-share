package main

// FileMetadata sent before the file content
type FileMetadata struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Hash string `json:"hash"`
}
