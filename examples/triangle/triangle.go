package main

import (
	_ "embed"
	"unsafe"

	"github.com/idaunis/mtg/metal"
)

var (
	device       *metal.MTLDevice
	commandQueue *metal.MTLCommandQueue
	metalLayer   *metal.CAMetalLayer
	pipeline     *metal.MTLRenderPipelineState
	vertexBuffer *metal.MTLBuffer

	//go:embed triangle.metal
	metalSource string
)

func makeBuffers() *metal.MTLBuffer {
	vertices := [][2]metal.Vector_float4{
		{{0.0, 0.5, 0, 1}, {1, 0, 0, 1}},
		{{-0.5, -0.5, 0, 1}, {0, 1, 0, 1}},
		{{0.5, -0.5, 0, 1}, {0, 0, 1, 1}},
	}
	count := len(vertices)
	size := unsafe.Sizeof(vertices[0]) * uintptr(count)
	data := unsafe.Pointer(&vertices[0][0])

	return device.NewBufferWithBytes(data, size, count, metal.MTLResourceCPUCacheModeDefaultCache)
}

func initDelegate(view *metal.MTKView) {
	device = view.Device()
	metalLayer = view.Layer()
	commandQueue = device.NewCommandQueue()
	vertexBuffer = makeBuffers()

	library := device.NewLibraryWithSource(metalSource)
	vertexFunc := library.NewFunctionWithName("vertex_main")
	fragmentFunc := library.NewFunctionWithName("fragment_main")

	pipelineDescriptor := metal.NewMTLRenderPipelineDescriptor()
	pipelineDescriptor.SetVertexFunction(vertexFunc)
	pipelineDescriptor.SetFragmentFunction(fragmentFunc)
	pipelineDescriptor.ColorAttachment(0).SetPixelFormat(metal.MTLPixelFormatBGRA8Unorm)

	pipeline = device.NewRenderPipelineStateWithDescriptor(pipelineDescriptor)
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
	commandEncoder.SetVertexBuffer(vertexBuffer, 0, 0)
	commandEncoder.DrawPrimitives(metal.MTLPrimitiveTypeTriangle, 0, 3)

	commandEncoder.EndEncoding()
	commandBuffer.PresentDrawable(drawable)
	commandBuffer.Commit()
}

func main() {
	metal.CreateApp()
	metal.RenderDelegate(metal.CreateWindow(), initDelegate, drawDelegate)
	metal.RunApp()
}
