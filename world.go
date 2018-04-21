package main

import "demo2/dijkstra"

type World struct {
	graph *dijkstra.Graph
	points map[uint]Point
}

