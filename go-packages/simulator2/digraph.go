package simulator2

// digraph - Describes a generic interface for representing directed graphs

// Vertex - interface for generic digraph vertex
type Vertex interface {
  GetID() uint
  GetEdges() []*Edge
}

// Edge - interface for directed weighted edge in digraph
type Edge interface {
  GetID() uint
  GetStartVert() *Vertex
  GetEndVert() *Vertex
  GetWeight() float64
}

// Digraph - interface for
type Digraph interface {
  GetVerts() map[uint]*Vertex
  GetEdges() map[uint]*Edge
}

// ShortestPath - solve for the shortest deighted directional path from start to end vertex
func (g Digraph) ShortestPath(startVertID, endVertID uint) (edgeIDs []uint, dist float64) {
  // TODO: implement this in a future branch. Assume Vertex and Edge interfaces are concurrency-safe
  return
}
