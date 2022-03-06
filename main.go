package main

import (
	"fmt"
	"math"
	"runtime"
	"snakegame/snakemodule"
	"strings"

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

var changeFoodPosition = true
var shouldMove = true

const (
	fromStart int8 = 1
	toStart   int8 = -1
)

// Movement direction
var direction = fromStart

// Movement horizontal or vertical
var horizontalMove = true

// mouse positions
// var firstMouse = true
// var lastX = float32(900)
// var lastY = float32(500)

// velocity settings
// var deltaTime = float32(0)
// var lastFrame = float32(0)

// prepare vertices
var vertices = []float32{
	0, 1, 0.0, // top left
	1, 1, 0.0, // top right
	1, 0, 0.0, // bottom right

	0, 1, 0.0, // top left
	1, 0, 0.0, // bottom right
	0, 0, 0.0, // bottom left
}

const (
	vertexShaderSource = `
	#version 410
    layout (location = 0) in vec3 aPos;

	uniform mat4 transformMatrix;

    void main()
    {
       gl_Position = transformMatrix*vec4(aPos.x, aPos.y, aPos.z, 1.0);
    }
	` + "\x00"

	fragmentShaderSource = `
	#version 410
	out vec4 FragmentColor;

	void main() {
		FragmentColor=vec4(1.0, 0.0, 0.0, 1.0);	
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

	// create buffers
	var vertexArrayObject uint32
	var vertexBufferObject uint32

	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindVertexArray(vertexArrayObject)

	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)

	// xOffset := float32(0)
	// yOffset := float32(0)
	startTime := glfw.GetTime()
	higherEdge := float32(cellsNumber-1) + 0.5
	lowerEdge := float32(0) - 0.5

	snake := snakemodule.InitSnake(3)

	var food snakemodule.Food

	fieldCells := make([]int, cellsNumber*cellsNumber)
	for i := 0; i < len(fieldCells); i++ {
		fieldCells[i] = i
	}

	timeToMove := false
	var xLimit, yLimit float32

	//////////////////////////////////////////////////////////////////
	// main loop
	for !window.ShouldClose() {
		endTime := glfw.GetTime()
		period := endTime - startTime

		if period >= 0.5 {
			startTime = endTime
			timeToMove = true
		}

		snakeHead := snake.GetHead()
		x := snakeHead.GetCoords().X()
		y := snakeHead.GetCoords().Y()

		if timeToMove {
			if horizontalMove {
				xLimit = x + float32(direction)*float32(period)
			} else {
				yLimit = y + float32(direction)*float32(period)
			}
			if xLimit >= higherEdge || xLimit <= lowerEdge || yLimit >= higherEdge || yLimit <= lowerEdge {
				fmt.Println(xLimit)
				shouldMove = false
			}
		}

		if shouldMove && timeToMove {
			timeToMove = false
			if horizontalMove {
				xLimit = x + float32(direction)*float32(period)
				x += float32(direction)
			} else {
				yLimit = y + float32(direction)*float32(period)
				y += float32(direction)
			}

		}

		if changeFoodPosition {
			possibleCells := snakemodule.GetPossibleCells(snake, fieldCells)
			food.SetPosition(possibleCells)
			changeFoodPosition = false
		}

		snake.Eat(&food, &changeFoodPosition)
		snake.Move(mgl32.Vec2{x, y})

		processInput(window)
		gl.ClearColor(0.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		food.Draw(program, vertexArrayObject, drawObject)
		snake.Draw(program, vertexArrayObject, drawObject)
		// drawObject(program, vertexArrayObject, xOffset, yOffset)

		glfw.PollEvents()
		window.SwapBuffers()
		glfw.SwapInterval(1)
	}
}

func drawObject(program, vertexArrayObject uint32, vec mgl32.Vec2) {
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

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		direction = 1
		horizontalMove = false
	}
	if window.GetKey(glfw.KeyDown) == glfw.Press {
		direction = -1
		horizontalMove = false
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		direction = -1
		horizontalMove = true
	}
	if window.GetKey(glfw.KeyRight) == glfw.Press {
		direction = 1
		horizontalMove = true
	}
}

func framebufferSizeCallback(window *glfw.Window, width, height int) {
	length := int32(math.Min(float64(width), float64(height)))
	startX := int32((width - int(length)) / 2)
	startY := int32((height - int(length)) / 2)
	gl.Viewport(startX, startY, length, length)
}
