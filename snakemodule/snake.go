package snakemodule

import (
	"math/rand"
	"snakegame/helpers"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

// Cell
type Cell struct {
	coords mgl32.Vec2
}

func (cell *Cell) GetCoords() mgl32.Vec2 {
	return cell.coords
}

type Food struct {
	cell Cell
}

func (food *Food) SetPosition(possibleCells []int) {
	if len(possibleCells) > 0 {
		seedValue := time.Now().UnixNano()
		seed := rand.NewSource(seedValue)
		randNew := rand.New(seed)
		chosenCell := possibleCells[randNew.Intn(len(possibleCells))]
		x, y := helpers.IndexToCoords(chosenCell)
		food.cell.coords = mgl32.Vec2{float32(x), float32(y)}
	}
}

func (food *Food) Draw(
	program,
	vao,
	texture uint32,
	draw func(program, vertexArrayObject, texture uint32, vec mgl32.Vec2),
) {
	position := food.cell.coords
	draw(program, vao, texture, position)
}

type Snake struct {
	body                  []Cell
	front                 mgl32.Vec2
	intersectionThreshold float32
}

func (snake *Snake) GetFront() mgl32.Vec2 {
	return snake.front
}

func (snake *Snake) SetFront(vec mgl32.Vec2) {
	snake.front = vec
}

func (snake *Snake) Move(vec mgl32.Vec2) {
	snakeBody := snake.body
	headIndex := len(snakeBody) - 1
	headCoords := snakeBody[headIndex].coords

	if vec.X() != headCoords.X() || vec.Y() != headCoords.Y() {
		for i := 0; i < headIndex; i++ {
			newCoords := snakeBody[i+1].coords
			snakeBody[i].coords = newCoords
		}
		snakeBody[headIndex].coords = vec
	}
}

func (snake *Snake) Eat(food Food) bool {
	snakeHead := snake.GetFront()
	foodCoords := food.cell.GetCoords()
	threshold := snake.intersectionThreshold
	if helpers.Distance(snakeHead, foodCoords) < threshold {
		snake.body = append(snake.body, food.cell)
		return true
	}
	return false
}

func (snake *Snake) Draw(
	program,
	vao,
	texture uint32,
	draw func(program, vertexArrayObject, texture uint32, vec mgl32.Vec2),
) {
	snakeBody := snake.body
	for i := 0; i < len(snakeBody); i++ {
		coords := snakeBody[i].coords
		draw(program, vao, texture, coords)
	}
}

func (snake *Snake) GetHead() Cell {
	snakeBody := snake.body
	return snakeBody[len(snakeBody)-1]
}

func (snake *Snake) CheckIntersection() bool {
	snakeBody := snake.body
	snakeHead := snake.GetFront()
	threshold := snake.intersectionThreshold
	for i := 0; i < len(snakeBody)-1; i++ {
		if helpers.Distance(snakeHead, snakeBody[i].coords) < threshold {
			// fmt.Println("snakeLength:", len(snakeBody))
			// fmt.Println("item#:", i)
			// fmt.Println("Distance:", helpers.Distance(snakeHead, snakeBody[i].coords))
			return true
		}
	}
	return false
}

func InitSnake(snakeLength int, intersectionThreshold float32) *Snake {
	var snake Snake
	snake.body = make([]Cell, snakeLength)
	for i := 0; i < snakeLength; i++ {
		snake.body[i].coords = mgl32.Vec2{float32(i), float32(0)}
	}
	snake.SetFront(snake.GetHead().coords)
	snake.intersectionThreshold = intersectionThreshold
	return &snake
}

func GetPossibleCells(snake *Snake, fieldCells []int) []int {
	busyCells := make([]int, len(snake.body))
	for i, val := range snake.body {
		busyCells[i] = helpers.CoordsToIndex(int(val.coords.X()), int(val.coords.Y()))
	}
	possibleCells := helpers.CellsDifference(fieldCells, busyCells)
	return possibleCells
}
