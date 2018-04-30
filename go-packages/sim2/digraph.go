package sim2

// digraph - Describes a generic interface for representing directed graphs

// Vertex - interface for generic digraph vertex.
type Vertex struct {
  ID uint
  Pos Coords
  AdjEdges []Edge
}

// Edge - interface for directed weighted edge in digraph.
type Edge struct {
  ID uint
  Start *Vertex
  End *Vertex
  Weight float64
}

// Digraph - interface for
type Digraph struct {
  Vertices map[uint]*Vertex  // map vertex ID to vertex reference
  Edges map[uint]*Edge  // map edge ID to edge reference
}

// ShortestPath - solve for the shortest deighted directional path from start to end vertex.
func (g Digraph) ShortestPath(startVertID, endVertID uint) (edgeIDs []uint, dist float64) {
  // TODO: implement this in a concurrency-safe manner: do not modify underlying values!
  return
}
