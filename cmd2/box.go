package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"unsafe"

	"test/metal"
)

var (
	device        *metal.MTLDevice
	library       *metal.MTLLibrary
	commandQueue  *metal.MTLCommandQueue
	metalLayer    *metal.CAMetalLayer
	pipeline      *metal.MTLRenderPipelineState
	vertexBuffer  *metal.MTLBuffer
	indexBuffer   *metal.MTLBuffer
	uniformBuffer *metal.MTLBuffer
)

type uniforms struct {
	modelViewProjectionMatrix metal.Matrix_float4x4
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
		3, 2, 6, 6, 7, 3, 4, 5, 1, 1, 0, 4, 4, 0, 3, 3, 7, 4, 1, 5, 6, 6, 2, 1, 0, 1, 2, 2, 3, 0, 7, 6, 5, 5, 4, 7,
	}

	uniforms := uniforms{
		modelViewProjectionMatrix: metal.Matrix_float4x4{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		},
	}

	size := unsafe.Sizeof(uniforms)
	b := (*[1 << 30]byte)(unsafe.Pointer(&uniforms))[0:size]
	fmt.Println("***", b)

	return device.NewBufferWithVectors2(vertices, metal.MTLResourceCPUCacheModeDefaultCache),
		device.NewBufferWithInts(indices, metal.MTLResourceCPUCacheModeDefaultCache),
		device.NewBufferWithBytes(unsafe.Pointer(&uniforms), unsafe.Sizeof(uniforms), 1, metal.MTLResourceCPUCacheModeDefaultCache)
}

func initDelegate(view *metal.MTKView) {
	device = view.Device()
	metalLayer = view.Layer()
	commandQueue = device.NewCommandQueue()

	vertexBuffer, indexBuffer, uniformBuffer = makeBuffers()

	vertexFunc := library.NewFunctionWithName("vertex_project")
	fragmentFunc := library.NewFunctionWithName("fragment_flatcolor")

	pipelineDescriptor := metal.NewMTLRenderPipelineDescriptor()
	pipelineDescriptor.SetVertexFunction(vertexFunc)
	pipelineDescriptor.SetFragmentFunction(fragmentFunc)
	pipelineDescriptor.ColorAttachment(0).SetPixelFormat(metal.MTLPixelFormatBGRA8Unorm)
	// pipelineDescriptor.SetDepthAttachmentPixelFormat(metal.MTLPixelFormatDepth32Float)

	pipeline = device.NewRenderPipelineStateWithDescriptor(pipelineDescriptor)

	fmt.Println("InitWithMetalKitView", view, device, metalLayer, commandQueue, library, vertexFunc, fragmentFunc, "pipeline:", pipeline)
}

func drawDelegate(view *metal.MTKView) {
	drawable := metalLayer.NextDrawable()
	texture := drawable.Texture()

	renderPassDescriptor := view.CurrentRenderPassDescriptor()

	renderPassDescriptor.ColorAttachment(&metal.ColorAttachment{
		Texture:     texture,
		LoadAction:  metal.MTLLoadActionClear,
		StoreAction: metal.MTLStoreActionStore,
		ClearColor:  metal.MTLClearColor{Red: 1, Green: 1, Blue: 0, Alpha: 1},
	})

	commandBuffer := commandQueue.CommandBuffer()
	commandEncoder := commandBuffer.RenderCommandEncoderWithDescriptor(renderPassDescriptor)

	commandEncoder.SetRenderPipelineState(pipeline)

	//[commandEncoder setDepthStencilState:self.depthStencilState];
	//[commandEncoder setFrontFacingWinding:MTLWindingCounterClockwise];
	//[commandEncoder setCullMode:MTLCullModeBack];

	commandEncoder.SetVertexBuffer(vertexBuffer, 0, 0)
	commandEncoder.SetVertexBuffer(uniformBuffer, 0, 1)
	commandEncoder.DrawIndexedPrimitives(metal.MTLPrimitiveTypeTriangle, indexBuffer.Count, metal.MTLIndexTypeUInt16, indexBuffer, 0)

	commandEncoder.EndEncoding()
	commandBuffer.PresentDrawable(drawable)
	commandBuffer.Commit()

	fmt.Println(metalLayer, drawable, renderPassDescriptor, texture, commandEncoder)
}

func libraryFromFile(device *metal.MTLDevice, name string) (*metal.MTLLibrary, error) {
	if device == nil {
		return nil, errors.New("device not initialized")
	}
	source, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return device.NewLibraryWithSource(string(source)), nil
}

func main() {
	var err error

	metal.CreateApp()
	w := metal.CreateWindow()

	library, err = libraryFromFile(w.Device(), "shaders.metal")
	if err != nil {
		fmt.Println(err)
		return
	}

	metal.RenderDelegate(w, initDelegate, drawDelegate)
	metal.RunApp()
}
