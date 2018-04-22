package agents

import (
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
	"demo2/dijkstra"

)

func GetWorldFromFile(filename string) (*World) {
	world := NewWorld()
	world.nodes, world.edges = readPointsFromFile(filename)
	world.graph = addNodesToGraph(world.nodes)
	return world;
}

func readPointsFromFile(filename string) (map[uint]*Node, map[uint]*Edge) {
	nodes := make(map[uint]*Node)
	edges := make(map[uint]*Edge)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line :=strings.Split(scanner.Text()," ")
		id, _ := strconv.Atoi(line[0])
		if _, ok := nodes[uint(id)]; !ok {
			nodes[uint(id)] = &Node{}
		}
		node := nodes[uint(id)]
		node.id = uint(id)
		x, y := splitCommaSepNumbers(line[1])
		node.coordinates.x = x
		node.coordinates.y = y
		for _, point := range line[2:] {
			nextNodeId, _ := strconv.Atoi(point)
			if _, ok := nodes[uint(nextNodeId)]; !ok {
				nodes[uint(nextNodeId)] = &Node{}
			}
			nextNode := nodes[uint(nextNodeId)]
			node.nextNodes = append(node.nextNodes, nextNode)

			if node.nextEdges == nil {
				node.nextEdges = make(map[*Node]*Edge)
			}
			edge := Edge{id:uint(len(edges)), startingNode:node, endingNode:nextNode}
			node.nextEdges[nextNode] = &edge
			edges[edge.id] = &edge
		}
	}

	for _, edge := range edges {
		edge.weight = edge.startingNode.coordinates.Distance(edge.endingNode.coordinates)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nodes, edges
}

func splitCommaSepNumbers(line string) (uint, uint) {
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
	return uint(numberOne), uint(numberTwo)
}

// TODO: toss this for our own concurrency safe implementation
func addNodesToGraph(nodes map[uint]*Node) (*dijkstra.Graph) {
	graph:=dijkstra.NewGraph()
	for k ,_ := range nodes {
		graph.AddVertex(int(k))
	}

	for _ , node := range nodes {
		for _, nextNode := range node.nextNodes {
			graph.AddArc(int(node.id), int(nextNode.id), int64(node.nextEdges[nextNode].weight))
		}
	}
	return graph
}