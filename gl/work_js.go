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
		ctx := &context3{context: context{canvas: canvas, ctx: gl2}}

		ctx.bind()

		return ctx, noopWorker{}
	}

	gl := canvas.Call("getContext", "webgl")
	if gl.Truthy() {
		ctx := &context{canvas: canvas, ctx: gl}

		ctx.bind()

		return ctx, noopWorker{}
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

var uint8Array = js.Global().Get("Uint8Array")

func jsBytes(data []byte) js.Value {
	if data == nil {
		return js.Null()
	}

	arr := uint8Array.New(len(data))
	js.CopyBytesToJS(arr, data)

	return arr
}

var float32Array = js.Global().Get("Float32Array")

func jsFloats(data []float32) js.Value {
	if data == nil {
		return js.Null()
	}

	arr := float32Array.New(len(data))

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

	activeTexture            js.Value
	bindBuffer               js.Value
	bindFramebuffer          js.Value
	bindTexture              js.Value
	bufferData               js.Value
	clear                    js.Value
	clearColor               js.Value
	disableVertexAttribArray js.Value
	drawElements             js.Value
	enableVertexAttribArray  js.Value
	uniform1f                js.Value
	uniform1i                js.Value
	uniform2f                js.Value
	uniform2i                js.Value
	uniform3f                js.Value
	uniform3i                js.Value
	uniform4f                js.Value
	uniform4i                js.Value
	uniformMatrix2fv         js.Value
	uniformMatrix3fv         js.Value
	uniformMatrix4fv         js.Value
	useProgram               js.Value
	vertexAttribPointer      js.Value
	viewport                 js.Value
}

var _ CanvasContext = (*context)(nil)
var _ CanvasContext = (*context3)(nil)

// TODO: var _ Context3 = (*context3)(nil)

func (c *context) Canvas() js.Value {
	return c.canvas
}

type context3 struct {
	context
}

func (e Enum) JSValue() js.Value {
	return js.ValueOf(int(e))
}

func (c *context) bind() {
	c.activeTexture = c.ctx.Get("activeTexture").Call("bind", c.ctx)
	c.bindBuffer = c.ctx.Get("bindBuffer").Call("bind", c.ctx)
	c.bindFramebuffer = c.ctx.Get("bindFramebuffer").Call("bind", c.ctx)
	c.bindTexture = c.ctx.Get("bindTexture").Call("bind", c.ctx)
	c.bufferData = c.ctx.Get("bufferData").Call("bind", c.ctx)
	c.clear = c.ctx.Get("clear").Call("bind", c.ctx)
	c.clearColor = c.ctx.Get("clearColor").Call("bind", c.ctx)
	c.disableVertexAttribArray = c.ctx.Get("disableVertexAttribArray").Call("bind", c.ctx)
	c.drawElements = c.ctx.Get("drawElements").Call("bind", c.ctx)
	c.enableVertexAttribArray = c.ctx.Get("enableVertexAttribArray").Call("bind", c.ctx)
	c.uniform1f = c.ctx.Get("uniform1f").Call("bind", c.ctx)
	c.uniform1i = c.ctx.Get("uniform1i").Call("bind", c.ctx)
	c.uniform2f = c.ctx.Get("uniform2f").Call("bind", c.ctx)
	c.uniform2i = c.ctx.Get("uniform2i").Call("bind", c.ctx)
	c.uniform3f = c.ctx.Get("uniform3f").Call("bind", c.ctx)
	c.uniform3i = c.ctx.Get("uniform3i").Call("bind", c.ctx)
	c.uniform4f = c.ctx.Get("uniform4f").Call("bind", c.ctx)
	c.uniform4i = c.ctx.Get("uniform4i").Call("bind", c.ctx)
	c.uniformMatrix2fv = c.ctx.Get("uniformMatrix2fv").Call("bind", c.ctx)
	c.uniformMatrix3fv = c.ctx.Get("uniformMatrix3fv").Call("bind", c.ctx)
	c.uniformMatrix4fv = c.ctx.Get("uniformMatrix4fv").Call("bind", c.ctx)
	c.useProgram = c.ctx.Get("useProgram").Call("bind", c.ctx)
	c.vertexAttribPointer = c.ctx.Get("vertexAttribPointer").Call("bind", c.ctx)
	c.viewport = c.ctx.Get("viewport").Call("bind", c.ctx)
}

func (c *context) ActiveTexture(texture Enum) {
	c.activeTexture.Invoke(texture)
}

func (c *context) AttachShader(p Program, s Shader) {
	c.ctx.Call("attachShader", p.Value, s.Value)
}

func (c *context) BindAttribLocation(p Program, a Attrib, name string) {
	c.ctx.Call("bindAttribLocation", p.Value, a.Value, name)
}

func (c *context) BindBuffer(target Enum, b Buffer) {
	c.bindBuffer.Invoke(target, b.Value)
}

func (c *context) BindFramebuffer(target Enum, fb Framebuffer) {
	c.bindFramebuffer.Invoke(target, fb.Value)
}

func (c *context) BindRenderbuffer(target Enum, rb Renderbuffer) {
	c.ctx.Call("bindRenderbuffer", target, rb.Value)
}

func (c *context) BindTexture(target Enum, t Texture) {
	c.bindTexture.Invoke(target, t.Value)
}

func (c *context) BindVertexArray(rb VertexArray) {
	c.ctx.Call("bindVertexArray", rb.Value)
}

func (c *context) BlendColor(red, green, blue, alpha float32) {
	c.ctx.Call("blendColor", red, green, blue, alpha)
}

func (c *context) BlendEquation(mode Enum) {
	c.ctx.Call("blendEquation", mode)
}

func (c *context) BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	c.ctx.Call("blendEquationSeparate", modeRGB, modeAlpha)
}

func (c *context) BlendFunc(sfactor, dfactor Enum) {
	c.ctx.Call("blendFunc", sfactor, dfactor)
}

func (c *context) BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	c.ctx.Call("blendFuncSeparate", sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha)
}

func (c *context) BufferData(target Enum, src []byte, usage Enum) {
	c.bufferData.Invoke(target, jsBytes(src), usage)
}

func (c *context) BufferInit(target Enum, size int, usage Enum) {
	c.ctx.Call("bufferInit", target, size, usage)
}

func (c *context) BufferSubData(target Enum, offset int, data []byte) {
	c.ctx.Call("bufferSubData", target, offset, jsBytes(data))
}

func (c *context) CheckFramebufferStatus(target Enum) Enum {
	return Enum(c.ctx.Call("checkFramebufferStatus", target).Int())
}

func (c *context) Clear(mask Enum) {
	c.clear.Invoke(mask)
}

func (c *context) ClearColor(red, green, blue, alpha float32) {
	c.clearColor.Invoke(red, green, blue, alpha)
}

func (c *context) ClearDepthf(d float32) {
	c.ctx.Call("clearDepth", d)
}

func (c *context) ClearStencil(s int) {
	c.ctx.Call("clearStencil", s)
}

func (c *context) ColorMask(red, green, blue, alpha bool) {
	c.ctx.Call("colorMask", red, green, blue, alpha)
}

func (c *context) CompileShader(s Shader) {
	c.ctx.Call("compileShader", s.Value)
}

func (c *context) CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte) {
	c.ctx.Call("compressedTexImage2D", target, level, internalformat, width, height, border, jsBytes(data))
}

func (c *context) CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte) {
	c.ctx.Call("compressedTexSubImage2D", target, level, xoffset, yoffset, width, height, format, jsBytes(data))
}

func (c *context) CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	c.ctx.Call("copyTexImage2D", target, level, internalformat, x, y, width, height, border)
}

func (c *context) CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	c.ctx.Call("copyTexSubImage2D", target, level, xoffset, yoffset, x, y, width, height)
}

func (c *context) CreateBuffer() Buffer {
	return Buffer{Value: c.ctx.Call("createBuffer")}
}

func (c *context) CreateFramebuffer() Framebuffer {
	return Framebuffer{Value: c.ctx.Call("createFramebuffer")}
}

func (c *context) CreateProgram() Program {
	return Program{Init: true, Value: c.ctx.Call("createProgram")}
}

func (c *context) CreateRenderbuffer() Renderbuffer {
	return Renderbuffer{Value: c.ctx.Call("createRenderbuffer")}
}

func (c *context) CreateShader(ty Enum) Shader {
	return Shader{Value: c.ctx.Call("createShader", ty)}
}

func (c *context) CreateTexture() Texture {
	return Texture{Value: c.ctx.Call("createTexture")}
}

func (c *context) CreateVertexArray() VertexArray {
	return VertexArray{Value: c.ctx.Call("createVertexArray")}
}

func (c *context) CullFace(mode Enum) {
	c.ctx.Call("cullFace", mode)
}

func (c *context) DeleteBuffer(v Buffer) {
	c.ctx.Call("deleteBuffer", v.Value)
}

func (c *context) DeleteFramebuffer(v Framebuffer) {
	c.ctx.Call("deleteFramebuffer", v.Value)
}

func (c *context) DeleteProgram(p Program) {
	c.ctx.Call("deleteProgram", p.Value)
}

func (c *context) DeleteRenderbuffer(v Renderbuffer) {
	c.ctx.Call("deleteRenderbuffer", v.Value)
}

func (c *context) DeleteShader(s Shader) {
	c.ctx.Call("deleteShader", s.Value)
}

func (c *context) DeleteTexture(v Texture) {
	c.ctx.Call("deleteTexture", v.Value)
}

func (c *context) DeleteVertexArray(v VertexArray) {
	c.ctx.Call("deleteVertexArray", v.Value)
}

func (c *context) DepthFunc(fn Enum) {
	c.ctx.Call("depthFunc", fn)
}

func (c *context) DepthMask(flag bool) {
	c.ctx.Call("depthMask", flag)
}

func (c *context) DepthRangef(n, f float32) {
	c.ctx.Call("depthRangef", n, f)
}

func (c *context) DetachShader(p Program, s Shader) {
	c.ctx.Call("detachShader", p.Value, s.Value)
}

func (c *context) Disable(cap Enum) {
	c.ctx.Call("disable", cap)
}

func (c *context) DisableVertexAttribArray(a Attrib) {
	c.disableVertexAttribArray.Invoke(a.Value)
}

func (c *context) DrawArrays(mode Enum, first, count int) {
	c.ctx.Call("drawArrays", mode, first, count)
}

func (c *context) DrawElements(mode Enum, count int, ty Enum, offset int) {
	c.drawElements.Invoke(mode, count, ty, offset)
}

func (c *context) Enable(cap Enum) {
	c.ctx.Call("enable", cap)
}

func (c *context) EnableVertexAttribArray(a Attrib) {
	c.enableVertexAttribArray.Invoke(a.Value)
}

func (c *context) Finish() {
	c.ctx.Call("finish")
}

func (c *context) Flush() {
	c.ctx.Call("flush")
}

func (c *context) FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer) {
	c.ctx.Call("framebufferRenderbuffer", target, attachment, rbTarget, rb.Value)
}

func (c *context) FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	c.ctx.Call("framebufferTexture2D", target, attachment, texTarget, t.Value, level)
}

func (c *context) FrontFace(mode Enum) {
	c.ctx.Call("frontFace", mode)
}

func (c *context) GenerateMipmap(target Enum) {
	c.ctx.Call("generateMipmap", target)
}

func (c *context) GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	active := c.ctx.Call("getActiveAttrib", p.Value, index)

	return active.Get("name").String(), active.Get("size").Int(), Enum(active.Get("type").Int())
}

func (c *context) GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	active := c.ctx.Call("getActiveUniform", p.Value, index)

	return active.Get("name").String(), active.Get("size").Int(), Enum(active.Get("type").Int())
}

func (c *context) GetAttachedShaders(p Program) []Shader {
	arr := c.ctx.Call("getAttachedShaders", p.Value)

	shaders := make([]Shader, arr.Length())
	for i := range shaders {
		shaders[i].Value = arr.Index(i)
	}

	return shaders
}

func (c *context) GetAttribLocation(p Program, name string) Attrib {
	return Attrib{Value: uint(c.ctx.Call("getAttribLocation", p.Value, name).Int())}
}

func (c *context) GetBooleanv(dst []bool, pname Enum) {
	arr := c.ctx.Call("getParameter", pname)

	for i := range dst {
		dst[i] = getParameter(arr, i).Bool()
	}
}

func (c *context) GetFloatv(dst []float32, pname Enum) {
	arr := c.ctx.Call("getParameter", pname)

	for i := range dst {
		dst[i] = float32(getParameter(arr, i).Float())
	}
}

func (c *context) GetIntegerv(dst []int32, pname Enum) {
	arr := c.ctx.Call("getParameter", pname)

	for i := range dst {
		dst[i] = int32(getParameter(arr, i).Int())
	}
}

func (c *context) GetInteger(pname Enum) int {
	return goInt(c.ctx.Call("getParameter", pname))
}

func (c *context) GetBufferParameteri(target, value Enum) int {
	return goInt(c.ctx.Call("getBufferParameter", target, value))
}

func (c *context) GetError() Enum {
	return Enum(c.ctx.Call("getError").Int())
}

func (c *context) GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	return goInt(c.ctx.Call("getFramebufferAttachmentParameter", target, attachment, pname))
}

func (c *context) GetProgrami(p Program, pname Enum) int {
	return goInt(c.ctx.Call("getProgramParameter", p.Value, pname))
}

func (c *context) GetProgramInfoLog(p Program) string {
	return c.ctx.Call("getProgramInfoLog", p.Value).String()
}

func (c *context) GetRenderbufferParameteri(target, pname Enum) int {
	return goInt(c.ctx.Call("getRenderbufferParameter", target, pname))
}

func (c *context) GetShaderi(s Shader, pname Enum) int {
	return goInt(c.ctx.Call("getShaderParameter", s.Value, pname))
}

func (c *context) GetShaderInfoLog(s Shader) string {
	return c.ctx.Call("getShaderInfoLog", s.Value).String()
}

func (c *context) GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeLow, rangeHigh, precision int) {
	format := c.ctx.Call("getShaderPrecisionFormat", shadertype, precisiontype)

	return format.Get("rangeMin").Int(), format.Get("rangeMax").Int(), format.Get("precision").Int()
}

func (c *context) GetShaderSource(s Shader) string {
	return c.ctx.Call("getShaderSource", s.Value).String()
}

func (c *context) GetString(pname Enum) string {
	return c.ctx.Call("getParameter", pname).String()
}

func (c *context) GetTexParameterfv(dst []float32, target, pname Enum) {
	arr := c.ctx.Call("getTexParameter", target, pname)

	for i := range dst {
		dst[i] = float32(getParameter(arr, i).Float())
	}
}

func (c *context) GetTexParameteriv(dst []int32, target, pname Enum) {
	arr := c.ctx.Call("getTexParameter", target, pname)

	for i := range dst {
		dst[i] = int32(getParameter(arr, i).Int())
	}
}

func (c *context) GetUniformfv(dst []float32, src Uniform, p Program) {
	arr := c.ctx.Call("getUniform", p.Value, src.Value)

	for i := range dst {
		dst[i] = float32(getParameter(arr, i).Float())
	}
}

func (c *context) GetUniformiv(dst []int32, src Uniform, p Program) {
	arr := c.ctx.Call("getUniform", p.Value, src.Value)

	for i := range dst {
		dst[i] = int32(getParameter(arr, i).Int())
	}
}

func (c *context) GetUniformLocation(p Program, name string) Uniform {
	return Uniform{Value: c.ctx.Call("getUniformLocation", p.Value, name)}
}

func (c *context) GetVertexAttribf(src Attrib, pname Enum) float32 {
	return float32(c.ctx.Call("getVertexAttrib", src.Value, pname).Float())
}

func (c *context) GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	arr := c.ctx.Call("getVertexAttrib", src.Value, pname)

	for i := range dst {
		dst[i] = float32(getParameter(arr, i).Float())
	}
}

func (c *context) GetVertexAttribi(src Attrib, pname Enum) int32 {
	return int32(c.ctx.Call("getVertexAttrib", src.Value, pname).Int())
}

func (c *context) GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	arr := c.ctx.Call("getVertexAttrib", src.Value, pname)

	for i := range dst {
		dst[i] = int32(getParameter(arr, i).Int())
	}
}

func (c *context) Hint(target, mode Enum) {
	c.ctx.Call("hint", target, mode)
}

func (c *context) IsBuffer(b Buffer) bool {
	return c.ctx.Call("isBuffer", b.Value).Bool()
}

func (c *context) IsEnabled(cap Enum) bool {
	return c.ctx.Call("isEnabled", cap).Bool()
}

func (c *context) IsFramebuffer(fb Framebuffer) bool {
	return c.ctx.Call("isFramebuffer", fb.Value).Bool()
}

func (c *context) IsProgram(p Program) bool {
	return c.ctx.Call("isProgram", p.Value).Bool()
}

func (c *context) IsRenderbuffer(rb Renderbuffer) bool {
	return c.ctx.Call("isRenderbuffer", rb.Value).Bool()
}

func (c *context) IsShader(s Shader) bool {
	return c.ctx.Call("isShader", s.Value).Bool()
}

func (c *context) IsTexture(t Texture) bool {
	return c.ctx.Call("isTexture", t.Value).Bool()
}

func (c *context) LineWidth(width float32) {
	c.ctx.Call("lineWidth", width)
}

func (c *context) LinkProgram(p Program) {
	c.ctx.Call("linkProgram", p.Value)
}

func (c *context) PixelStorei(pname Enum, param int32) {
	c.ctx.Call("pixelStorei", pname, param)
}

func (c *context) PolygonOffset(factor, units float32) {
	c.ctx.Call("polygonOffset", factor, units)
}

func (c *context) ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
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

func (c *context) ReleaseShaderCompiler() {
	// not in WebGL
}

func (c *context) RenderbufferStorage(target, internalFormat Enum, width, height int) {
	c.ctx.Call("renderbufferStorage", target, internalFormat, width, height)
}

func (c *context) SampleCoverage(value float32, invert bool) {
	c.ctx.Call("sampleCoverage", value, invert)
}

func (c *context) Scissor(x, y, width, height int32) {
	c.ctx.Call("scissor", x, y, width, height)
}

func (c *context) ShaderSource(s Shader, src string) {
	c.ctx.Call("shaderSource", s.Value, src)
}

func (c *context) StencilFunc(fn Enum, ref int, mask uint32) {
	c.ctx.Call("stencilFunc", fn, ref, mask)
}

func (c *context) StencilFuncSeparate(face, fn Enum, ref int, mask uint32) {
	c.ctx.Call("stencilFuncSeparate", face, fn, ref, mask)
}

func (c *context) StencilMask(mask uint32) {
	c.ctx.Call("stencilMask", mask)
}

func (c *context) StencilMaskSeparate(face Enum, mask uint32) {
	c.ctx.Call("stencilMaskSeparate", face, mask)
}

func (c *context) StencilOp(fail, zfail, zpass Enum) {
	c.ctx.Call("stencilOp", fail, zfail, zpass)
}

func (c *context) StencilOpSeparate(face, sfail, dpfail, dppass Enum) {
	c.ctx.Call("stencilOpSeparate", face, sfail, dpfail, dppass)
}

func (c *context) TexImage2D(target Enum, level int, internalFormat int, width, height int, format Enum, ty Enum, data []byte) {
	c.ctx.Call("texImage2D", target, level, internalFormat, width, height, 0, format, ty, jsBytes(data))
}

func (c *context) TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	c.ctx.Call("texSubImage2D", target, level, x, y, width, height, format, ty, jsBytes(data))
}

func (c *context) TexParameterf(target, pname Enum, param float32) {
	c.ctx.Call("texParameterf", target, pname, param)
}

func (c *context) TexParameterfv(target, pname Enum, params []float32) {
	c.ctx.Call("texParameterf", target, pname, params[0])
}

func (c *context) TexParameteri(target, pname Enum, param int) {
	c.ctx.Call("texParameteri", target, pname, param)
}

func (c *context) TexParameteriv(target, pname Enum, params []int32) {
	c.ctx.Call("texParameteri", target, pname, params[0])
}

func (c *context) Uniform1f(dst Uniform, v float32) {
	c.uniform1f.Invoke(dst.Value, v)
}

func (c *context) Uniform1fv(dst Uniform, src []float32) {
	c.uniform1f.Invoke(dst.Value, src[0])
}

func (c *context) Uniform1i(dst Uniform, v int) {
	c.uniform1i.Invoke(dst.Value, v)
}

func (c *context) Uniform1iv(dst Uniform, src []int32) {
	c.uniform1i.Invoke(dst.Value, src[0])
}

func (c *context) Uniform2f(dst Uniform, v0, v1 float32) {
	c.uniform2f.Invoke(dst.Value, v0, v1)
}

func (c *context) Uniform2fv(dst Uniform, src []float32) {
	c.uniform2f.Invoke(dst.Value, src[0], src[1])
}

func (c *context) Uniform2i(dst Uniform, v0, v1 int) {
	c.uniform2i.Invoke(dst.Value, v0, v1)
}

func (c *context) Uniform2iv(dst Uniform, src []int32) {
	c.uniform2i.Invoke(dst.Value, src[0], src[1])
}

func (c *context) Uniform3f(dst Uniform, v0, v1, v2 float32) {
	c.uniform3f.Invoke(dst.Value, v0, v1, v2)
}

func (c *context) Uniform3fv(dst Uniform, src []float32) {
	c.uniform3f.Invoke(dst.Value, src[0], src[1], src[2])
}

func (c *context) Uniform3i(dst Uniform, v0, v1, v2 int32) {
	c.uniform3i.Invoke(dst.Value, v0, v1, v2)
}

func (c *context) Uniform3iv(dst Uniform, src []int32) {
	c.uniform3i.Invoke(dst.Value, src[0], src[1], src[2])
}

func (c *context) Uniform4f(dst Uniform, v0, v1, v2, v3 float32) {
	c.uniform4f.Invoke(dst.Value, v0, v1, v2, v3)
}

func (c *context) Uniform4fv(dst Uniform, src []float32) {
	c.uniform4f.Invoke(dst.Value, src[0], src[1], src[2], src[3])
}

func (c *context) Uniform4i(dst Uniform, v0, v1, v2, v3 int32) {
	c.uniform4i.Invoke(dst.Value, v0, v1, v2, v3)
}

func (c *context) Uniform4iv(dst Uniform, src []int32) {
	c.uniform4i.Invoke(dst.Value, src[0], src[1], src[2], src[3])
}

func (c *context) UniformMatrix2fv(dst Uniform, src []float32) {
	c.uniformMatrix2fv.Invoke(dst.Value, false, jsFloats(src))
}

func (c *context) UniformMatrix3fv(dst Uniform, src []float32) {
	c.uniformMatrix3fv.Invoke(dst.Value, false, jsFloats(src))
}

func (c *context) UniformMatrix4fv(dst Uniform, src []float32) {
	c.uniformMatrix4fv.Invoke(dst.Value, false, jsFloats(src))
}

func (c *context) UseProgram(p Program) {
	c.useProgram.Invoke(p.Value)
}

func (c *context) ValidateProgram(p Program) {
	c.ctx.Call("validateProgram", p.Value)
}

func (c *context) VertexAttrib1f(dst Attrib, x float32) {
	c.ctx.Call("vertexAttrib1f", dst.Value, x)
}

func (c *context) VertexAttrib1fv(dst Attrib, src []float32) {
	c.ctx.Call("vertexAttrib1fv", dst.Value, jsFloats(src))
}

func (c *context) VertexAttrib2f(dst Attrib, x, y float32) {
	c.ctx.Call("vertexAttrib2f", dst.Value, x, y)
}

func (c *context) VertexAttrib2fv(dst Attrib, src []float32) {
	c.ctx.Call("vertexAttrib2fv", dst.Value, jsFloats(src))
}

func (c *context) VertexAttrib3f(dst Attrib, x, y, z float32) {
	c.ctx.Call("vertexAttrib3f", dst.Value, x, y, z)
}

func (c *context) VertexAttrib3fv(dst Attrib, src []float32) {
	c.ctx.Call("vertexAttrib3fv", dst.Value, jsFloats(src))
}

func (c *context) VertexAttrib4f(dst Attrib, x, y, z, w float32) {
	c.ctx.Call("vertexAttrib4f", dst.Value, x, y, z, w)
}

func (c *context) VertexAttrib4fv(dst Attrib, src []float32) {
	c.ctx.Call("vertexAttrib4fv", dst.Value, jsFloats(src))
}

func (c *context) VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	c.vertexAttribPointer.Invoke(dst.Value, size, ty, normalized, stride, offset)
}

func (c *context) Viewport(x, y, width, height int) {
	c.viewport.Invoke(x, y, width, height)
}
