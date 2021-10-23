package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sync"
	"unsafe"

	"github.com/idaunis/mtg/metal"
	"github.com/idaunis/mtg/obj"
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
	samplerState      *metal.MTLSamplerState
	diffuseTexture    *metal.MTLTexture
	wg                sync.WaitGroup
)

type uniforms struct {
	modelViewProjectionMatrix metal.Matrix_float4x4 // 4 x 16 = 64
	modelViewMatrix           metal.Matrix_float4x4 // 4 x 16 = 64
	normalMatrix              metal.Matrix_float3x3 // 3 x 12 = 36 + 12
}

// these are packed floats so are not 16-bytes aligned
type vertex struct {
	position  metal.Vector_float4 // 4 float32 (4 x 4bytes = 16 bytes)
	normal    metal.Vector_float4
	texCoords metal.Vector_float2
}

func makeTexture() *metal.MTLTexture {
	// width := 100
	// height := 100
	// imageData := []byte{0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255}
	minmapped := true
	imageName := "spot_texture.png"

	reader, err := os.Open(imageName)
	if err != nil {
		log.Fatal(err)
	}
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	size := m.Bounds().Size()

	width := size.X
	height := size.Y
	bytesPerRow := width * 4
	imageData := make([]byte, bytesPerRow*height)

	i := 0
	for y := size.Y - 1; y >= 0; y-- {
		for x := 0; x < size.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			imageData[i], imageData[i+1], imageData[i+2], imageData[i+3] = byte(r), byte(g), byte(b), byte(a)
			i += 4
		}
	}

	textureDescriptor := metal.Texture2DDescriptorWithPixelFormat(metal.MTLPixelFormatRGBA8Unorm, width, height, minmapped)
	textureDescriptor.SetUsage(metal.MTLTextureUsageShaderRead)

	texture := device.NewTextureWithDescriptor(textureDescriptor)
	texture.SetLabel(imageName)

	region := metal.MTLRegionMake2D(0, 0, width, height)
	texture.ReplaceRegion(region, 0, imageData, bytesPerRow)

	commandBuffer := commandQueue.CommandBuffer()
	blitEncoder := commandBuffer.BlitCommandEncoder()

	blitEncoder.GenerateMipmapsForTexture(texture)
	blitEncoder.EndEncoding()

	commandBuffer.Commit()
	commandBuffer.WaitUntilCompleted()

	return texture
}

func makeBuffers() (*metal.MTLBuffer, *metal.MTLBuffer, *metal.MTLBuffer) {
	uniforms := uniforms{}
	vertices := []vertex{}
	model, _ := obj.Parse("spot.obj")
	group := model.GetGroup(1)

	// group.GenerateNormals()
	group.EachVertex(func(p obj.PackedVertex) {
		vertices = append(vertices, vertex{p.Position, p.Normal, p.TexCoord})
	})
	indices := group.Indices

	return device.NewBufferWithBytes(unsafe.Pointer(&vertices[0]), unsafe.Sizeof(vertices[0]), len(vertices), metal.MTLResourceCPUCacheModeDefaultCache),
		device.NewBufferWithBytes(unsafe.Pointer(&indices[0]), unsafe.Sizeof(indices[0]), len(indices), metal.MTLResourceCPUCacheModeDefaultCache),
		device.NewBufferWithBytes(unsafe.Pointer(&uniforms), unsafe.Sizeof(uniforms), 1, metal.MTLResourceCPUCacheModeDefaultCache)
}

func initDelegate(view *metal.MTKView) {
	device = view.Device()
	metalLayer = view.Layer()

	view.SetDepthStencilPixelFormat(metal.MTLPixelFormatDepth32Float)

	vertexBuffer, indexBuffer, uniformBuffer = makeBuffers()

	pipelineDescriptor := metal.NewMTLRenderPipelineDescriptor()
	pipelineDescriptor.SetVertexFunction(library.NewFunctionWithName("vertex_project"))
	pipelineDescriptor.SetFragmentFunction(library.NewFunctionWithName("fragment_light"))
	pipelineDescriptor.ColorAttachment(0).SetPixelFormat(metal.MTLPixelFormatRGBA8Unorm)
	pipelineDescriptor.SetDepthAttachmentPixelFormat(metal.MTLPixelFormatDepth32Float)

	depthStencilDescriptor := metal.NewMTLDepthStencilDescriptor()
	depthStencilDescriptor.SetDepthCompareFunction(metal.MTLCompareFunctionLess)
	depthStencilDescriptor.SetDepthWriteEnabled(true)
	depthStencilState = device.NewDepthStencilStateWithDescriptor(depthStencilDescriptor)

	pipeline = device.NewRenderPipelineStateWithDescriptor(pipelineDescriptor)

	commandQueue = device.NewCommandQueue()

	diffuseTexture = makeTexture()

	samplerDescriptor := metal.NewMTLSamplerDescriptor()
	samplerDescriptor.SetSAddressMode(metal.MTLSamplerAddressModeClampToEdge)
	samplerDescriptor.SetTAddressMode(metal.MTLSamplerAddressModeClampToEdge)
	samplerDescriptor.SetMinFilter(metal.MTLSamplerMinMagFilterNearest)
	samplerDescriptor.SetMagFilter(metal.MTLSamplerMinMagFilterLinear)
	samplerDescriptor.SetMipFilter(metal.MTLSamplerMipFilterLinear)
	samplerState = device.NewSamplerStateWithDescriptor(samplerDescriptor)
}

var (
	time      float32
	rotationX float32
	rotationY float32
	rotationZ float32
)

func updateUniforms(layer *metal.CAMetalLayer, yPos float64) {
	scaleFactor := float32(1.0)
	xAxis := metal.Vector_float3{1, 0, 0}
	yAxis := metal.Vector_float3{0, 1, 0}
	zAxis := metal.Vector_float3{0, 0, 1}
	xRot := metal.Matrix_float4x4_rotation(xAxis, rotationX)
	yRot := metal.Matrix_float4x4_rotation(yAxis, rotationY)
	zRot := metal.Matrix_float4x4_rotation(zAxis, rotationZ)

	scale := metal.Matrix_float4x4_uniform_scale(scaleFactor)

	modelMatrix := metal.Matrix_multiply(metal.Matrix_multiply(xRot, metal.Matrix_multiply(zRot, yRot)), scale)

	cameraTranslation := metal.Vector3(0, yPos, -1.0)
	viewMatrix := metal.Matrix_float4x4_translation(cameraTranslation)

	aspect := layer.DrawableSize().Width / layer.DrawableSize().Height
	fov := float32((2 * math.Pi) / 5)
	near := float32(.1)
	far := float32(100)

	projectionMatrix := metal.Matrix_float4x4_perspective(aspect, fov, near, far)

	uniforms := uniforms{}
	uniforms.modelViewMatrix = metal.Matrix_multiply(viewMatrix, modelMatrix)
	uniforms.modelViewProjectionMatrix = metal.Matrix_multiply(projectionMatrix, uniforms.modelViewMatrix)
	uniforms.normalMatrix = metal.Matrix_float4x4_extract_linear(uniforms.modelViewMatrix)

	uniformBuffer.ContentsCopy(unsafe.Pointer(&uniforms), unsafe.Sizeof(uniforms), 0)
}

func drawDelegate(view *metal.MTKView) {
	drawable := metalLayer.NextDrawable()
	texture := drawable.Texture()

	duration := float32(1) / 60
	rotationY += duration * (math.Pi / 2)
	time += duration
	updateUniforms(view.Layer(), 0)

	wg.Add(1)
	commandBuffer := commandQueue.CommandBuffer()
	commandBuffer.AddCompletedHandler(func(cb *metal.MTLCommandBuffer) {
		wg.Done()
	})

	renderPassDescriptor := view.CurrentRenderPassDescriptor()
	renderPassDescriptor.ColorAttachment(&metal.ColorAttachment{
		Texture:     texture,
		LoadAction:  metal.MTLLoadActionClear,
		StoreAction: metal.MTLStoreActionStore,
		ClearColor:  metal.MTLClearColor{Red: 1, Green: 1, Blue: 1, Alpha: 1},
	})
	commandEncoder := commandBuffer.RenderCommandEncoderWithDescriptor(renderPassDescriptor)

	commandEncoder.SetRenderPipelineState(pipeline)

	commandEncoder.SetDepthStencilState(depthStencilState)
	commandEncoder.SetFrontFacingWinding(metal.MTLWindingCounterClockwise)
	commandEncoder.SetCullMode(metal.MTLCullModeBack)
	commandEncoder.SetFragmentTexture(diffuseTexture, 0)
	commandEncoder.SetFragmentSamplerState(samplerState, 0)

	commandEncoder.SetVertexBuffer(vertexBuffer, 0, 0)
	commandEncoder.SetVertexBuffer(uniformBuffer, 0, 1)
	commandEncoder.DrawIndexedPrimitives(metal.MTLPrimitiveTypeTriangle, indexBuffer.Count, metal.MTLIndexTypeUInt16, indexBuffer, 0)

	commandEncoder.EndEncoding()

	commandBuffer.PresentDrawable(drawable)
	commandBuffer.Commit()
	wg.Wait()
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
	metal.CreateApp()
	w := metal.CreateWindow()

	var err error
	library, err = libraryFromFile(w.Device(), "shaders.metal")
	if err != nil {
		fmt.Println(err)
		return
	}

	metal.RenderDelegate(w, initDelegate, drawDelegate)
	metal.RunApp()
}
