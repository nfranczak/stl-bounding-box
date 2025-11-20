package main

import (
	"fmt"
	"os"

	stl "github.com/nfranczak/stl-bounding-box"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: stl-bounding-box <file.stl>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	bbox, err := stl.CalculateBoundingBoxFromFile(filePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	width, height, depth := bbox.Dimensions()

	fmt.Printf("Bounding Box:\n")
	fmt.Printf("  Min: (%.2f, %.2f, %.2f)\n", bbox.MinX, bbox.MinY, bbox.MinZ)
	fmt.Printf("  Max: (%.2f, %.2f, %.2f)\n", bbox.MaxX, bbox.MaxY, bbox.MaxZ)
	fmt.Printf("  Dimensions: (%.2f, %.2f, %.2f)\n", width, height, depth)
	fmt.Printf("  Center: (%.2f, %.2f, %.2f)\n", bbox.Center.X, bbox.Center.Y, bbox.Center.Z)
	fmt.Printf("  Volume: %.2f\n", bbox.Volume())
}
