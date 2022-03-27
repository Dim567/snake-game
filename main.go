package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"snakegame/helpers"
	"snakegame/snakemodule"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Texture struct {
	id uint32
}

func (texture *Texture) Load(imgPath string) {
	imgBytes, width, height := helpers.LoadImage(imgPath)
	imgBytes = helpers.ReflectImageVertically(imgBytes, width, true)
	gl.GenTextures(1, &texture.id)
	gl.BindTexture(gl.TEXTURE_2D, texture.id)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(imgBytes))
	gl.GenerateMipmap(gl.TEXTURE_2D)
}

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

var pauseGame = false
var gameLevel = 0
var eatenFoodCounter = 0
var progressSaved = false
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

const (
	vertexShaderSource = `
	#version 410
    layout (location = 0) in vec3 aPos;
	layout (location = 1) in vec2 aTexCoord;

	out vec2 texCoord;

	uniform mat4 transformMatrix;

    void main()
    {
       gl_Position = transformMatrix*vec4(aPos.x, aPos.y, aPos.z, 1.0);
	   texCoord = aTexCoord;
    }
	` + "\x00"

	fragmentShaderSource = `
	#version 410
	in vec2 texCoord;

	out vec4 FragmentColor;

	uniform sampler2D texture1;

	void main() {
		// FragmentColor=vec4(1.0, 0.0, 0.0, 1.0);	
		FragmentColor=texture(texture1, texCoord);
	}
	` + "\x00"
)

func main() {
	runtime.LockOSThread()

	defer glfw.Terminate()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Snake game", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(framebufferSizeCallback)
	window.SetKeyCallback(keyInputCallback)

	// init gl
	err = gl.Init()
	if err != nil {
		panic(err)
	}

	// create vertex shader
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	source, free := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, source, nil)
	free()
	gl.CompileShader(vertexShader)
	var status int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vertexShader, logLength, nil, gl.Str(log))

		fmt.Errorf("failed to compile %v: %v", source, log)
	}

	// create fragment shader
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	source, free = gl.Strs(fragmentShaderSource)
	gl.ShaderSource(fragmentShader, 1, source, nil)
	free()
	gl.CompileShader(fragmentShader)
	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragmentShader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fragmentShader, logLength, nil, gl.Str(log))

		fmt.Errorf("failed to compile %v: %v", source, log)
	}

	// create program
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	var linkStatus int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &linkStatus)
	if linkStatus == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		fmt.Errorf("failed to link programm: %v", log)
	}

	// Create buffers
	var vertexArrayObject uint32
	var vertexBufferObject uint32

	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindVertexArray(vertexArrayObject)

	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, uintptr(12))
	gl.EnableVertexAttribArray(1)

	// Create and load textures
	var snakeTexture Texture
	snakeTexture.Load("snake-skin1.png")

	var backgroundTexture Texture
	backgroundTexture.Load("background.png")

	var gameOverTexture Texture
	gameOverTexture.Load("game-over.png")

	var levelTexture0 Texture
	levelTexture0.Load("level0.png")

	var levelTexture1 Texture
	levelTexture1.Load("level1.png")

	var levelTexture2 Texture
	levelTexture2.Load("level2.png")

	var levelTexture3 Texture
	levelTexture3.Load("level3.png")

	var finishLevelTexture Texture
	finishLevelTexture.Load("finish_level.png")

	levelTextures := [4]uint32{
		levelTexture0.id,
		levelTexture1.id,
		levelTexture2.id,
		// levelTexture3.id,
		finishLevelTexture.id,
	}

	for i := 0; i < len(fieldCells); i++ {
		fieldCells[i] = i
	}

	resetGame(0, 3)

	// main loop
	for !window.ShouldClose() {
		gl.ClearColor(0.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()

		switch {
		case gameOver:
			showLevel = true
			drawBackground(program, vertexArrayObject, gameOverTexture.id)
		case showLevel:
			if !progressSaved {
				saveProgress(gameLevel)
				progressSaved = true
			}
			if startLevel {
				gameLevel = 0
				timeWindow = 0.5
			}
			startLevel = false
			resetLevel = true
			textureItem := levelTextures[gameLevel]
			drawBackground(program, vertexArrayObject, textureItem)
		case pauseGame:
			timeToMove = false
			startTime = glfw.GetTime()
			drawBackground(program, vertexArrayObject, backgroundTexture.id)
			food.Draw(program, vertexArrayObject, snakeTexture.id, drawObject)
			snake.Draw(program, vertexArrayObject, snakeTexture.id, drawObject)
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
				gameLevel == 3 {
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
				if eatenFoodCounter == 3 {
					eatenFoodCounter = 0
					timeWindow -= 0.1
					gameLevel += 1
					progressSaved = false
					showLevel = true
				} else {
					setFoodPosition(fieldCells)
				}
			}

			drawBackground(program, vertexArrayObject, backgroundTexture.id)

			if period < (2*timeWindow/7) || period > (5*timeWindow/7) {
				food.Draw(program, vertexArrayObject, snakeTexture.id, drawObject)
			}
			snake.Draw(program, vertexArrayObject, snakeTexture.id, drawObject)
		}

		window.SwapBuffers()
		glfw.SwapInterval(1)
	}
}

func drawObject(program, vertexArrayObject, texture uint32, vec mgl32.Vec2) {
	scaleFactor := float32(2.0 / cellsNumber)
	scale := mgl32.Scale3D(scaleFactor, scaleFactor, 1)
	xPos := vec.X()*scaleFactor - 1
	yPos := vec.Y()*scaleFactor - 1
	translate := mgl32.Translate3D(xPos, yPos, 0)
	transform := translate.Mul4(scale)
	draw(program, vertexArrayObject, texture, transform)
}

func drawBackground(program, vertexArrayObject, texture uint32) {
	scale := mgl32.Scale3D(2, 2, 1)
	translate := mgl32.Translate3D(-1, -1, 0)
	transform := translate.Mul4(scale)
	draw(program, vertexArrayObject, texture, transform)
}

func draw(program, vertexArrayObject, texture uint32, transform mgl32.Mat4) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transformMatrix\x00")), 1, false, &transform[0])
	gl.UseProgram(program)
	gl.BindVertexArray(vertexArrayObject)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)
}

func keyInputCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if !pauseGame && !gameOver && !showLevel {
		if key == glfw.KeyUp && action == glfw.Press {
			if !horizontalMove && direction == -1 {
				return
			}
			direction = 1
			horizontalMove = false
		}
		if key == glfw.KeyDown && action == glfw.Press {
			if !horizontalMove && direction == 1 {
				return
			}
			direction = -1
			horizontalMove = false
		}
		if key == glfw.KeyLeft && action == glfw.Press {
			if horizontalMove && direction == 1 {
				return
			}
			direction = -1
			horizontalMove = true
		}
		if key == glfw.KeyRight && action == glfw.Press {
			if horizontalMove && direction == -1 {
				return
			}
			direction = 1
			horizontalMove = true
		}
	}

	if showLevel {
		if key == glfw.KeyEnter && action == glfw.Press {
			showLevel = false
		}
	}

	if gameOver {
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
}

func framebufferSizeCallback(window *glfw.Window, width, height int) {
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

	startTime = glfw.GetTime()
}

func setFoodPosition(fieldCells []int) {
	possibleCells := snakemodule.GetPossibleCells(snake, fieldCells)
	food.SetPosition(possibleCells)
}

func saveProgress(level int) {
	currentLevel := fmt.Sprintf("%d", level)
	os.WriteFile("progress.txt", []byte(currentLevel), 0644)
}
