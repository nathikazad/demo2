package main

import (
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
	"demo2/dijkstra"
)

type Point struct {
	id uint
	x uint
	y  uint
	adjacent_points []uint
}

func getWorldFromFile(filename string) (*World) {
	world := new(World)
	world.points = readPointsFromFile(filename)
	world.graph = convertPointsToGraph(world.points)
	return world;
}

func readPointsFromFile(filename string) (map[uint]Point) {
	points := make(map[uint]Point)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line :=strings.Split(scanner.Text()," ")
		index, _ := strconv.Atoi(line[0])
		x, y := splitCommaSepNumbers(line[1])
		adjacent_points := []uint{}
		for _, point := range line[2:] {
			temp, _ := strconv.Atoi(point)
			adjacent_points = append(adjacent_points, uint(temp))
		}
		points[uint(index)] = Point{x: x, y: y, adjacent_points:adjacent_points}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return points
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

func convertPointsToGraph(points map[uint]Point) (*dijkstra.Graph) {
	graph:=dijkstra.NewGraph()
	for k ,_ := range points {
		graph.AddVertex(int(k))
	}

	for index , point := range points {
		for _, a_point := range point.adjacent_points {
			distance := int64(point.Distance(points[a_point]))
			graph.AddArc(int(index), int(a_point), distance)
		}
	}
	return graph
}