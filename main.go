package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"snakegame/graphics"
	"snakegame/snakemodule"
	"strconv"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Window initial sizes
const (
	windowWidth  = 800
	windowHeight = 800
)

// Cells number
const cellsNumber = 10

var gameOver = false
var showLevel = true
var resetLevel = false
var loadLevel = false

const levelsNumber = 5

var startGame = true
var pauseGame = false
var gameLevel int = 0
var eatenFoodCounter int = 0
var progressSaved = true
var startLevel = true

const (
	fromStart int8 = 1
	toStart   int8 = -1
)

// Movement direction
var direction = fromStart

// Time settings
var startTime float64
var endTime float64
var period float32

// Speed settings
var timeWindow float32 = 0.5
var intersectionThreshold float32
var higherEdge, lowerEdge float32
var timeToMove = false

// Field settings
var fieldCells = make([]int, cellsNumber*cellsNumber)

var snake *snakemodule.Snake
var food snakemodule.Food

// Movement horizontal or vertical
var horizontalMove = true

var foodWasEaten = false

func main() {
	runtime.LockOSThread()

	defer graphics.Terminate()
	err := graphics.Init("Snake game", windowWidth, windowHeight)
	if err != nil {
		panic(err)
	}
	graphics.SetResizeWindowCallback(resizeWindowCallback)
	graphics.SetKeyInputCallback(keyInputCallback)

	// Create and load textures
	var snakeTexture = graphics.LoadTexture("snake_skin.png")
	var backgroundTexture = graphics.LoadTexture("background.png")
	var gameOverTexture = graphics.LoadTexture("game_over.png")
	var levelTexture0 = graphics.LoadTexture("level_1.png")
	var levelTexture1 = graphics.LoadTexture("level_2.png")
	var levelTexture2 = graphics.LoadTexture("level_3.png")
	var levelTexture3 = graphics.LoadTexture("level_4.png")
	var finishLevelTexture = graphics.LoadTexture("finish.png")
	var startGameTexture = graphics.LoadTexture("start_game.png")

	levelTextures := [levelsNumber]uint32{
		levelTexture0,
		levelTexture1,
		levelTexture2,
		levelTexture3,
		finishLevelTexture,
	}

	for i := 0; i < len(fieldCells); i++ {
		fieldCells[i] = i
	}

	resetGame(0, 3)
	gameLogic := func() {
		switch {
		case startGame:
			showLevel = true
			drawBackground(startGameTexture)
		case gameOver:
			showLevel = true
			drawBackground(gameOverTexture)
		case showLevel:
			if !progressSaved {
				saveProgress(gameLevel)
				progressSaved = true
			}
			if startLevel {
				gameLevel = 0
				timeWindow = getTimeWindow(gameLevel)
				if loadLevel {
					gameLevel = loadProgress()
					timeWindow = getTimeWindow(gameLevel)
					loadLevel = false
				}
			}
			startLevel = false
			resetLevel = true
			textureItem := levelTextures[gameLevel]
			drawBackground(textureItem)
		case pauseGame:
			timeToMove = false
			startTime = glfw.GetTime()
			drawBackground(backgroundTexture)
			food.Draw(snakeTexture, drawObject)
			snake.Draw(snakeTexture, drawObject)
		case resetLevel:
			resetLevel = false
			resetGame(gameLevel, 3)
			fallthrough
		default:
			endTime = glfw.GetTime()
			period = float32(endTime - startTime)

			if period >= timeWindow {
				startTime = endTime
				timeToMove = true
			}

			snakeHead := snake.GetHead()
			x, y := snakeHead.GetCoords().Elem()
			frontX, frontY := x, y
			if horizontalMove {
				frontX += float32(direction) * float32(period)
			} else {
				frontY += float32(direction) * float32(period)
			}
			snake.SetFront(mgl32.Vec2{frontX, frontY})

			if frontX >= higherEdge ||
				frontX <= lowerEdge ||
				frontY >= higherEdge ||
				frontY <= lowerEdge ||
				snake.CheckIntersection() ||
				gameLevel == levelsNumber-1 {
				gameOver = true
			}

			if timeToMove {
				timeToMove = false
				if horizontalMove {
					x += float32(direction)
				} else {
					y += float32(direction)
				}

				foodWasEaten = snake.Eat(food)
				snake.Move(mgl32.Vec2{x, y})
			}

			if foodWasEaten {
				foodWasEaten = false
				eatenFoodCounter += 1
				if eatenFoodCounter == getFoodLimit(gameLevel) {
					gameLevel += 1
					timeWindow = getTimeWindow(gameLevel)
					progressSaved = false
					showLevel = true
				} else {
					setFoodPosition(fieldCells)
				}
			}

			drawBackground(backgroundTexture)

			if period < (2*timeWindow/7) || period > (5*timeWindow/7) {
				food.Draw(snakeTexture, drawObject)
			}
			snake.Draw(snakeTexture, drawObject)
		}
	}

	graphics.MainLoop(gameLogic)
}

func drawObject(texture uint32, vec mgl32.Vec2) {
	scaleFactor := float32(2.0 / cellsNumber)
	scale := mgl32.Scale3D(scaleFactor, scaleFactor, 1)
	xPos := vec.X()*scaleFactor - 1
	yPos := vec.Y()*scaleFactor - 1
	translate := mgl32.Translate3D(xPos, yPos, 0)
	transform := translate.Mul4(scale)
	graphics.Draw(texture, transform)
}

func drawBackground(texture uint32) {
	scale := mgl32.Scale3D(2, 2, 1)
	translate := mgl32.Translate3D(-1, -1, 0)
	transform := translate.Mul4(scale)
	graphics.Draw(texture, transform)
}

func keyInputCallback(key graphics.KeyValue, action graphics.KeyAction) {
	if !pauseGame && !gameOver && !showLevel {
		if (key == graphics.KeyW || key == graphics.KeyUp) && action == graphics.Press {
			if !horizontalMove && direction == -1 {
				return
			}
			direction = 1
			horizontalMove = false
		}
		if (key == graphics.KeyS || key == graphics.KeyDown) && action == graphics.Press {
			if !horizontalMove && direction == 1 {
				return
			}
			direction = -1
			horizontalMove = false
		}
		if (key == graphics.KeyA || key == graphics.KeyLeft) && action == graphics.Press {
			if horizontalMove && direction == 1 {
				return
			}
			direction = -1
			horizontalMove = true
		}
		if (key == graphics.KeyD || key == graphics.KeyRight) && action == graphics.Press {
			if horizontalMove && direction == -1 {
				return
			}
			direction = 1
			horizontalMove = true
		}
	}

	if showLevel && !startGame {
		if key == graphics.KeyEnter && action == graphics.Press {
			showLevel = false
		}
	}

	if gameOver && !startGame {
		if key == graphics.KeyR && action == graphics.Press {
			gameOver = false
			startLevel = true
		}
		if key == graphics.KeyL && action == graphics.Press {
			gameOver = false
		}
	}

	if !gameOver && !showLevel {
		if key == graphics.KeySpace && action == graphics.Press {
			pauseGame = !pauseGame
		}
	}

	if startGame {
		if key == graphics.KeyEnter && action == graphics.Press {
			startGame = false
		}
		if key == graphics.KeyL && action == graphics.Press {
			startGame = false
			loadLevel = true
		}
	}
}

func resizeWindowCallback(width, height int) (startX, startY, newWidth, newHeight int32) {
	length := int32(math.Min(float64(width), float64(height)))
	startX = int32((width - int(length)) / 2)
	startY = int32((height - int(length)) / 2)
	newWidth = length
	newHeight = length
	return
}

func resetGame(level int, snakeLength int) {
	timeToMove = false
	direction = fromStart
	horizontalMove = true
	gameLevel = level

	intersectionThreshold = 1 - timeWindow
	higherEdge = float32(cellsNumber-1) + timeWindow
	lowerEdge = float32(0) - timeWindow

	snake = snakemodule.InitSnake(snakeLength, intersectionThreshold)
	setFoodPosition(fieldCells)
	eatenFoodCounter = 0

	startTime = glfw.GetTime()
}

func setFoodPosition(fieldCells []int) {
	possibleCells := snakemodule.GetPossibleCells(snake, fieldCells)
	food.SetPosition(possibleCells)
}

func saveProgress(level int) {
	currentLevel := fmt.Sprintf("%d", level)
	err := os.WriteFile("progress.txt", []byte(currentLevel), 0644)
	if err != nil {
		panic(err)
	}
}

func loadProgress() int {
	levelStr, err := os.ReadFile("progress.txt")
	if err != nil {
		panic(err)
	}
	level, err := strconv.Atoi(string(levelStr))
	if err != nil {
		panic(err)
	}
	return level
}

func getTimeWindow(level int) float32 {
	res := float32(0.5) - float32(level)/10
	return res
}

func getFoodLimit(level int) int {
	limit := 15 + 5*level
	return limit
}
