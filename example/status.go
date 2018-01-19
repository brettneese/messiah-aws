package main

// StatusInfo stores a summary of information for the health check endpoint.
type StatusInfo struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

var version = "1.0.0"
var name = "Messiah Example API"

// GetStatus derives the system status
func GetStatus() StatusInfo {
	status := StatusInfo{
		Version: version,
		Name:    name,
	}

	return status
}
