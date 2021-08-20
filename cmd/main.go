package main

import (
	"fmt"
	"io/ioutil"
	"unsafe"

	"test/metal"
)

var (
	device       *metal.MTLDevice
	commandQueue *metal.MTLCommandQueue
	metalLayer   *metal.CAMetalLayer
)

func makeBuffers() {
	vertices := []metal.MBEVertex{
		{Position: metal.Vector_float4{0.0, 0.5, 0, 1}, Color: metal.Vector_float4{1, 0, 0, 1}},
		{Position: metal.Vector_float4{-0.5, -0.5, 0, 1}, Color: metal.Vector_float4{0, 1, 0, 1}},
		{Position: metal.Vector_float4{0.5, -0.5, 0, 1}, Color: metal.Vector_float4{0, 0, 1, 1}},
	}
	device.NewBufferWithBytes(vertices, int(unsafe.Sizeof(vertices[0]))*len(vertices), metal.MTLResourceCPUCacheModeDefaultCache)
}

func initDelegate(view *metal.MTKView) {
	device = view.Device()
	metalLayer = view.Layer()
	commandQueue = device.NewCommandQueue()

	makeBuffers()

	source, err := ioutil.ReadFile("shaders.metal")
	if err != nil {
		fmt.Print(err)
	}

	library := device.NewLibraryWithSource(string(source))
	vertexFunc := library.NewFunctionWithName("vertex_main")
	fragmentFunc := library.NewFunctionWithName("fragment_main")

	fmt.Println("InitWithMetalKitView", view, device, metalLayer, commandQueue, library, vertexFunc, fragmentFunc)
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

	// [commandEncoder setRenderPipelineState:self.pipeline]
	// [commandEncoder setVertexBuffer:self.vertexBuffer offset:0 atIndex:0]
	// [commandEncoder drawPrimitives:MTLPrimitiveTypeTriangle vertexStart:0 vertexCount:3];

	commandEncoder.EndEncoding()
	commandBuffer.PresentDrawable(drawable)
	commandBuffer.Commit()

	fmt.Println(metalLayer, drawable, renderPassDescriptor, texture, commandEncoder)
}

func main() {
	metal.CreateApp()
	w := metal.CreateWindow()
	metal.RenderDelegate(w, initDelegate, drawDelegate)
	metal.RunApp()
}
