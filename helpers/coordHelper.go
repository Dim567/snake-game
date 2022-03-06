package helpers

import "github.com/go-gl/mathgl/mgl32"

func CellsDifference(firstSlice, secondSlice []int) []int {
	diffCells := make([]int, 0, len(secondSlice))
	cellsMap := make(map[int]bool)
	for _, val := range secondSlice {
		cellsMap[val] = true
	}
	for _, val := range firstSlice {
		if _, ok := cellsMap[val]; !ok {
			diffCells = append(diffCells, val)
		}
	}
	return diffCells
}

func CoordsToIndex(x, y int) int {
	return y*10 + x
}

func IndexToCoords(i int) (x, y int) {
	x = i % 10
	y = i / 10
	return
}

func Distance(coord1, coord2 mgl32.Vec2) float32 {
	return coord1.Sub(coord2).Len()
}
