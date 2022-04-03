package graphics

import (
	"fmt"
	"snakegame/helpers"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var window *glfw.Window
var program, vertexArrayObject uint32
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
		FragmentColor=texture(texture1, texCoord);
	}
	` + "\x00"
)

func Init(windowName string, windowWidth, windowHeight int) error {
	err := glfw.Init()
	if err != nil {
		return err
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err = glfw.CreateWindow(windowWidth, windowHeight, windowName, nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	err = gl.Init()
	if err != nil {
		return err
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	vertexShader, err := createShader("vertex", vertexShaderSource)
	if err != nil {
		return err
	}

	fragmentShader, err := createShader("fragment", fragmentShaderSource)
	if err != nil {
		return err
	}
	program, err = createProgram(vertexShader, fragmentShader)
	if err != nil {
		return err
	}

	vertexArrayObject = createVAO(vertices)

	return nil
}

func Terminate() {
	glfw.Terminate()
}

func SetResizeWindowCallback(callback func(width, height int) (startX, startY, newWidth, newHeight int32)) {
	framebufferSizeCallback := func(w *glfw.Window, width int, height int) {
		startX, startY, nWidth, nHeight := callback(width, height)
		gl.Viewport(startX, startY, nWidth, nHeight)
	}
	window.SetFramebufferSizeCallback(framebufferSizeCallback)
}

func SetKeyInputCallback(keyInputCallback glfw.KeyCallback) {
	window.SetKeyCallback(keyInputCallback)
}

func MainLoop(gameLogic func()) {
	for !window.ShouldClose() {
		gl.ClearColor(0.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gameLogic()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func createShader(shaderType, shaderSource string) (uint32, error) {
	var shType uint32
	if shaderType == "vertex" {
		shType = gl.VERTEX_SHADER
	}
	if shaderType == "fragment" {
		shType = gl.FRAGMENT_SHADER
	}
	shader := gl.CreateShader(shType)
	source, free := gl.Strs(shaderSource)
	gl.ShaderSource(shader, 1, source, nil)
	free()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
}

func createProgram(vertexShader, fragmentShader uint32) (uint32, error) {
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
		return 0, fmt.Errorf("failed to link programm: %v", log)
	}
	return program, nil
}

func createVAO(vertices []float32) uint32 {
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
	return vertexArrayObject
}

func LoadTexture(imgPath string) uint32 {
	var texture uint32
	imgBytes, width, height := helpers.LoadImage(imgPath)
	imgBytes = helpers.ReflectImageVertically(imgBytes, width, true)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(imgBytes))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return texture
}

func Draw(texture uint32, transform mgl32.Mat4) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transformMatrix\x00")), 1, false, &transform[0])
	gl.UseProgram(program)
	gl.BindVertexArray(vertexArrayObject)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)
}
