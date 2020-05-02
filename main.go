package main

import (
	"io/ioutil"
	"log"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const width = 800
const height = 600

func main() {

	var triangle = []float32{
		//X,Y,Z,R,G,B
		-1.0, 0.0, 0.0, 1.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0, 0.0, 1.0,
		0.00, 1.0, 0.0, 0.0, 1.0, 1.0,
		0.1, 0.0, 0.0, 1.0, 1.0, 1.0,
		0.0, 0.0, 0.0, 0.5, 0.5, 0.5,
		0.05, 0.1, 0.0, 1.0, 1.0, 1.0,
	}
	/*var triangle2 = []float32{
		0.1, 0.0, 0.0,
		0.0, 0.0, 0.0,
		0.05, 0.1, 0.0,
	}*/
	if err := glfw.Init(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(width, height, "OpenGL 2 triangles", nil, nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	window.MakeContextCurrent()
	if err = gl.Init(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer glfw.Terminate()
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	tmp, err := getShaderSource("shaders/simple_shader.glsl")
	shaderSource, freeFN := gl.Strs(tmp)
	defer freeFN()
	gl.ShaderSource(vertexShader, 1, shaderSource, nil)
	if err != nil {
		log.Println(err)
	}
	tmp1, err := getShaderSource("shaders/fragment_shader.glsl")
	if err != nil {
		log.Println(err)
	}
	shaderSource, freeFNs := gl.Strs(tmp1)
	defer freeFNs()
	gl.ShaderSource(fragmentShader, 1, shaderSource, nil)
	gl.CompileShader(vertexShader)
	var succ int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &succ)
	//infolog := new(uint8)

	if succ == 0 {
		var length int32
		buffer := make([]byte, 5*1024)
		gl.GetShaderInfoLog(vertexShader, 1024, &length, (*uint8)(unsafe.Pointer(&buffer[0])))
		log.Printf(string(buffer[:]))
	}
	gl.CompileShader(fragmentShader)
	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &succ)
	//infolog := new(uint8)

	if succ == 0 {
		var length int32
		buffer := make([]byte, 5*1024)
		gl.GetShaderInfoLog(vertexShader, 1024, &length, (*uint8)(unsafe.Pointer(&buffer[0])))
		log.Printf(string(buffer[:]))
	}
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	log.Println(gl.GetError())
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triangle)*4, gl.Ptr(triangle), gl.STATIC_DRAW)
	log.Println(gl.GetError())

	/*gl.BindBuffer(gl.ARRAY_BUFFER, vbo1)
	gl.BufferData(gl.ARRAY_BUFFER, len(triangle2)*4, gl.Ptr(triangle2), gl.STATIC_DRAW)
	log.Println(gl.GetError())*/

	gl.BindVertexArray(vao)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(0))
	//gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(12))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(12))
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)
	gl.ClearColor(0, 0, 0, 1)

	//var timeValue float64

	for !window.ShouldClose() {
		//timeValue = glfw.GetTime()
		//color := math.Sin(timeValue)/2.0 + 0.5
		/*oof := gl.GetUniformLocation(program, gl.Str("theColor"+"\x00"))
		if oof == -1 {

			log.Println("Can't find uniform location")
			defer gl.DeleteProgram(program)
			os.Exit(1)
		}*/

		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.BindVertexArray(vao)
		gl.UseProgram(program)
		//gl.Uniform4f(oof, 0.0, float32(color), 0.0, 0.0)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		//gl.DrawArrays(gl.TRIANGLES, 1, 3)
		gl.BindVertexArray(0)
		//log.Println(gl.GetError())
		glfw.PollEvents()
		window.SwapBuffers()
	}

}

func getShaderSource(path string) (string, error) {

	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(buffer), nil

}
