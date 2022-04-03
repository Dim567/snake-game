package graphics

import (
	"fmt"
	"snakegame/helpers"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type KeyAction glfw.Action

const (
	Release KeyAction = KeyAction(glfw.Release)
	Press   KeyAction = KeyAction(glfw.Press)
	Repeat  KeyAction = KeyAction(glfw.Repeat)
)

type KeyValue glfw.Key

const (
	KeyUnknown      KeyValue = KeyValue(glfw.KeyUnknown)
	KeySpace        KeyValue = KeyValue(glfw.KeySpace)
	KeyApostrophe   KeyValue = KeyValue(glfw.KeyApostrophe)
	KeyComma        KeyValue = KeyValue(glfw.KeyComma)
	KeyMinus        KeyValue = KeyValue(glfw.KeyMinus)
	KeyPeriod       KeyValue = KeyValue(glfw.KeyPeriod)
	KeySlash        KeyValue = KeyValue(glfw.KeySlash)
	Key0            KeyValue = KeyValue(glfw.Key0)
	Key1            KeyValue = KeyValue(glfw.Key1)
	Key2            KeyValue = KeyValue(glfw.Key2)
	Key3            KeyValue = KeyValue(glfw.Key3)
	Key4            KeyValue = KeyValue(glfw.Key4)
	Key5            KeyValue = KeyValue(glfw.Key5)
	Key6            KeyValue = KeyValue(glfw.Key6)
	Key7            KeyValue = KeyValue(glfw.Key7)
	Key8            KeyValue = KeyValue(glfw.Key8)
	Key9            KeyValue = KeyValue(glfw.Key9)
	KeySemicolon    KeyValue = KeyValue(glfw.KeySemicolon)
	KeyEqual        KeyValue = KeyValue(glfw.KeyEqual)
	KeyA            KeyValue = KeyValue(glfw.KeyA)
	KeyB            KeyValue = KeyValue(glfw.KeyB)
	KeyC            KeyValue = KeyValue(glfw.KeyC)
	KeyD            KeyValue = KeyValue(glfw.KeyD)
	KeyE            KeyValue = KeyValue(glfw.KeyE)
	KeyF            KeyValue = KeyValue(glfw.KeyF)
	KeyG            KeyValue = KeyValue(glfw.KeyG)
	KeyH            KeyValue = KeyValue(glfw.KeyH)
	KeyI            KeyValue = KeyValue(glfw.KeyI)
	KeyJ            KeyValue = KeyValue(glfw.KeyJ)
	KeyK            KeyValue = KeyValue(glfw.KeyK)
	KeyL            KeyValue = KeyValue(glfw.KeyL)
	KeyM            KeyValue = KeyValue(glfw.KeyM)
	KeyN            KeyValue = KeyValue(glfw.KeyN)
	KeyO            KeyValue = KeyValue(glfw.KeyO)
	KeyP            KeyValue = KeyValue(glfw.KeyP)
	KeyQ            KeyValue = KeyValue(glfw.KeyQ)
	KeyR            KeyValue = KeyValue(glfw.KeyR)
	KeyS            KeyValue = KeyValue(glfw.KeyS)
	KeyT            KeyValue = KeyValue(glfw.KeyT)
	KeyU            KeyValue = KeyValue(glfw.KeyY)
	KeyV            KeyValue = KeyValue(glfw.KeyV)
	KeyW            KeyValue = KeyValue(glfw.KeyW)
	KeyX            KeyValue = KeyValue(glfw.KeyX)
	KeyY            KeyValue = KeyValue(glfw.KeyY)
	KeyZ            KeyValue = KeyValue(glfw.KeyZ)
	KeyLeftBracket  KeyValue = KeyValue(glfw.KeyLeftBracket)
	KeyBackslash    KeyValue = KeyValue(glfw.KeyBackslash)
	KeyRightBracket KeyValue = KeyValue(glfw.KeyRightBracket)
	KeyGraveAccent  KeyValue = KeyValue(glfw.KeyGraveAccent)
	KeyWorld1       KeyValue = KeyValue(glfw.KeyWorld1)
	KeyWorld2       KeyValue = KeyValue(glfw.KeyWorld2)
	KeyEscape       KeyValue = KeyValue(glfw.KeyEscape)
	KeyEnter        KeyValue = KeyValue(glfw.KeyEnter)
	KeyTab          KeyValue = KeyValue(glfw.KeyTab)
	KeyBackspace    KeyValue = KeyValue(glfw.KeyBackspace)
	KeyInsert       KeyValue = KeyValue(glfw.KeyInsert)
	KeyDelete       KeyValue = KeyValue(glfw.KeyDelete)
	KeyRight        KeyValue = KeyValue(glfw.KeyRight)
	KeyLeft         KeyValue = KeyValue(glfw.KeyLeft)
	KeyDown         KeyValue = KeyValue(glfw.KeyDown)
	KeyUp           KeyValue = KeyValue(glfw.KeyUp)
	KeyPageUp       KeyValue = KeyValue(glfw.KeyPageUp)
	KeyPageDown     KeyValue = KeyValue(glfw.KeyPageDown)
	KeyHome         KeyValue = KeyValue(glfw.KeyHome)
	KeyEnd          KeyValue = KeyValue(glfw.KeyEnd)
	KeyCapsLock     KeyValue = KeyValue(glfw.KeyCapsLock)
	KeyScrollLock   KeyValue = KeyValue(glfw.KeyScrollLock)
	KeyNumLock      KeyValue = KeyValue(glfw.KeyNumLock)
	KeyPrintScreen  KeyValue = KeyValue(glfw.KeyPrintScreen)
	KeyPause        KeyValue = KeyValue(glfw.KeyPause)
	KeyF1           KeyValue = KeyValue(glfw.KeyF1)
	KeyF2           KeyValue = KeyValue(glfw.KeyF2)
	KeyF3           KeyValue = KeyValue(glfw.KeyF3)
	KeyF4           KeyValue = KeyValue(glfw.KeyF4)
	KeyF5           KeyValue = KeyValue(glfw.KeyF5)
	KeyF6           KeyValue = KeyValue(glfw.KeyF6)
	KeyF7           KeyValue = KeyValue(glfw.KeyF7)
	KeyF8           KeyValue = KeyValue(glfw.KeyF8)
	KeyF9           KeyValue = KeyValue(glfw.KeyF9)
	KeyF10          KeyValue = KeyValue(glfw.KeyF10)
	KeyF11          KeyValue = KeyValue(glfw.KeyF11)
	KeyF12          KeyValue = KeyValue(glfw.KeyF12)
	KeyF13          KeyValue = KeyValue(glfw.KeyF13)
	KeyF14          KeyValue = KeyValue(glfw.KeyF14)
	KeyF15          KeyValue = KeyValue(glfw.KeyF15)
	KeyF16          KeyValue = KeyValue(glfw.KeyF16)
	KeyF17          KeyValue = KeyValue(glfw.KeyF17)
	KeyF18          KeyValue = KeyValue(glfw.KeyF18)
	KeyF19          KeyValue = KeyValue(glfw.KeyF19)
	KeyF20          KeyValue = KeyValue(glfw.KeyF20)
	KeyF21          KeyValue = KeyValue(glfw.KeyF21)
	KeyF22          KeyValue = KeyValue(glfw.KeyF22)
	KeyF23          KeyValue = KeyValue(glfw.KeyF23)
	KeyF24          KeyValue = KeyValue(glfw.KeyF24)
	KeyF25          KeyValue = KeyValue(glfw.KeyF25)
	KeyKP0          KeyValue = KeyValue(glfw.KeyKP0)
	KeyKP1          KeyValue = KeyValue(glfw.KeyKP1)
	KeyKP2          KeyValue = KeyValue(glfw.KeyKP2)
	KeyKP3          KeyValue = KeyValue(glfw.KeyKP3)
	KeyKP4          KeyValue = KeyValue(glfw.KeyKP4)
	KeyKP5          KeyValue = KeyValue(glfw.KeyKP5)
	KeyKP6          KeyValue = KeyValue(glfw.KeyKP6)
	KeyKP7          KeyValue = KeyValue(glfw.KeyKP7)
	KeyKP8          KeyValue = KeyValue(glfw.KeyKP8)
	KeyKP9          KeyValue = KeyValue(glfw.KeyKP9)
	KeyKPDecimal    KeyValue = KeyValue(glfw.KeyKPDecimal)
	KeyKPDivide     KeyValue = KeyValue(glfw.KeyKPDivide)
	KeyKPMultiply   KeyValue = KeyValue(glfw.KeyKPMultiply)
	KeyKPSubtract   KeyValue = KeyValue(glfw.KeyKPSubtract)
	KeyKPAdd        KeyValue = KeyValue(glfw.KeyKPAdd)
	KeyKPEnter      KeyValue = KeyValue(glfw.KeyKPEnter)
	KeyKPEqual      KeyValue = KeyValue(glfw.KeyKPEqual)
	KeyLeftShift    KeyValue = KeyValue(glfw.KeyLeftShift)
	KeyLeftControl  KeyValue = KeyValue(glfw.KeyLeftControl)
	KeyLeftAlt      KeyValue = KeyValue(glfw.KeyLeftAlt)
	KeyLeftSuper    KeyValue = KeyValue(glfw.KeyLeftSuper)
	KeyRightShift   KeyValue = KeyValue(glfw.KeyRightShift)
	KeyRightControl KeyValue = KeyValue(glfw.KeyRightControl)
	KeyRightAlt     KeyValue = KeyValue(glfw.KeyRightAlt)
	KeyRightSuper   KeyValue = KeyValue(glfw.KeyRightSuper)
	KeyMenu         KeyValue = KeyValue(glfw.KeyMenu)
	KeyLast         KeyValue = KeyValue(glfw.KeyLast)
)

type ShaderType int

const (
	Vertex ShaderType = iota
	Fragment
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

	vertexShader, err := createShader(Vertex, vertexShaderSource)
	if err != nil {
		return err
	}

	fragmentShader, err := createShader(Fragment, fragmentShaderSource)
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

func SetKeyInputCallback(callback func(keyValue KeyValue, keyAction KeyAction)) {
	keyInputCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		callback(KeyValue(key), KeyAction(action))
	}
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

func createShader(shaderType ShaderType, shaderSource string) (uint32, error) {
	var shType uint32
	if shaderType == Vertex {
		shType = gl.VERTEX_SHADER
	}
	if shaderType == Fragment {
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
