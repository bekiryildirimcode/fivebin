package main

import (
	"fmt"
	"runtime"
)

// printMemUsage prints current memory usage statistics
// Useful for monitoring memory efficiency during development
func printMemUsage(stage string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("[%s]\n", stage)
	fmt.Printf("Alloc = %v KB\n", m.Alloc/1024)
	fmt.Printf("TotalAlloc = %v KB\n", m.TotalAlloc/1024)
	fmt.Printf("Sys = %v KB\n", m.Sys/1024)
	fmt.Printf("NumGC = %v\n\n", m.NumGC)
}

// formatMemorySize formats bytes into a human-readable string
func formatMemorySize(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
