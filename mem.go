package main

import (
	"fmt"
	"log"
	"runtime"
)

var (
	// debug
	m   runtime.MemStats
	kib uint64
	mib uint64
)

func MonitorMemory() {
	// debug memory usage
	runtime.ReadMemStats(&m)

	// https://forum.golangbridge.org/t/how-can-i-know-limit-memory-size-of-golang-application-sys-heap-stack/20070/2
	kib = m.HeapAlloc / 1024
	mib = kib / 1024
	fmt.Printf("\tHeap allocated = %v MiB (%vKiB)\n", mib, kib)

	if mib > MaxMemAlloc {
		log.Fatal(fmt.Sprintf("Heap allocation %dMiB exceeded max %dMiB", mib, MaxMemAlloc))
	}
}
