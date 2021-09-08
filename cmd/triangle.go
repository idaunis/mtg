package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"unsafe"

	"github.com/idaunis/mtg/metal"
)

var (
	device       *metal.MTLDevice
	library      *metal.MTLLibrary
	commandQueue *metal.MTLCommandQueue
	metalLayer   *metal.CAMetalLayer
	pipeline     *metal.MTLRenderPipelineState
	vertexBuffer *metal.MTLBuffer
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

	// b := (*[1 << 30]byte)(unsafe.Pointer(&vertices[0][0]))[0:size]

	// return device.NewBufferWithVectors2(vertices, metal.MTLResourceCPUCacheModeDefaultCache)
	return device.NewBufferWithBytes(data, size, count, metal.MTLResourceCPUCacheModeDefaultCache)
}

func initDelegate(view *metal.MTKView) {
	device = view.Device()
	metalLayer = view.Layer()
	commandQueue = device.NewCommandQueue()

	vertexBuffer = makeBuffers()

	vertexFunc := library.NewFunctionWithName("vertex_main")
	fragmentFunc := library.NewFunctionWithName("fragment_main")

	pipelineDescriptor := metal.NewMTLRenderPipelineDescriptor()
	pipelineDescriptor.SetVertexFunction(vertexFunc)
	pipelineDescriptor.SetFragmentFunction(fragmentFunc)
	pipelineDescriptor.ColorAttachment(0).SetPixelFormat(metal.MTLPixelFormatBGRA8Unorm)

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
	commandEncoder.SetVertexBuffer(vertexBuffer, 0, 0)
	commandEncoder.DrawPrimitives(metal.MTLPrimitiveTypeTriangle, 0, 3)

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

	library, err = libraryFromFile(w.Device(), "triangle.metal")
	if err != nil {
		fmt.Println(err)
		return
	}

	metal.RenderDelegate(w, initDelegate, drawDelegate)
	metal.RunApp()
}
