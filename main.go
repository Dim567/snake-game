package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"snakegame/graphics"
	"snakegame/snakemodule"
	"strconv"

	"github.com/go-gl/gl/v4.1-core/gl"
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

// prepare vertices
var vertices = []float32{
	//vertices coords              texture coords
	0, 1, 0.0 /* top left */, 0.0, 1.0,
	1, 1, 0.0 /* top right */, 1.0, 1.0,
	1, 0, 0.0 /* bottom right */, 1.0, 0.0,

	0, 1, 0.0 /* top left */, 0.0, 1.0,
	1, 0, 0.0 /* bottom right */, 1.0, 0.0,
	0, 0, 0.0 /* bottom left */, 0.0, 0.0,
}

func main() {
	runtime.LockOSThread()

	defer graphics.Terminate()
	err := graphics.Init("Snake game", windowWidth, windowHeight)
	if err != nil {
		panic(err)
	}
	graphics.SetResizeWindowCallback(resizeWindowCallback)
	graphics.SetKeyInputCallback(keyInputCallback)

	program, err := graphics.CreateProgram()
	if err != nil {
		panic(err)
	}

	vertexArrayObject := graphics.CreateVAO(vertices)

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
	cb := func() {
		switch {
		case startGame:
			showLevel = true
			drawBackground(program, vertexArrayObject, startGameTexture)
		case gameOver:
			showLevel = true
			drawBackground(program, vertexArrayObject, gameOverTexture)
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
			drawBackground(program, vertexArrayObject, textureItem)
		case pauseGame:
			timeToMove = false
			startTime = glfw.GetTime()
			drawBackground(program, vertexArrayObject, backgroundTexture)
			food.Draw(program, vertexArrayObject, snakeTexture, drawObject)
			snake.Draw(program, vertexArrayObject, snakeTexture, drawObject)
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

				// need to refactor this
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

			drawBackground(program, vertexArrayObject, backgroundTexture)

			if period < (2*timeWindow/7) || period > (5*timeWindow/7) {
				food.Draw(program, vertexArrayObject, snakeTexture, drawObject)
			}
			snake.Draw(program, vertexArrayObject, snakeTexture, drawObject)
		}
	}

	graphics.MainLoop(cb)
}

func drawObject(program, vertexArrayObject, texture uint32, vec mgl32.Vec2) {
	scaleFactor := float32(2.0 / cellsNumber)
	scale := mgl32.Scale3D(scaleFactor, scaleFactor, 1)
	xPos := vec.X()*scaleFactor - 1
	yPos := vec.Y()*scaleFactor - 1
	translate := mgl32.Translate3D(xPos, yPos, 0)
	transform := translate.Mul4(scale)
	graphics.Draw(program, vertexArrayObject, texture, transform)
}

func drawBackground(program, vertexArrayObject, texture uint32) {
	scale := mgl32.Scale3D(2, 2, 1)
	translate := mgl32.Translate3D(-1, -1, 0)
	transform := translate.Mul4(scale)
	graphics.Draw(program, vertexArrayObject, texture, transform)
}

func keyInputCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if !pauseGame && !gameOver && !showLevel {
		if (key == glfw.KeyW || key == glfw.KeyUp) && action == glfw.Press {
			if !horizontalMove && direction == -1 {
				return
			}
			direction = 1
			horizontalMove = false
		}
		if (key == glfw.KeyS || key == glfw.KeyDown) && action == glfw.Press {
			if !horizontalMove && direction == 1 {
				return
			}
			direction = -1
			horizontalMove = false
		}
		if (key == glfw.KeyA || key == glfw.KeyLeft) && action == glfw.Press {
			if horizontalMove && direction == 1 {
				return
			}
			direction = -1
			horizontalMove = true
		}
		if (key == glfw.KeyD || key == glfw.KeyRight) && action == glfw.Press {
			if horizontalMove && direction == -1 {
				return
			}
			direction = 1
			horizontalMove = true
		}
	}

	if showLevel && !startGame {
		if key == glfw.KeyEnter && action == glfw.Press {
			showLevel = false
		}
	}

	if gameOver && !startGame {
		if key == glfw.KeyR && action == glfw.Press {
			gameOver = false
			startLevel = true
		}
		if key == glfw.KeyL && action == glfw.Press {
			gameOver = false
		}
	}

	if !gameOver && !showLevel {
		if key == glfw.KeySpace && action == glfw.Press {
			pauseGame = !pauseGame
		}
	}

	if startGame {
		if key == glfw.KeyEnter && action == glfw.Press {
			startGame = false
		}
		if key == glfw.KeyL && action == glfw.Press {
			startGame = false
			loadLevel = true
		}
	}
}

func resizeWindowCallback(window *glfw.Window, width, height int) {
	length := int32(math.Min(float64(width), float64(height)))
	startX := int32((width - int(length)) / 2)
	startY := int32((height - int(length)) / 2)
	gl.Viewport(startX, startY, length, length)
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
