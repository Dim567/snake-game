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

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, uintptr(12))
	gl.EnableVertexAttribArray(1)

	// load image for textures
	imgBytes, imgWidth, imgHeight := loadImage("awesomeface.png")
	imgBytes = reflectImageVertically(imgBytes, imgWidth, true)

	// create texture
	var texture1 uint32
	gl.GenTextures(1, &texture1)
	gl.BindTexture(gl.TEXTURE_2D, texture1)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, imgWidth, imgHeight, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(imgBytes))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	var food snakemodule.Food

	fieldCells := make([]int, cellsNumber*cellsNumber)
	for i := 0; i < len(fieldCells); i++ {
		fieldCells[i] = i
	}

	timeToMove := false
	timeWindow := float32(0.4)
	intersectionThreshold := 1 - timeWindow
	higherEdge := float32(cellsNumber-1) + timeWindow
	lowerEdge := float32(0) - timeWindow

	snake := snakemodule.InitSnake(3, intersectionThreshold)

	startTime := glfw.GetTime()
	//////////////////////////////////////////////////////////////////
	// main loop
	for !window.ShouldClose() {
		endTime := glfw.GetTime()
		period := float32(endTime - startTime)

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
			snake.CheckIntersection() { // move this into SetFront
			shouldMove = false
		}

		if shouldMove && timeToMove {
			timeToMove = false
			snake.Eat(&food, &changeFoodPosition)
			if horizontalMove {
				x += float32(direction)
			} else {
				y += float32(direction)
			}

			// need to refactor this
			snake.Move(mgl32.Vec2{x, y})
		}

		if changeFoodPosition {
			possibleCells := snakemodule.GetPossibleCells(snake, fieldCells)
			food.SetPosition(possibleCells)
			changeFoodPosition = false
		}

		processInput(window)
		gl.ClearColor(0.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		food.Draw(program, vertexArrayObject, texture1, drawObject)
		snake.Draw(program, vertexArrayObject, texture1, drawObject)
		// drawObject(program, vertexArrayObject, xOffset, yOffset)

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

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		if !horizontalMove && direction == -1 {
			return
		}
		direction = 1
		horizontalMove = false
	}
	if window.GetKey(glfw.KeyDown) == glfw.Press {
		if !horizontalMove && direction == 1 {
			return
		}
		direction = -1
		horizontalMove = false
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		if horizontalMove && direction == 1 {
			return
		}
		direction = -1
		horizontalMove = true
	}
	if window.GetKey(glfw.KeyRight) == glfw.Press {
		if horizontalMove && direction == -1 {
			return
		}
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
	// fmt.Print(imageData)
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
