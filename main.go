package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// widow initial sizes
const (
	windowWidth  = 800
	windowHeight = 800
)

// mouse positions
var firstMouse = true
var lastX = float32(900)
var lastY = float32(500)

// velocity settings
var deltaTime = float32(0)
var lastFrame = float32(0)

// prepare vertices
var vertices = []float32{
	-0.5, 0.5, 0.0,
	0.5, -0.5, 0.0,
	0.5, 0.5, 0.0,
	-0.5, 0.5, 0.0,
	0.5, -0.5, 0.0,
	-0.5, -0.5, 0.0,
}

const (
	vertexShaderSource = `
	#version 410
    layout (location = 0) in vec3 aPos;

    void main()
    {
       gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
    }
	` + "\x00"

	fragmentShaderSource = `
	#version 410
	out vec4 FragmentColor;

	void main() {
		FragmentColor=vec4(1.0, 1.0, 0.0, 1.0);	
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
	// if {
	// 	var success
	// 	glGetProgramiv(shader, GL_LINK_STATUS, &success);
	// 	if(!success)
	// 	{
	// 		glGetProgramInfoLog(shader, 1024, NULL, infoLog);
	// 		std::cout << "ERROR::PROGRAM_LINKING_ERROR of type: " << type << "\n" << infoLog << "\n -- --------------------------------------------------- -- " << std::endl;
	// 	}
	// }
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

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

	// main loop
	for !window.ShouldClose() {
		currentFrame := glfw.GetTime()
		deltaTime = float32(currentFrame) - lastFrame
		lastFrame = float32(currentFrame)
		processInput(window)
		gl.ClearColor(0.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		drawObject(program, vertexArrayObject)

		glfw.PollEvents()
		window.SwapBuffers()
		glfw.SwapInterval(1)
	}
}

func getAngle(velocity, time float32) float32 {
	path := velocity * time
	roundsCount := path / 360
	fullRoundsCount := float32(int(roundsCount))
	angle := (roundsCount - fullRoundsCount) * 360

	return angle
}

func drawObject(program, vertexArrayObject uint32) {
	gl.UseProgram(program)
	gl.BindVertexArray(vertexArrayObject)

	// drawing ////////////////////////////////////////
	// gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyW) == glfw.Press {
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
	}
}

func framebufferSizeCallback(window *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}
