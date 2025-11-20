# stl-bounding-box

A Go library for calculating bounding boxes from STL (Stereolithography) files. Supports both binary and ASCII STL formats with automatic format detection.

## Features

- **Dual Format Support**: Automatically detects and parses both binary and ASCII STL files
- **Flexible API**: Use file paths or `io.Reader` for maximum flexibility
- **Vector Math Integration**: Uses `r3.Vec` from gonum for 3D vector operations
- **Helper Methods**: Built-in methods for dimensions and volume calculations
- **Pre-calculated Center**: Bounding box center is automatically computed and stored

## Installation

```bash
go get github.com/nfranczak/stl-bounding-box
```

## Usage

### As a Library

```go
package main

import (
    "fmt"
    "log"

    stl "github.com/nfranczak/stl-bounding-box"
)

func main() {
    // Option 1: From file path
    bbox, err := stl.CalculateBoundingBoxFromFile("model.stl")
    if err != nil {
        log.Fatal(err)
    }

    // Option 2: From io.Reader
    file, _ := os.Open("model.stl")
    defer file.Close()
    bbox, err := stl.CalculateBoundingBox(file)
    if err != nil {
        log.Fatal(err)
    }

    // Access bounding box properties
    fmt.Printf("Min: (%.2f, %.2f, %.2f)\n", bbox.MinX, bbox.MinY, bbox.MinZ)
    fmt.Printf("Max: (%.2f, %.2f, %.2f)\n", bbox.MaxX, bbox.MaxY, bbox.MaxZ)
    fmt.Printf("Center: (%.2f, %.2f, %.2f)\n", bbox.Center.X, bbox.Center.Y, bbox.Center.Z)

    // Use helper methods
    width, height, depth := bbox.Dimensions()
    fmt.Printf("Dimensions: %.2f x %.2f x %.2f\n", width, height, depth)
    fmt.Printf("Volume: %.2f\n", bbox.Volume())
}
```

### As a CLI Tool

```bash
go run main.go model.stl
```

Output:
```
Bounding Box:
  Min: (0.00, 0.00, 0.00)
  Max: (100.00, 50.00, 25.00)
  Dimensions: (100.00, 50.00, 25.00)
  Center: (50.00, 25.00, 12.50)
  Volume: 125000.00
```

## API Reference

### Types

#### `BoundingBox`
```go
type BoundingBox struct {
    MinX, MinY, MinZ float32
    MaxX, MaxY, MaxZ float32
    Center           r3.Vec
}
```

#### `Triangle`
```go
type Triangle struct {
    Normal   r3.Vec
    Vertices [3]r3.Vec
}
```

### Functions

#### `CalculateBoundingBoxFromFile(filePath string) (*BoundingBox, error)`
Reads an STL file from the given path and returns its bounding box. Automatically detects binary or ASCII format.

#### `CalculateBoundingBox(r io.Reader) (*BoundingBox, error)`
Reads an STL file from an `io.Reader` and returns its bounding box. Useful for working with streams, HTTP responses, or embedded files.

### Methods

#### `(bb *BoundingBox) Dimensions() (width, height, depth float32)`
Returns the width, height, and depth of the bounding box.

#### `(bb *BoundingBox) Volume() float32`
Returns the volume of the bounding box.

## STL Format Support

This library supports both STL format variants:

- **Binary STL**: The standard binary format with 80-byte header, triangle count, and packed vertex data
- **ASCII STL**: The text-based format with `solid`, `facet`, `vertex`, and `endfacet` keywords

Format detection is automatic - you don't need to specify which format you're using.

## Dependencies

- [gonum.org/v1/gonum](https://github.com/gonum/gonum) - For `r3.Vec` 3D vector type

## License

MIT