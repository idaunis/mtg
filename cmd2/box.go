package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"unsafe"

	"github.com/idaunis/mtg/metal"
)

var (
	device            *metal.MTLDevice
	library           *metal.MTLLibrary
	commandQueue      *metal.MTLCommandQueue
	metalLayer        *metal.CAMetalLayer
	pipeline          *metal.MTLRenderPipelineState
	vertexBuffer      *metal.MTLBuffer
	indexBuffer       *metal.MTLBuffer
	uniformBuffer     *metal.MTLBuffer
	depthStencilState *metal.MTLDepthStencilState
)

type uniforms struct {
	modelViewMatrix  metal.Matrix_float4x4
	projectionMatrix metal.Matrix_float4x4
}

func makeBuffers() (*metal.MTLBuffer, *metal.MTLBuffer, *metal.MTLBuffer) {
	vertices := [][2]metal.Vector_float4{
		{{-1, 1, 1, 1}, {0, 1, 1, 1}},
		{{-1, -1, 1, 1}, {0, 0, 1, 1}},
		{{1, -1, 1, 1}, {1, 0, 1, 1}},
		{{1, 1, 1, 1}, {1, 1, 1, 1}},
		{{-1, 1, -1, 1}, {0, 1, 0, 1}},
		{{-1, -1, -1, 1}, {0, 0, 0, 1}},
		{{1, -1, -1, 1}, {1, 0, 0, 1}},
		{{1, 1, -1, 1}, {1, 1, 0, 1}},
	}

	indices := []metal.Uint16{
		3, 2, 6, 6, 7, 3,
		4, 5, 1, 1, 0, 4,
		4, 0, 3, 3, 7, 4,
		1, 5, 6, 6, 2, 1,
		0, 1, 2, 2, 3, 0,
		7, 6, 5, 5, 4, 7,
	}

	uniforms := uniforms{}

	size := unsafe.Sizeof(uniforms)
	b := (*[1 << 30]byte)(unsafe.Pointer(&uniforms))[0:size]
	alignedUniformsSize := align256(size)
	fmt.Println("***", b, size, alignedUniformsSize)

	fmt.Println("***", unsafe.Sizeof(indices[0])*uintptr(len(indices)))

	return device.NewBufferWithVectors2(vertices, metal.MTLResourceCPUCacheModeDefaultCache),
		// device.NewBufferWithInts(indices, metal.MTLResourceCPUCacheModeDefaultCache),
		device.NewBufferWithBytes(unsafe.Pointer(&indices[0]), unsafe.Sizeof(indices[0]), len(indices), metal.MTLResourceCPUCacheModeDefaultCache),
		device.NewBufferWithBytes(unsafe.Pointer(&uniforms), unsafe.Sizeof(uniforms), 1, metal.MTLResourceCPUCacheModeDefaultCache)
}

func align256(size uintptr) uintptr {
	return uintptr((int(size) + 0xFF) & -0x100)
}

func initDelegate(view *metal.MTKView) {
	device = view.Device()
	metalLayer = view.Layer()

	view.SetDepthStencilPixelFormat(metal.MTLPixelFormatDepth32Float)

	vertexBuffer, indexBuffer, uniformBuffer = makeBuffers()

	pipelineDescriptor := metal.NewMTLRenderPipelineDescriptor()
	pipelineDescriptor.SetVertexFunction(library.NewFunctionWithName("vertex_project"))
	pipelineDescriptor.SetFragmentFunction(library.NewFunctionWithName("fragment_flatcolor"))
	pipelineDescriptor.ColorAttachment(0).SetPixelFormat(metal.MTLPixelFormatBGRA8Unorm)
	pipelineDescriptor.SetDepthAttachmentPixelFormat(metal.MTLPixelFormatDepth32Float)

	depthStencilDescriptor := metal.NewMTLDepthStencilDescriptor()
	depthStencilDescriptor.SetDepthCompareFunction(metal.MTLCompareFunctionLess)
	depthStencilDescriptor.SetDepthWriteEnabled(true)
	depthStencilState = device.NewDepthStencilStateWithDescriptor(depthStencilDescriptor)

	pipeline = device.NewRenderPipelineStateWithDescriptor(pipelineDescriptor)

	commandQueue = device.NewCommandQueue()

	fmt.Println("InitWithMetalKitView", view, device, metalLayer, commandQueue, library, "pipeline:", pipeline)
}

var (
	time      float32
	rotationX float32
	rotationY float32
)

func updateUniforms(layer *metal.CAMetalLayer) {
	duration := float32(1) / 60
	time += duration
	rotationX += duration * (math.Pi / 2)
	rotationY += duration * (math.Pi / 3)
	scaleFactor := float32(math.Sin(5*float64(time)))*0.25 + 1
	xAxis := metal.Vector_float3{1, 0, 0}
	yAxis := metal.Vector_float3{0, 1, 0}
	xRot := metal.Matrix_float4x4_rotation(xAxis, rotationX)
	yRot := metal.Matrix_float4x4_rotation(yAxis, rotationY)

	scale := metal.Matrix_float4x4_uniform_scale(scaleFactor)

	modelMatrix := metal.Matrix_multiply(metal.Matrix_multiply(xRot, yRot), scale)
	cameraTranslation := metal.Vector_float3{0, 0, -5}
	viewMatrix := metal.Matrix_float4x4_translation(cameraTranslation)

	aspect := layer.DrawableSize().Width / layer.DrawableSize().Height
	fov := float32((2 * math.Pi) / 5)
	near := float32(1)
	far := float32(100)

	projectionMatrix := metal.Matrix_float4x4_perspective(aspect, fov, near, far)

	uniforms := uniforms{
		projectionMatrix: projectionMatrix,
		modelViewMatrix:  metal.Matrix_multiply(viewMatrix, modelMatrix),
	}

	uniformBuffer.ContentsCopy(unsafe.Pointer(&uniforms), unsafe.Sizeof(uniforms), 0)
}

func drawDelegate(view *metal.MTKView) {
	drawable := metalLayer.NextDrawable()
	texture := drawable.Texture()

	updateUniforms(view.Layer())

	renderPassDescriptor := view.CurrentRenderPassDescriptor()

	renderPassDescriptor.ColorAttachment(&metal.ColorAttachment{
		Texture:     texture,
		LoadAction:  metal.MTLLoadActionClear,
		StoreAction: metal.MTLStoreActionStore,
		ClearColor:  metal.MTLClearColor{Red: 0.1, Green: .1, Blue: .2, Alpha: 1},
	})

	commandBuffer := commandQueue.CommandBuffer()
	commandEncoder := commandBuffer.RenderCommandEncoderWithDescriptor(renderPassDescriptor)

	commandEncoder.SetRenderPipelineState(pipeline)

	commandEncoder.SetDepthStencilState(depthStencilState)
	commandEncoder.SetFrontFacingWinding(metal.MTLWindingCounterClockwise)
	commandEncoder.SetCullMode(metal.MTLCullModeBack)

	commandEncoder.SetVertexBuffer(vertexBuffer, 0, 0)
	commandEncoder.SetVertexBuffer(uniformBuffer, 0, 1)
	commandEncoder.DrawIndexedPrimitives(metal.MTLPrimitiveTypeTriangle, indexBuffer.Count, metal.MTLIndexTypeUInt16, indexBuffer, 0)

	commandEncoder.EndEncoding()

	commandBuffer.PresentDrawable(drawable)
	commandBuffer.Commit()
}

func libraryFromFile(device *metal.MTLDevice, name string) (*metal.MTLLibrary, error) {
	if device == nil {
		return nil, errors.New("device not initialized")
	}
	source, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	// TODO: have two methods one to load lib as source and another to load a compiled metallib
	// as a reminder to compile a metallib we can do:
	// > xcrun -sdk macosx metal -c shaders.metal -o shaders.air
	// > xcrun -sdk macosx metallib shaders.air -o shaders.metallib

	return device.NewLibraryWithSource(string(source)), nil
}

func main() {
	var err error

	metal.CreateApp()
	w := metal.CreateWindow()

	library, err = libraryFromFile(w.Device(), "box.metal")
	if err != nil {
		fmt.Println(err)
		return
	}

	metal.RenderDelegate(w, initDelegate, drawDelegate)
	metal.RunApp()
}
