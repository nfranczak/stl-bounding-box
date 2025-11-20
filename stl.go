package stl

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/spatial/r3"
)

// Triangle represents a triangle in an STL file
type Triangle struct {
	Normal   r3.Vec
	Vertices [3]r3.Vec
}

// BoundingBox stores the min and max coordinates of a 3D model
type BoundingBox struct {
	MinX, MinY, MinZ float32
	MaxX, MaxY, MaxZ float32
	Center           r3.Vec
}

// Dimensions returns the width, height, and depth of the bounding box
func (bb *BoundingBox) Dimensions() (width, height, depth float32) {
	return bb.MaxX - bb.MinX, bb.MaxY - bb.MinY, bb.MaxZ - bb.MinZ
}

// Volume returns the volume of the bounding box
func (bb *BoundingBox) Volume() float32 {
	w, h, d := bb.Dimensions()
	return w * h * d
}

// CalculateBoundingBoxFromFile reads an STL file from the given path
// and returns its bounding box. Supports both binary and ASCII STL formats.
func CalculateBoundingBoxFromFile(filePath string) (*BoundingBox, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	return CalculateBoundingBox(file)
}

// CalculateBoundingBox reads an STL file from the given io.Reader
// and returns its bounding box. Supports both binary and ASCII STL formats.
// The function automatically detects the format.
func CalculateBoundingBox(r io.Reader) (*BoundingBox, error) {
	// Read first 80 bytes to check if it's ASCII or binary
	header := make([]byte, 80)
	n, err := io.ReadFull(r, header)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	// Check if it's ASCII by looking for "solid" keyword
	headerStr := string(header[:n])
	if strings.HasPrefix(strings.TrimSpace(headerStr), "solid") {
		// Might be ASCII, need to verify by checking if "facet" follows
		return parseASCII(io.MultiReader(strings.NewReader(headerStr), r))
	}

	// Binary STL format
	return parseBinary(io.MultiReader(strings.NewReader(headerStr), r))
}

// binaryTriangle is used for reading binary STL format (float32)
type binaryTriangle struct {
	Normal   [3]float32
	Vertices [3][3]float32
}

// parseBinary parses a binary STL file
func parseBinary(r io.Reader) (*BoundingBox, error) {
	// Skip 80-byte header
	header := make([]byte, 80)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	// Read number of triangles
	var numTriangles uint32
	if err := binary.Read(r, binary.LittleEndian, &numTriangles); err != nil {
		return nil, fmt.Errorf("error reading number of triangles: %w", err)
	}

	bbox := &BoundingBox{
		MinX: math.MaxFloat32, MinY: math.MaxFloat32, MinZ: math.MaxFloat32,
		MaxX: -math.MaxFloat32, MaxY: -math.MaxFloat32, MaxZ: -math.MaxFloat32,
	}

	for i := 0; i < int(numTriangles); i++ {
		var binTriangle binaryTriangle
		if err := binary.Read(r, binary.LittleEndian, &binTriangle); err != nil {
			return nil, fmt.Errorf("error reading triangle %d: %w", i, err)
		}

		// Convert to r3.Vec
		vertices := [3]r3.Vec{
			{X: float64(binTriangle.Vertices[0][0]), Y: float64(binTriangle.Vertices[0][1]), Z: float64(binTriangle.Vertices[0][2])},
			{X: float64(binTriangle.Vertices[1][0]), Y: float64(binTriangle.Vertices[1][1]), Z: float64(binTriangle.Vertices[1][2])},
			{X: float64(binTriangle.Vertices[2][0]), Y: float64(binTriangle.Vertices[2][1]), Z: float64(binTriangle.Vertices[2][2])},
		}

		updateBoundingBox(bbox, vertices[:])

		// Skip 2-byte attribute byte count
		var attributeByteCount uint16
		if err := binary.Read(r, binary.LittleEndian, &attributeByteCount); err != nil {
			return nil, fmt.Errorf("error reading attribute byte count: %w", err)
		}
	}

	// Calculate center
	bbox.Center = r3.Vec{
		X: float64((bbox.MinX + bbox.MaxX) / 2),
		Y: float64((bbox.MinY + bbox.MaxY) / 2),
		Z: float64((bbox.MinZ + bbox.MaxZ) / 2),
	}

	return bbox, nil
}

// parseASCII parses an ASCII STL file
func parseASCII(r io.Reader) (*BoundingBox, error) {
	scanner := bufio.NewScanner(r)
	bbox := &BoundingBox{
		MinX: math.MaxFloat32, MinY: math.MaxFloat32, MinZ: math.MaxFloat32,
		MaxX: -math.MaxFloat32, MaxY: -math.MaxFloat32, MaxZ: -math.MaxFloat32,
	}

	var currentTriangle [3]r3.Vec
	vertexIndex := 0
	inFacet := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)

		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "facet":
			inFacet = true
			vertexIndex = 0
		case "vertex":
			if !inFacet || len(fields) < 4 {
				return nil, fmt.Errorf("invalid vertex line: %s", line)
			}
			if vertexIndex >= 3 {
				return nil, fmt.Errorf("too many vertices in facet")
			}

			x, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing x coordinate: %w", err)
			}
			y, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing y coordinate: %w", err)
			}
			z, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing z coordinate: %w", err)
			}

			currentTriangle[vertexIndex] = r3.Vec{
				X: x,
				Y: y,
				Z: z,
			}
			vertexIndex++
		case "endfacet":
			if vertexIndex != 3 {
				return nil, fmt.Errorf("incomplete triangle, got %d vertices", vertexIndex)
			}
			updateBoundingBox(bbox, currentTriangle[:])
			inFacet = false
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Check if we found any triangles
	if bbox.MinX == math.MaxFloat32 {
		return nil, fmt.Errorf("no triangles found in STL file")
	}

	// Calculate center
	bbox.Center = r3.Vec{
		X: float64((bbox.MinX + bbox.MaxX) / 2),
		Y: float64((bbox.MinY + bbox.MaxY) / 2),
		Z: float64((bbox.MinZ + bbox.MaxZ) / 2),
	}

	return bbox, nil
}

// updateBoundingBox updates the bounding box with the given vertices
func updateBoundingBox(bbox *BoundingBox, vertices []r3.Vec) {
	for _, vertex := range vertices {
		x, y, z := float32(vertex.X), float32(vertex.Y), float32(vertex.Z)

		if x < bbox.MinX {
			bbox.MinX = x
		}
		if y < bbox.MinY {
			bbox.MinY = y
		}
		if z < bbox.MinZ {
			bbox.MinZ = z
		}

		if x > bbox.MaxX {
			bbox.MaxX = x
		}
		if y > bbox.MaxY {
			bbox.MaxY = y
		}
		if z > bbox.MaxZ {
			bbox.MaxZ = z
		}
	}
}
