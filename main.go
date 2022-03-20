package main

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"os"
	"runtime"
	"snakegame/snakemodule"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	_ "image/jpeg"
	_ "image/png"
)

type Texture struct {
	id uint32
}

func (texture *Texture) Load(imgPath string) {
	imgBytes, width, height := loadImage(imgPath)
	imgBytes = reflectImageVertically(imgBytes, width, true)
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

// var restartGame = true
var pauseGame = false
var gameLevel = 0
var eatenFoodCounter = 0

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
var levelStartTime float64

// Speed settings
var timeWindow float32
var intersectionThreshold float32
var higherEdge, lowerEdge float32
var timeToMove = false

// Field settings
var fieldCells = make([]int, cellsNumber*cellsNumber)

var snake *snakemodule.Snake
var food snakemodule.Food

var showLevel = true

// Movement horizontal or vertical
var horizontalMove = true

var foodWasEaten = false

// mouse positions
// var firstMouse = true
// var lastX = float32(900)
// var lastY = float32(500)

// velocity settings
// var deltaTime = float32(0)
// var lastFrame = float32(0)

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

	// init glfw
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
	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

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

	levelTextures := [4]uint32{
		levelTexture0.id,
		levelTexture1.id,
		levelTexture2.id,
		levelTexture3.id,
	}

	for i := 0; i < len(fieldCells); i++ {
		fieldCells[i] = i
	}

	resetGame(0, 0.5, 3)

	//////////////////////////////////////////////////////////////////
	// main loop
	for !window.ShouldClose() {
		if pauseGame {
			period = 0
			timeToMove = false
			startTime = glfw.GetTime()
		}

		if gameLevel == 4 {
			gameOver = true
		}

		if eatenFoodCounter == 10 {
			newLevel := gameLevel + 1
			newTimeWindow := timeWindow - 0.1
			resetGame(newLevel, newTimeWindow, 3)
		}

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

		// frontX, frontY = snake.GetFront().Elem()
		if frontX >= higherEdge ||
			frontX <= lowerEdge ||
			frontY >= higherEdge ||
			frontY <= lowerEdge ||
			snake.CheckIntersection() { // move this into SetFront
			gameOver = true
		}

		gl.ClearColor(0.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		if gameOver {
			drawBackground(program, vertexArrayObject, gameOverTexture.id)
		} else if showLevel {
			textureItem := levelTextures[gameLevel]
			drawBackground(program, vertexArrayObject, textureItem)
			timer := glfw.GetTime() - levelStartTime
			if timer > 1 {
				showLevel = false
			}
		} else {
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
				eatenFoodCounter += 1
				setFoodPosition(fieldCells)
				foodWasEaten = false
			}

			drawBackground(program, vertexArrayObject, backgroundTexture.id)

			if period < (2*timeWindow/7) || period > (5*timeWindow/7) {
				food.Draw(program, vertexArrayObject, snakeTexture.id, drawObject)
			}
			snake.Draw(program, vertexArrayObject, snakeTexture.id, drawObject)
		}

		glfw.PollEvents()
		window.SwapBuffers()
		glfw.SwapInterval(1)
	}
}

func drawObject(program, vertexArrayObject, texture uint32, vec mgl32.Vec2) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	scaleFactor := float32(2.0 / cellsNumber)
	scale := mgl32.Scale3D(scaleFactor, scaleFactor, 1)
	xPos := vec.X()*scaleFactor - 1
	yPos := vec.Y()*scaleFactor - 1
	translate := mgl32.Translate3D(xPos, yPos, 0)
	transform := translate.Mul4(scale)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transformMatrix\x00")), 1, false, &transform[0])
	gl.UseProgram(program)
	gl.BindVertexArray(vertexArrayObject)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)
}

func drawBackground(program, vertexArrayObject, texture uint32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	scale := mgl32.Scale3D(2, 2, 1)
	translate := mgl32.Translate3D(-1, -1, 0)
	transform := translate.Mul4(scale)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transformMatrix\x00")), 1, false, &transform[0])
	gl.UseProgram(program)
	gl.BindVertexArray(vertexArrayObject)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)
}

func keyInputCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if !pauseGame {
		if key == glfw.KeyUp && action == glfw.Press {
			// fmt.Println("UP")
			if !horizontalMove && direction == -1 {
				return
			}
			direction = 1
			horizontalMove = false
		}
		if key == glfw.KeyDown && action == glfw.Press {
			// fmt.Println("DOWN")
			if !horizontalMove && direction == 1 {
				return
			}
			direction = -1
			horizontalMove = false
		}
		if key == glfw.KeyLeft && action == glfw.Press {
			// fmt.Println("LEFT")
			if horizontalMove && direction == 1 {
				return
			}
			direction = -1
			horizontalMove = true
		}
		if key == glfw.KeyRight && action == glfw.Press {
			// fmt.Println("RIGHT")
			if horizontalMove && direction == -1 {
				return
			}
			direction = 1
			horizontalMove = true
		}
	}

	if key == glfw.KeyR && action == glfw.Press {
		resetGame(0, 0.5, 3)
	}

	if key == glfw.KeySpace && action == glfw.Press {
		// fmt.Println("PAUSE")
		pauseGame = !pauseGame
	}
}

func framebufferSizeCallback(window *glfw.Window, width, height int) {
	length := int32(math.Min(float64(width), float64(height)))
	startX := int32((width - int(length)) / 2)
	startY := int32((height - int(length)) / 2)
	gl.Viewport(startX, startY, length, length)
}

func loadImage(path string) ([]uint8, int32, int32) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	return rgba.Pix, width, height
}

func reflectImageVertically(imageData []uint8, width int32, alfa bool) []uint8 {
	reflected := make([]uint8, 0, len(imageData))
	var stride int
	if alfa {
		stride = int(width * 4)
	} else {
		stride = int(width * 3)
	}

	for i := len(imageData) - stride; i >= 0; i = i - stride {
		for j := i; j < stride+i; j++ {
			reflected = append(reflected, imageData[j])
		}
	}
	return reflected
}

func resetGame(level int, period float32, snakeLength int) {
	pauseGame = false
	gameOver = false
	timeToMove = false
	direction = fromStart
	horizontalMove = true
	eatenFoodCounter = 0
	gameLevel = level
	showLevel = true

	timeWindow = period
	intersectionThreshold = 1 - period
	fmt.Println(intersectionThreshold)
	higherEdge = float32(cellsNumber-1) + period
	lowerEdge = float32(0) - period

	snake = snakemodule.InitSnake(snakeLength, intersectionThreshold)
	setFoodPosition(fieldCells)

	startTime = glfw.GetTime()
	levelStartTime = glfw.GetTime()
}

func setFoodPosition(fieldCells []int) {
	possibleCells := snakemodule.GetPossibleCells(snake, fieldCells)
	food.SetPosition(possibleCells)
}
