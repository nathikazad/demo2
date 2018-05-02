package sim2

import (
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
)

// digraph - Describes an implementation of a simple weighted directed graph with underlying coords

// Vertex - struct for generic digraph vertex.
type Vertex struct {
  ID uint
  Pos Coords
  AdjEdges []Edge
}

// Edge - struct for directed weighted edge in digraph.
type Edge struct {
  ID uint
  Start *Vertex
  End *Vertex
  Weight float64
}

// Digraph - struct for Digraph object.
type Digraph struct {
  Vertices map[uint]*Vertex  // map vertex ID to vertex reference
  Edges map[uint]*Edge  // map edge ID to edge reference
}

// NewDigraph - Constructor for valid Digraph object.
func NewDigraph() *Digraph {
  d := new(Digraph)
  d.Vertices = make(map[uint]*Vertex)
  d.Edges = make(map[uint]*Edge)
  return d
}

// GetDigraphFromFile - Populate Digraph object from fomratted 'map' file.
func GetDigraphFromFile(fname string) (d *Digraph) {
  d = NewDigraph()

  file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

  scanner := bufio.NewScanner(file)

  // Scan input file line by line
  for scanner.Scan() {
    line := strings.Split(scanner.Text()," ")

    // Fetch id of Vertex described on this line
		idRead, _ := strconv.Atoi(line[0])
    id := uint(idRead)

    // Construct and populate new Vertex
    if _, ok := d.Vertices[id]; !ok {
      d.Vertices[id] = new(Vertex)
    }
    vert := d.Vertices[id]
    vert.ID = id

    // Parse input Vertex coordinates
    x, y := splitCSV(line[1])
    vert.Pos.X = x
    vert.Pos.Y = y

    // Parse adjacent vertices
		for _, point := range line[2:] {
      // Construct and populate new adjacent Vertex
      idReadNext, _ := strconv.Atoi(point)
      idNext := uint(idReadNext)
      if _, ok := d.Vertices[idNext]; !ok {
        d.Vertices[idNext] = new(Vertex)
      }
      vertNext := d.Vertices[idNext]

      // Construct and populate new edge between adjacent vertices
      edge := new(Edge)
      edge.ID = uint(len(d.Edges))
      edge.Start = vert
      edge.End = vertNext
      d.Edges[edge.ID] = edge
    }

    // Set starting edge weight based on distance
    for _, edge := range d.Edges {
      edge.Weight = edge.Start.Pos.Distance(edge.End.Pos)
    }
  }

  return d
}

// splitCSV - Split a CSV pair into constituent values.
func splitCSV(line string) (float64, float64) {
	numbers := strings.Split(line,",")
	if len(numbers) != 2 {
		log.Fatal("Line does not have format <number1>, <number2> ")
	}
	numberOne, err := strconv.Atoi(numbers[0])
	if err != nil {
		log.Fatal("Line does not have format <number1>, <number2> ")
	}
	numberTwo, err := strconv.Atoi(numbers[1])
	if err != nil {
		log.Fatal("Line does not have format <number1>, <number2> ")
	}
	return float64(numberOne), float64(numberTwo)
}

// ShortestPath - solve for the shortest deighted directional path from start to end vertex.
func (g Digraph) ShortestPath(startVertID, endVertID uint) (edgeIDs []uint, dist float64) {
  // TODO: implement this in a concurrency-safe manner: do not modify underlying values!
  return
}
