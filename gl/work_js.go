// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build js,wasm

package gl

import "syscall/js"

type CanvasContext interface {
	Canvas() js.Value
}

func NewContext() (Context, Worker) {
	canvas := js.Global().Get("document").Call("createElement", "canvas")

	canvas.Set("width", js.Global().Get("innerWidth"))
	canvas.Set("height", js.Global().Get("innerHeight"))

	gl2 := canvas.Call("getContext", "webgl2")
	if gl2.Truthy() {
		return context3{context: context{canvas: canvas, ctx: gl2}}, noopWorker{}
	}

	gl := canvas.Call("getContext", "webgl")
	if gl.Truthy() {
		return context{canvas: canvas, ctx: gl}, noopWorker{}
	}

	panic("unable to get WebGL context")
}

func Version() string {
	canvas := js.Global().Get("document").Call("createElement", "canvas")

	canvas.Set("width", 1)
	canvas.Set("height", 1)

	if canvas.Call("getContext", "webgl2").Truthy() {
		return "GL_ES_3_0"
	}

	return "GL_ES_2_0"
}

func jsBytes(data []byte) js.Value {
	arr := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(arr, data)

	return arr
}

func jsInts(data []int32) js.Value {
	arr := js.Global().Get("Int32Array").New(len(data))

	for i, n := range data {
		arr.SetIndex(i, n)
	}

	return arr
}

func jsFloats(data []float32) js.Value {
	arr := js.Global().Get("Float32Array").New(len(data))

	for i, f := range data {
		arr.SetIndex(i, f)
	}

	return arr
}

func goInt(value js.Value) int {
	if value.Type() == js.TypeBoolean {
		if value.Bool() {
			return 1
		}

		return 0
	}

	return value.Int()
}

func getParameter(arr js.Value, index int) js.Value {
	if arr.Type() == js.TypeObject {
		return arr.Index(index)
	}

	if index == 0 {
		return arr
	}

	return js.Undefined()
}

type noopWorker struct{}

func (noopWorker) DoWork() {
	// work is done synchronously
}

func (noopWorker) WorkAvailable() <-chan struct{} {
	return nil
}

type context struct {
	canvas js.Value
	ctx    js.Value
}

var _ CanvasContext = context{}
var _ CanvasContext = context3{}

func (c context) Canvas() js.Value {
	return c.canvas
}

type context3 struct {
	context
}

func (c context) ActiveTexture(texture Enum) {
	c.ctx.Call("activeTexture", texture)
}

func (c context) AttachShader(p Program, s Shader) {
	c.ctx.Call("attachShader", p.Value, s.Value)
}

func (c context) BindAttribLocation(p Program, a Attrib, name string) {
	c.ctx.Call("bindAttribLocation", p.Value, a.Value, name)
}

func (c context) BindBuffer(target Enum, b Buffer) {
	c.ctx.Call("bindBuffer", target, b.Value)
}

func (c context) BindFramebuffer(target Enum, fb Framebuffer) {
	c.ctx.Call("bindFramebuffer", target, fb.Value)
}

func (c context) BindRenderbuffer(target Enum, rb Renderbuffer) {
	c.ctx.Call("bindRenderbuffer", target, rb.Value)
}

func (c context) BindTexture(target Enum, t Texture) {
	c.ctx.Call("bindTexture", target, t.Value)
}

func (c context) BindVertexArray(rb VertexArray) {
	c.ctx.Call("bindVertexArray", rb.Value)
}

func (c context) BlendColor(red, green, blue, alpha float32) {
	c.ctx.Call("blendColor", red, green, blue, alpha)
}

func (c context) BlendEquation(mode Enum) {
	c.ctx.Call("blendEquation", mode)
}

func (c context) BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	c.ctx.Call("blendEquationSeparate", modeRGB, modeAlpha)
}

func (c context) BlendFunc(sfactor, dfactor Enum) {
	c.ctx.Call("blendFunc", sfactor, dfactor)
}

func (c context) BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	c.ctx.Call("blendFuncSeparate", sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha)
}

func (c context) BufferData(target Enum, src []byte, usage Enum) {
	c.ctx.Call("bufferData", target, jsBytes(src), usage)
}

func (c context) BufferInit(target Enum, size int, usage Enum) {
	c.ctx.Call("bufferInit", target, size, usage)
}

func (c context) BufferSubData(target Enum, offset int, data []byte) {
	c.ctx.Call("target", offset, jsBytes(data))
}

func (c context) CheckFramebufferStatus(target Enum) Enum {
	return Enum(c.ctx.Call("checkFramebufferStatus", target).Int())
}

func (c context) Clear(mask Enum) {
	c.ctx.Call("clear", mask)
}

func (c context) ClearColor(red, green, blue, alpha float32) {
	c.ctx.Call("clearColor", red, green, blue, alpha)
}

func (c context) ClearDepthf(d float32) {
	c.ctx.Call("clearDepth", d)
}

func (c context) ClearStencil(s int) {
	c.ctx.Call("clearStencil", s)
}

func (c context) ColorMask(red, green, blue, alpha bool) {
	c.ctx.Call("colorMask", red, green, blue, alpha)
}

func (c context) CompileShader(s Shader) {
	c.ctx.Call("compileShader", s.Value)
}

func (c context) CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte) {
	c.ctx.Call("compressedTexImage2D", target, level, internalformat, width, height, border, jsBytes(data))
}

func (c context) CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte) {
	c.ctx.Call("compressedTexSubImage2D", target, level, xoffset, yoffset, width, height, format, jsBytes(data))
}

func (c context) CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	c.ctx.Call("copyTexImage2D", target, level, internalformat, x, y, width, height, border)
}

func (c context) CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	c.ctx.Call("copyTexSubImage2D", target, level, xoffset, yoffset, x, y, width, height)
}

func (c context) CreateBuffer() Buffer {
	return Buffer{Value: c.ctx.Call("createBuffer")}
}

func (c context) CreateFramebuffer() Framebuffer {
	return Framebuffer{Value: c.ctx.Call("createFramebuffer")}
}

func (c context) CreateProgram() Program {
	return Program{Init: true, Value: c.ctx.Call("createProgram")}
}

func (c context) CreateRenderbuffer() Renderbuffer {
	return Renderbuffer{Value: c.ctx.Call("createRenderbuffer")}
}

func (c context) CreateShader(ty Enum) Shader {
	return Shader{Value: c.ctx.Call("createShader", ty)}
}

func (c context) CreateTexture() Texture {
	return Texture{Value: c.ctx.Call("createTexture")}
}

func (c context) CreateVertexArray() VertexArray {
	return VertexArray{Value: c.ctx.Call("createVertexArray")}
}

func (c context) CullFace(mode Enum) {
	c.ctx.Call("cullFace", mode)
}

func (c context) DeleteBuffer(v Buffer) {
	c.ctx.Call("deleteBuffer", v.Value)
}

func (c context) DeleteFramebuffer(v Framebuffer) {
	c.ctx.Call("deleteFramebuffer", v.Value)
}

func (c context) DeleteProgram(p Program) {
	c.ctx.Call("deleteProgram", p.Value)
}

func (c context) DeleteRenderbuffer(v Renderbuffer) {
	c.ctx.Call("deleteRenderbuffer", v.Value)
}

func (c context) DeleteShader(s Shader) {
	c.ctx.Call("deleteShader", s.Value)
}

func (c context) DeleteTexture(v Texture) {
	c.ctx.Call("deleteTexture", v.Value)
}

func (c context) DeleteVertexArray(v VertexArray) {
	c.ctx.Call("deleteVertexArray", v.Value)
}

func (c context) DepthFunc(fn Enum) {
	c.ctx.Call("depthFunc", fn)
}

func (c context) DepthMask(flag bool) {
	c.ctx.Call("depthMask", flag)
}

func (c context) DepthRangef(n, f float32) {
	c.ctx.Call("depthRangef", n, f)
}

func (c context) DetachShader(p Program, s Shader) {
	c.ctx.Call("detachShader", p.Value, s.Value)
}

func (c context) Disable(cap Enum) {
	c.ctx.Call("disable", cap)
}

func (c context) DisableVertexAttribArray(a Attrib) {
	c.ctx.Call("disableVertexAttribArray", a.Value)
}

func (c context) DrawArrays(mode Enum, first, count int) {
	c.ctx.Call("drawArrays", mode, first, count)
}

func (c context) DrawElements(mode Enum, count int, ty Enum, offset int) {
	c.ctx.Call("darwElements", mode, count, ty, offset)
}

func (c context) Enable(cap Enum) {
	c.ctx.Call("enable", cap)
}

func (c context) EnableVertexAttribArray(a Attrib) {
	c.ctx.Call("enableVertexAttribArray", a.Value)
}

func (c context) Finish() {
	c.ctx.Call("finish")
}

func (c context) Flush() {
	c.ctx.Call("flush")
}

func (c context) FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer) {
	c.ctx.Call("framebufferRenderbuffer", target, attachment, rbTarget, rb.Value)
}

func (c context) FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	c.ctx.Call("framebufferTexture2D", target, attachment, texTarget, t.Value, level)
}

func (c context) FrontFace(mode Enum) {
	c.ctx.Call("frontFace", mode)
}

func (c context) GenerateMipmap(target Enum) {
	c.ctx.Call("generateMipmap", target)
}

func (c context) GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	active := c.ctx.Call("getActiveAttrib", p.Value, index)

	return active.Get("name").String(), active.Get("size").Int(), Enum(active.Get("type").Int())
}

func (c context) GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	active := c.ctx.Call("getActiveUniform", p.Value, index)

	return active.Get("name").String(), active.Get("size").Int(), Enum(active.Get("type").Int())
}

func (c context) GetAttachedShaders(p Program) []Shader {
	arr := c.ctx.Call("getAttachedShaders", p.Value)

	shaders := make([]Shader, arr.Length())
	for i := range shaders {
		shaders[i].Value = arr.Index(i)
	}

	return shaders
}

func (c context) GetAttribLocation(p Program, name string) Attrib {
	return Attrib{Value: uint(c.ctx.Call("getAttribLocation", p.Value, name).Int())}
}

func (c context) GetBooleanv(dst []bool, pname Enum) {
	arr := c.ctx.Call("getParameter", pname)

	for i := range dst {
		dst[i] = getParameter(arr, i).Bool()
	}
}

func (c context) GetFloatv(dst []float32, pname Enum) {
	arr := c.ctx.Call("getParameter", pname)

	for i := range dst {
		dst[i] = float32(getParameter(arr, i).Float())
	}
}

func (c context) GetIntegerv(dst []int32, pname Enum) {
	arr := c.ctx.Call("getParameter", pname)

	for i := range dst {
		dst[i] = int32(getParameter(arr, i).Int())
	}
}

func (c context) GetInteger(pname Enum) int {
	return goInt(c.ctx.Call("getParameter", pname))
}

func (c context) GetBufferParameteri(target, value Enum) int {
	return goInt(c.ctx.Call("getBufferParameter", target, value))
}

func (c context) GetError() Enum {
	return Enum(c.ctx.Call("getError").Int())
}

func (c context) GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	return goInt(c.ctx.Call("getFramebufferAttachmentParameter", target, attachment, pname))
}

func (c context) GetProgrami(p Program, pname Enum) int {
	return goInt(c.ctx.Call("getProgramParameter", p.Value, pname))
}

func (c context) GetProgramInfoLog(p Program) string {
	return c.ctx.Call("getProgramInfoLog", p.Value).String()
}

func (c context) GetRenderbufferParameteri(target, pname Enum) int {
	return goInt(c.ctx.Call("getRenderbufferParameter", target, pname))
}

func (c context) GetShaderi(s Shader, pname Enum) int {
	return goInt(c.ctx.Call("getShaderParameter", s.Value, pname))
}

func (c context) GetShaderInfoLog(s Shader) string {
	return c.ctx.Call("getShaderInfoLog", s.Value).String()
}

func (c context) GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeLow, rangeHigh, precision int) {
	format := c.ctx.Call("getShaderPrecisionFormat", shadertype, precisiontype)

	return format.Get("rangeMin").Int(), format.Get("rangeMax").Int(), format.Get("precision").Int()
}

func (c context) GetShaderSource(s Shader) string {
	return c.ctx.Call("getShaderSource", s.Value).String()
}

func (c context) GetString(pname Enum) string {
	return c.ctx.Call("getParameter", pname).String()
}

func (c context) GetTexParameterfv(dst []float32, target, pname Enum) {
	arr := c.ctx.Call("getTexParameter", target, pname)

	for i := range dst {
		dst[i] = float32(getParameter(arr, i).Float())
	}
}

func (c context) GetTexParameteriv(dst []int32, target, pname Enum) {
	arr := c.ctx.Call("getTexParameter", target, pname)

	for i := range dst {
		dst[i] = int32(getParameter(arr, i).Int())
	}
}

func (c context) GetUniformfv(dst []float32, src Uniform, p Program) {
	arr := c.ctx.Call("getUniform", p.Value, src.Value)

	for i := range dst {
		dst[i] = float32(getParameter(arr, i).Float())
	}
}

func (c context) GetUniformiv(dst []int32, src Uniform, p Program) {
	arr := c.ctx.Call("getUniform", p.Value, src.Value)

	for i := range dst {
		dst[i] = int32(getParameter(arr, i).Int())
	}
}

func (c context) GetUniformLocation(p Program, name string) Uniform {
	return Uniform{Value: c.ctx.Call("getUniformLocation", p.Value, name)}
}

func (c context) GetVertexAttribf(src Attrib, pname Enum) float32 {
	return float32(c.ctx.Call("getVertexAttrib", src.Value, pname).Float())
}

func (c context) GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	arr := c.ctx.Call("getVertexAttrib", src.Value, pname)

	for i := range dst {
		dst[i] = float32(getParameter(arr, i).Float())
	}
}

func (c context) GetVertexAttribi(src Attrib, pname Enum) int32 {
	return int32(c.ctx.Call("getVertexAttrib", src.Value, pname).Int())
}

func (c context) GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	arr := c.ctx.Call("getVertexAttrib", src.Value, pname)

	for i := range dst {
		dst[i] = int32(getParameter(arr, i).Int())
	}
}

func (c context) Hint(target, mode Enum) {
	c.ctx.Call("hint", target, mode)
}

func (c context) IsBuffer(b Buffer) bool {
	return c.ctx.Call("isBuffer", b.Value).Bool()
}

func (c context) IsEnabled(cap Enum) bool {
	return c.ctx.Call("isEnabled", cap).Bool()
}

func (c context) IsFramebuffer(fb Framebuffer) bool {
	return c.ctx.Call("isFramebuffer", fb.Value).Bool()
}

func (c context) IsProgram(p Program) bool {
	return c.ctx.Call("isProgram", p.Value).Bool()
}

func (c context) IsRenderbuffer(rb Renderbuffer) bool {
	return c.ctx.Call("isRenderbuffer", rb.Value).Bool()
}

func (c context) IsShader(s Shader) bool {
	return c.ctx.Call("isShader", s.Value).Bool()
}

func (c context) IsTexture(t Texture) bool {
	return c.ctx.Call("isTexture", t.Value).Bool()
}

func (c context) LineWidth(width float32) {
	c.ctx.Call("lineWidth", width)
}

func (c context) LinkProgram(p Program) {
	c.ctx.Call("linkProgram", p.Value)
}

func (c context) PixelStorei(pname Enum, param int32) {
	c.ctx.Call("pixelStorei", pname, param)
}

func (c context) PolygonOffset(factor, units float32) {
	c.ctx.Call("polygonOffset", factor, units)
}

func (c context) ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
	jsBuf := js.Global().Get("Uint8Array").New(len(dst))
	jsTypedBuf := jsBuf

	switch ty {
	case UNSIGNED_SHORT_5_6_5, UNSIGNED_SHORT_4_4_4_4, UNSIGNED_SHORT_5_5_5_1:
		jsTypedBuf = js.Global().Get("Uint16Array").New(jsBuf.Get("buffer"), 0)
	case FLOAT:
		jsTypedBuf = js.Global().Get("Float32Array").New(jsBuf.Get("buffer"), 0)
	}

	c.ctx.Call("readPixels", x, y, width, height, format, ty, jsTypedBuf)

	js.CopyBytesToGo(dst, jsBuf)
}

func (c context) ReleaseShaderCompiler() {
	// not in WebGL
}

func (c context) RenderbufferStorage(target, internalFormat Enum, width, height int) {
	c.ctx.Call("renderbufferStorage", target, internalFormat, width, height)
}

func (c context) SampleCoverage(value float32, invert bool) {
	c.ctx.Call("sampleCoverage", value, invert)
}

func (c context) Scissor(x, y, width, height int32) {
	c.ctx.Call("scissor", x, y, width, height)
}

func (c context) ShaderSource(s Shader, src string) {
	c.ctx.Call("shaderSource", s.Value, src)
}

func (c context) StencilFunc(fn Enum, ref int, mask uint32) {
	c.ctx.Call("stencilFunc", fn, ref, mask)
}

func (c context) StencilFuncSeparate(face, fn Enum, ref int, mask uint32) {
	c.ctx.Call("stencilFuncSeparate", face, fn, ref, mask)
}

func (c context) StencilMask(mask uint32) {
	c.ctx.Call("stencilMask", mask)
}

func (c context) StencilMaskSeparate(face Enum, mask uint32) {
	c.ctx.Call("stencilMaskSeparate", face, mask)
}

func (c context) StencilOp(fail, zfail, zpass Enum) {
	c.ctx.Call("stencilOp", fail, zfail, zpass)
}

func (c context) StencilOpSeparate(face, sfail, dpfail, dppass Enum) {
	c.ctx.Call("stencilOpSeparate", face, sfail, dpfail, dppass)
}

func (c context) TexImage2D(target Enum, level int, internalFormat int, width, height int, format Enum, ty Enum, data []byte) {
	c.ctx.Call("texImage2D", target, level, internalFormat, width, height, format, ty, jsBytes(data))
}

func (c context) TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	c.ctx.Call("texSubImage2D", target, level, x, y, width, height, format, ty, jsBytes(data))
}

func (c context) TexParameterf(target, pname Enum, param float32) {
	c.ctx.Call("texParameterf", target, pname, param)
}

func (c context) TexParameterfv(target, pname Enum, params []float32) {
	c.ctx.Call("texParameterf", target, pname, params[0])
}

func (c context) TexParameteri(target, pname Enum, param int) {
	c.ctx.Call("texParameteri", target, pname, param)
}

func (c context) TexParameteriv(target, pname Enum, params []int32) {
	c.ctx.Call("texParameteri", target, pname, params[0])
}

func (c context) Uniform1f(dst Uniform, v float32) {
	c.ctx.Call("uniform1f", dst.Value, v)
}

func (c context) Uniform1fv(dst Uniform, src []float32) {
	c.ctx.Call("uniform1fv", dst.Value, jsFloats(src))
}

func (c context) Uniform1i(dst Uniform, v int) {
	c.ctx.Call("uniform1i", dst.Value, v)
}

func (c context) Uniform1iv(dst Uniform, src []int32) {
	c.ctx.Call("uniform1iv", dst.Value, jsInts(src))
}

func (c context) Uniform2f(dst Uniform, v0, v1 float32) {
	c.ctx.Call("uniform2f", dst.Value, v0, v1)
}

func (c context) Uniform2fv(dst Uniform, src []float32) {
	c.ctx.Call("uniform2fv", dst.Value, jsFloats(src))
}

func (c context) Uniform2i(dst Uniform, v0, v1 int) {
	c.ctx.Call("uniform2i", dst.Value, v0, v1)
}

func (c context) Uniform2iv(dst Uniform, src []int32) {
	c.ctx.Call("uniform2iv", dst.Value, jsInts(src))
}

func (c context) Uniform3f(dst Uniform, v0, v1, v2 float32) {
	c.ctx.Call("uniform3f", dst.Value, v0, v1, v2)
}

func (c context) Uniform3fv(dst Uniform, src []float32) {
	c.ctx.Call("uniform3fv", dst, jsFloats(src))
}

func (c context) Uniform3i(dst Uniform, v0, v1, v2 int32) {
	c.ctx.Call("uniform3i", dst.Value, v0, v1, v2)
}

func (c context) Uniform3iv(dst Uniform, src []int32) {
	c.ctx.Call("uniform3iv", dst.Value, jsInts(src))
}

func (c context) Uniform4f(dst Uniform, v0, v1, v2, v3 float32) {
	c.ctx.Call("uniform4f", dst.Value, v0, v1, v2, v3)
}

func (c context) Uniform4fv(dst Uniform, src []float32) {
	c.ctx.Call("uniform4fv", dst.Value, jsFloats(src))
}

func (c context) Uniform4i(dst Uniform, v0, v1, v2, v3 int32) {
	c.ctx.Call("uniform4i", dst.Value, v0, v1, v2, v3)
}

func (c context) Uniform4iv(dst Uniform, src []int32) {
	c.ctx.Call("uniform4iv", dst.Value, jsInts(src))
}

func (c context) UniformMatrix2fv(dst Uniform, src []float32) {
	c.ctx.Call("uniformMatrix2fv", dst.Value, jsFloats(src))
}

func (c context) UniformMatrix3fv(dst Uniform, src []float32) {
	c.ctx.Call("uniformMatrix3fv", dst.Value, jsFloats(src))
}

func (c context) UniformMatrix4fv(dst Uniform, src []float32) {
	c.ctx.Call("uniformMatrix4fv", dst.Value, jsFloats(src))
}

func (c context) UseProgram(p Program) {
	c.ctx.Call("useProgram", p.Value)
}

func (c context) ValidateProgram(p Program) {
	c.ctx.Call("validateProgram", p.Value)
}

func (c context) VertexAttrib1f(dst Attrib, x float32) {
	c.ctx.Call("vertexAttrib1f", dst.Value, x)
}

func (c context) VertexAttrib1fv(dst Attrib, src []float32) {
	c.ctx.Call("vertexAttrib1fv", dst.Value, jsFloats(src))
}

func (c context) VertexAttrib2f(dst Attrib, x, y float32) {
	c.ctx.Call("vertexAttrib2f", dst.Value, x, y)
}

func (c context) VertexAttrib2fv(dst Attrib, src []float32) {
	c.ctx.Call("vertexAttrib2fv", dst.Value, jsFloats(src))
}

func (c context) VertexAttrib3f(dst Attrib, x, y, z float32) {
	c.ctx.Call("vertexAttrib3f", dst.Value, x, y, z)
}

func (c context) VertexAttrib3fv(dst Attrib, src []float32) {
	c.ctx.Call("vertexAttrib3fv", dst.Value, jsFloats(src))
}

func (c context) VertexAttrib4f(dst Attrib, x, y, z, w float32) {
	c.ctx.Call("vertexAttrib4f", dst.Value, x, y, z, w)
}

func (c context) VertexAttrib4fv(dst Attrib, src []float32) {
	c.ctx.Call("vertexAttrib4fv", dst.Value, jsFloats(src))
}

func (c context) VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	c.ctx.Call("vertexAttribPointer", dst.Value, size, ty, normalized, stride, offset)
}

func (c context) Viewport(x, y, width, height int) {
	c.ctx.Call("viewport", x, y, width, height)
}
