# mtg - Metal API support for Go

mtg is a Go package that adds support of the [Metal Framework](https://developer.apple.com/documentation/metal?language=objc) API for your Go applications.

mtg is a light weight cgo wrapper around the objective-c metal api, that allows to quickly build applications with metal.

mtg aims to be low-level, fast, and performant, and methods mimic as much as possible the metal API naming conventions used in swift and objective-c.

## Example Usage

```go
package main

import (
	"unsafe"

	"github.com/idaunis/mtg/metal"
)

var (
	device       *metal.MTLDevice
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

	return device.NewBufferWithBytes(data, size, count, metal.MTLResourceCPUCacheModeDefaultCache)
}

func initDelegate(view *metal.MTKView) {
	device = view.Device()
	metalLayer = view.Layer()
	commandQueue = device.NewCommandQueue()
	vertexBuffer = makeBuffers()

	metalSource := `
		using namespace metal;

		struct Vertex {
			float4 position [[position]];
			float4 color;
		};

		vertex Vertex vertex_main(const device Vertex *vertices [[buffer(0)]], uint vid [[vertex_id]]) {
			return vertices[vid];
		}

		fragment float4 fragment_main(Vertex inVertex [[stage_in]]) {
			return inVertex.color;
		}
	`
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

```

See the [examples](https://github.com/idaunis/mtg/examples) for more detailed examples (loading .obj files and textures).

![teapot example](https://github.com/idaunis/mtg/blob/main/examples/teapot/preview.png?raw=true)


## Debugging

To enable metal API validation for debugging set the following env var:

```
export METAL_DEVICE_WRAPPER_TYPE=1
```

## Features

* Command setup: queue/buffer/encoder including complete handler
* Rendering: vertex/fragment function support, render pipelines to draw primitives, loading textures, minmap generation
* Vector / Matrix / SIMD type support and basic matrix operations
* Support loading 3D objects with normals and UV mapping (Wavefront .obj file)
