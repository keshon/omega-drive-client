// From https://github.com/goxjs/gl
package nanovgo

import (
	"os"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// CreateBuffer creates a buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenBuffers.xhtml
func CreateBuffer() uint32 {
	var b uint32
	gl.GenBuffers(1, &b)
	return b
}

// CreateTexture creates a texture object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenTextures.xhtml
func CreateTexture() uint32 {
	var t uint32
	gl.GenTextures(1, &t)
	return t
}

// GetAttribLocation returns the location of an attribute variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttribLocation.xhtml
func GetAttribLocation(p uint32, name string) uint32 {
	location := gl.GetAttribLocation(p, gl.Str(name+"\x00"))
	err := gl.GetError()
	if err != gl.NO_ERROR {
		dumpLog("Error %08x after %s\n", int(err), "GetAttribLocation")
		os.Exit(0)
	}
	return uint32(location)
}

// GetProgrami returns a parameter value for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramiv.xhtml
func GetProgrami(p uint32, pname uint32) int {
	var result int32
	gl.GetProgramiv(p, pname, &result)
	return int(result)
}

// GetProgramInfoLog returns the information log for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramInfoLog.xhtml
func GetProgramInfoLog(p uint32) string {
	var logLength int32
	gl.GetProgramiv(p, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetProgramInfoLog(p, logLength, nil, &logBuffer[0])
	return gl.GoStr(&logBuffer[0])
}

// GetRenderbufferParameteri returns a parameter value for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderiv.xhtml
func GetShaderi(s uint32, pname uint32) int {
	var result int32
	gl.GetShaderiv(s, pname, &result)
	return int(result)
}

// GetShaderInfoLog returns the information log for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderInfoLog.xhtml
func GetShaderInfoLog(s uint32) string {
	var logLength int32
	gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetShaderInfoLog(s, logLength, nil, &logBuffer[0])
	return gl.GoStr(&logBuffer[0])
}

// ShaderSource sets the source code of s to the given source code.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glShaderSource.xhtml
func ShaderSource(s uint32, src string) {
	glsource, free := gl.Strs(src + "\x00")
	gl.ShaderSource(s, 1, glsource, nil)
	free()
}

// TexImage2D writes a 2D texture image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexImage2D.xhtml
func TexImage2D(target uint32, level int, width, height int, format uint32, ty uint32, data []byte) {
	p := unsafe.Pointer(nil)
	if len(data) > 0 {
		p = gl.Ptr(&data[0])
	}
	gl.TexImage2D(target, int32(level), int32(format), int32(width), int32(height), 0, format, ty, p)
}
