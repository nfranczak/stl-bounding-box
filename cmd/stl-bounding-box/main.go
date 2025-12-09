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
	fmt.Printf("  Min: (%.5f, %.5f, %.5f)\n", bbox.MinX, bbox.MinY, bbox.MinZ)
	fmt.Printf("  Max: (%.5f, %.5f, %.5f)\n", bbox.MaxX, bbox.MaxY, bbox.MaxZ)
	fmt.Printf("  Dimensions: (%.5f, %.5f, %.5f)\n", width, height, depth)
	fmt.Printf("  Center: (%.5f, %.5f, %.5f)\n", bbox.Center.X, bbox.Center.Y, bbox.Center.Z)
	fmt.Printf("  Volume: %.5f\n", bbox.Volume())
}
