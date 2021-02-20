package client

import (
	"fmt"
	"runtime"
)

// SystemInfo provides easy access to metadata about the client.
type SystemInfo struct {
	BaseURL  string
	OS       string
	Arch     string
}

func (s SystemInfo) String() string {
	return fmt.Sprintf("{baseURL: %v, os: %v, arch: %v}", s.BaseURL, s.OS, s.Arch)
}

// GetSystemInfo returns information about the client.
func GetSystemInfo() SystemInfo {
	return SystemInfo{BaseURL: baseURL, OS: runtime.GOOS, Arch: runtime.GOARCH}
}
