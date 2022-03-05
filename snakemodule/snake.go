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
	vao uint32,
	draw func(program, vertexArrayObject uint32, vec mgl32.Vec2),
) {
	position := food.cell.coords
	draw(program, vao, position)
}

type Snake struct {
	body []Cell
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

func (snake *Snake) Eat(food *Food, changeFoodPosition *bool) {
	snakeHead := snake.GetHead()
	if snakeHead.coords.X() == food.cell.coords.X() && snakeHead.coords.Y() == food.cell.coords.Y() {
		snake.body = append(snake.body, food.cell)
		*changeFoodPosition = true
	}
}

func (snake *Snake) Draw(
	program,
	vao uint32,
	draw func(program, vertexArrayObject uint32, vec mgl32.Vec2),
) {
	snakeBody := snake.body
	for i := 0; i < len(snakeBody); i++ {
		coords := snakeBody[i].coords
		draw(program, vao, coords)
	}
}

func (snake *Snake) GetHead() Cell {
	snakeBody := snake.body
	return snakeBody[len(snakeBody)-1]
}

func InitSnake(n int) *Snake {
	var snake Snake
	snake.body = make([]Cell, n)
	for i := 0; i < len(snake.body); i++ {
		snake.body[i].coords = mgl32.Vec2{float32(i), float32(0)}
	}
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
