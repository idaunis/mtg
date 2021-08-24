package metal

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Metal -framework MetalKit
#include "renderer2.h"
#include <stdio.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type CAMetalLayer struct {
	ptr unsafe.Pointer
}

func (s *CAMetalLayer) NextDrawable() *CAMetalDrawable {
	ptr := C.CAMetalLayer_nextDrawable(s.ptr)
	return &CAMetalDrawable{ptr}
}

func (s *CAMetalLayer) DrawableSize() CGSize {
	size := C.CAMetalLayer_drawableSize(s.ptr)
	return CGSize{float32(size.width), float32(size.height)}
}

type CGSize struct {
	Width  float32
	Height float32
}

type CAMetalDrawable struct {
	ptr unsafe.Pointer
}

func (s *CAMetalDrawable) Texture() *MTLTexture {
	ptr := C.CAMetalDrawable_texture(s.ptr)
	return &MTLTexture{ptr}
}

type MTLTexture struct {
	ptr unsafe.Pointer
}

type MTKView struct {
	ptr unsafe.Pointer
}

func (s *MTKView) CurrentRenderPassDescriptor() *MTLRenderPassDescriptor {
	ptr := C.MTKView_currentRenderPassDescriptor(s.ptr)
	return &MTLRenderPassDescriptor{ptr}
}

func (s *MTKView) Device() *MTLDevice {
	ptr := C.MTKView_device(s.ptr)
	return &MTLDevice{ptr}
}

func (s *MTKView) Layer() *CAMetalLayer {
	ptr := C.MTKView_layer(s.ptr)
	return &CAMetalLayer{ptr}
}

type MTLDevice struct {
	ptr unsafe.Pointer
}

func (s *MTLDevice) NewCommandQueue() *MTLCommandQueue {
	ptr := C.MTLDevice_newCommandQueue(s.ptr)
	return &MTLCommandQueue{ptr}
}

func (s *MTLDevice) NewLibraryWithSource(source string) *MTLLibrary {
	ptr := C.MTLDevice_newLibraryWithSource(s.ptr, C.CString(source), nil)
	return &MTLLibrary{ptr}
}

func (s *MTLDevice) NewBufferWithBytes(data unsafe.Pointer, size uintptr, count int, options MTLResourceOptions) *MTLBuffer {
	fmt.Println(C.int(size) * C.int(count))
	ptr := C.MTLDevice_newBufferWithBytes(s.ptr, data, C.int(size)*C.int(count), C.MTLResourceOptions(options))
	return &MTLBuffer{ptr, count}
}

func (s *MTLDevice) NewBufferWithByteArray(data []byte, size int, count int, options MTLResourceOptions) *MTLBuffer {
	ptr := C.MTLDevice_newBufferWithBytes(s.ptr, unsafe.Pointer(&data[0]), C.int(size), C.MTLResourceOptions(options))
	return &MTLBuffer{ptr, count}
}

func (s *MTLDevice) NewBufferWithVectors(vertices []Vector_float4, options MTLResourceOptions) *MTLBuffer {
	size := int(unsafe.Sizeof(vertices[0])) * len(vertices)
	ptr := C.MTLDevice_newBufferWithVectors(s.ptr, (*C.vector_float4)(&(vertices[0])), C.int(size), C.MTLResourceOptions(options))
	return &MTLBuffer{ptr, len(vertices)}
}

func (s *MTLDevice) NewBufferWithVectors2(vertices [][2]Vector_float4, options MTLResourceOptions) *MTLBuffer {
	size := int(unsafe.Sizeof(vertices[0])) * len(vertices)
	ptr := C.MTLDevice_newBufferWithVectors(s.ptr, (*C.vector_float4)(&(vertices[0][0])), C.int(size), C.MTLResourceOptions(options))
	return &MTLBuffer{ptr, len(vertices)}
}

func (s *MTLDevice) NewBufferWithInts(vertices []Uint16, options MTLResourceOptions) *MTLBuffer {
	size := int(unsafe.Sizeof(vertices[0])) * len(vertices)
	ptr := C.MTLDevice_newBufferWithInts(s.ptr, (*C.uint16_t)(&(vertices[0])), C.int(size), C.MTLResourceOptions(options))
	return &MTLBuffer{ptr, len(vertices)}
}

func (s *MTLDevice) NewRenderPipelineStateWithDescriptor(pipelineDescriptor *MTLRenderPipelineDescriptor) *MTLRenderPipelineState {
	ptr := C.MTLDevice_newRenderPipelineStateWithDescriptor(s.ptr, pipelineDescriptor.ptr)
	return &MTLRenderPipelineState{ptr}
}

func (s *MTLDevice) NewDepthStencilStateWithDescriptor(stencilDescriptor *MTLDepthStencilDescriptor) *MTLDepthStencilState {
	ptr := C.MTLDevice_newDepthStencilStateWithDescriptor(s.ptr, stencilDescriptor.ptr)
	return &MTLDepthStencilState{ptr}
}

type MTLDepthStencilState struct {
	ptr unsafe.Pointer
}

type MTLBuffer struct {
	ptr   unsafe.Pointer
	Count int
}

func (s *MTLBuffer) Contents() unsafe.Pointer {
	return C.MTLBuffer_contents(s.ptr)
}

func (s *MTLBuffer) ContentsCopy(data unsafe.Pointer, size uintptr, offset uintptr) unsafe.Pointer {
	return C.memcpy(unsafe.Pointer(uintptr(s.Contents())+offset), data, C.ulong(size))
}

type MTLRenderPipelineState struct {
	ptr unsafe.Pointer
}

type MTLLibrary struct {
	ptr unsafe.Pointer
}

func (s *MTLLibrary) NewFunctionWithName(name string) *MTLFunction {
	ptr := C.MTLLibrary_newFunctionWithName(s.ptr, C.CString(name))
	return &MTLFunction{ptr}
}

type MTLFunction struct {
	ptr unsafe.Pointer
}

type MTLCommandQueue struct {
	ptr unsafe.Pointer
}

func (s *MTLCommandQueue) CommandBuffer() *MTLCommandBuffer {
	ptr := C.MTLCommandQueue_commandBuffer(s.ptr)
	return &MTLCommandBuffer{ptr}
}

type MTLCommandBuffer struct {
	ptr unsafe.Pointer
}

func (s *MTLCommandBuffer) RenderCommandEncoderWithDescriptor(passDescriptor *MTLRenderPassDescriptor) *MTLRenderCommandEncoder {
	ptr := C.MTLCommandBuffer_renderCommandEncoderWithDescriptor(s.ptr, passDescriptor.ptr)
	return &MTLRenderCommandEncoder{ptr}
}

func (s *MTLCommandBuffer) PresentDrawable(drawable *CAMetalDrawable) {
	C.MTLCommandBuffer_presentDrawable(s.ptr, drawable.ptr)
}

func (s *MTLCommandBuffer) Commit() {
	C.MTLCommandBuffer_commit(s.ptr)
}

type MTLRenderCommandEncoder struct {
	ptr unsafe.Pointer
}

func (s *MTLRenderCommandEncoder) SetRenderPipelineState(ps *MTLRenderPipelineState) {
	C.MTLRenderCommandEncoder_setRenderPipelineState(s.ptr, ps.ptr)
}

func (s *MTLRenderCommandEncoder) SetDepthStencilState(dss *MTLDepthStencilState) {
	C.MTLRenderCommandEncoder_setDepthStencilState(s.ptr, dss.ptr)
}

func (s *MTLRenderCommandEncoder) SetFrontFacingWinding(winding MTLWinding) {
	C.MTLRenderCommandEncoder_setFrontFacingWinding(s.ptr, C.MTLWinding(winding))
}

func (s *MTLRenderCommandEncoder) SetCullMode(cullmode MTLCullMode) {
	C.MTLRenderCommandEncoder_setCullMode(s.ptr, C.MTLCullMode(cullmode))
}

func (s *MTLRenderCommandEncoder) SetVertexBuffer(ps *MTLBuffer, offset int, atIndex int) {
	C.MTLRenderCommandEncoder_setVertexBuffer(s.ptr, ps.ptr, C.int(offset), C.int(atIndex))
}

func (s *MTLRenderCommandEncoder) DrawPrimitives(primitiveType MTLPrimitiveType, vertexStart int, vertexCount int) {
	C.MTLRenderCommandEncoder_drawPrimitives(s.ptr, C.MTLPrimitiveType(primitiveType), C.int(vertexStart), C.int(vertexCount))
}

func (s *MTLRenderCommandEncoder) DrawIndexedPrimitives(primitiveType MTLPrimitiveType, indexCount int, indexType MTLIndexType, indexBuffer *MTLBuffer, indexBufferOffset int) {
	C.MTLRenderCommandEncoder_drawIndexedPrimitives(s.ptr, C.MTLPrimitiveType(primitiveType), C.int(indexCount), C.MTLIndexType(indexType), indexBuffer.ptr, C.int(indexBufferOffset))
}

func (s *MTLRenderCommandEncoder) EndEncoding() {
	C.MTLRenderCommandEncoder_endEncoding(s.ptr)
}

type MTLRenderPassDescriptor struct {
	ptr unsafe.Pointer
}

func (s *MTLRenderPassDescriptor) ColorAttachment(ca *ColorAttachment) {
	ptr := C.MTLRenderPassDescriptor_colorAttachments(s.ptr, 0)
	C.colorAttachments_set_loadAction(ptr, C.MTLLoadAction(ca.LoadAction))
	C.colorAttachments_set_storeAction(ptr, C.MTLStoreAction(ca.StoreAction))
	C.colorAttachments_set_clearColor(ptr, C.MTLClearColor{C.double(ca.ClearColor.Red), C.double(ca.ClearColor.Green), C.double(ca.ClearColor.Blue), C.double(ca.ClearColor.Alpha)})
	C.colorAttachments_set_texture(ptr, ca.Texture.ptr)
}

type MTLRenderPipelineDescriptor struct {
	ptr unsafe.Pointer
}

func (s *MTLRenderPipelineDescriptor) SetVertexFunction(fn *MTLFunction) {
	C.MTLRenderPipelineDescriptor_set_vertexFunction(s.ptr, fn.ptr)
}

func (s *MTLRenderPipelineDescriptor) SetFragmentFunction(fn *MTLFunction) {
	C.MTLRenderPipelineDescriptor_set_fragmentFunction(s.ptr, fn.ptr)
}

func (s *MTLRenderPipelineDescriptor) SetDepthAttachmentPixelFormat(pixelFormat MTLPixelFormat) {
	C.MTLRenderPipelineDescriptor_set_depthAttachmentPixelFormat(s.ptr, C.MTLPixelFormat(pixelFormat))
}

type MTLRenderPipelineColorAttachmentDescriptor struct {
	ptr unsafe.Pointer
}

func (s *MTLRenderPipelineColorAttachmentDescriptor) SetPixelFormat(pixelFormat MTLPixelFormat) {
	C.colorAttachments_set_pixelFormat(s.ptr, C.MTLPixelFormat(pixelFormat))
}

func (s *MTLRenderPipelineDescriptor) ColorAttachment(idx int) *MTLRenderPipelineColorAttachmentDescriptor {
	ptr := C.MTLRenderPipelineDescriptor_colorAttachments(s.ptr, C.int(idx))
	return &MTLRenderPipelineColorAttachmentDescriptor{ptr}
}

func NewMTLRenderPipelineDescriptor() *MTLRenderPipelineDescriptor {
	ptr := C.MTLRenderPipelineDescriptor_new()
	return &MTLRenderPipelineDescriptor{ptr}
}

type MTLDepthStencilDescriptor struct {
	ptr unsafe.Pointer
}

func (s *MTLDepthStencilDescriptor) SetDepthCompareFunction(depthCompareFunction MTLCompareFunction) {
	C.MTLDepthStencilDescriptor_set_depthCompareFunction(s.ptr, C.MTLCompareFunction(depthCompareFunction))
}

func (s *MTLDepthStencilDescriptor) SetDepthWriteEnabled(enabled bool) {
	C.MTLDepthStencilDescriptor_set_depthWriteEnabled(s.ptr, C._Bool(enabled))
}

func NewMTLDepthStencilDescriptor() *MTLDepthStencilDescriptor {
	ptr := C.MTLDepthStencilDescriptor_new()
	return &MTLDepthStencilDescriptor{ptr}
}

func Vector4(a, b, c, d float64) Vector_float4 {
	return Vector_float4{C.float(a), C.float(b), C.float(c), C.float(d)}
}

func Vector3(a, b, c float64) Vector_float3 {
	return Vector_float3{C.float(a), C.float(b), C.float(c)}
}

func Vector2(a, b float64) Vector_float2 {
	return Vector_float2{C.float(a), C.float(b)}
}

type (
	MTLLoadAction      C.MTLLoadAction
	MTLStoreAction     C.MTLStoreAction
	Vector_float4      C.vector_float4
	Vector_float3      C.vector_float3
	Vector_float2      C.vector_float2
	Matrix_float4x4    C.matrix_float4x4
	Float              C.float
	Uint16             C.uint16_t
	MTLPixelFormat     C.MTLPixelFormat
	MTLResourceOptions C.MTLResourceOptions
	MTLPrimitiveType   C.MTLPrimitiveType
	MTLIndexType       C.MTLIndexType
	MTLCompareFunction C.MTLCompareFunction
	MTLWinding         C.MTLWinding
	MTLCullMode        C.MTLCullMode
)

type MTLClearColor struct {
	Red   float64
	Green float64
	Blue  float64
	Alpha float64
}

const (
	MTLPixelFormatBGRA8Unorm   MTLPixelFormat = C.MTLPixelFormatBGRA8Unorm
	MTLPixelFormatDepth32Float MTLPixelFormat = C.MTLPixelFormatDepth32Float

	MTLLoadActionClear       MTLLoadAction      = C.MTLLoadActionClear
	MTLStoreActionStore      MTLStoreAction     = C.MTLStoreActionStore
	MTLPrimitiveTypeTriangle MTLPrimitiveType   = C.MTLPrimitiveTypeTriangle
	MTLIndexTypeUInt16       MTLIndexType       = C.MTLIndexTypeUInt16
	MTLCompareFunctionLess   MTLCompareFunction = C.MTLCompareFunctionLess

	MTLWindingCounterClockwise MTLWinding  = C.MTLWindingCounterClockwise
	MTLCullModeBack            MTLCullMode = C.MTLCullModeBack

	MTLResourceCPUCacheModeDefaultCache = C.MTLResourceCPUCacheModeDefaultCache
)

type ColorAttachment struct {
	Texture     *MTLTexture
	LoadAction  MTLLoadAction
	StoreAction MTLStoreAction
	ClearColor  MTLClearColor
}

type MBEVertex struct {
	Position Vector_float4
	Color    Vector_float4
}

//export renderInitWithMetalKitView
func renderInitWithMetalKitView(mkViewPtr unsafe.Pointer) {
	viewAddr := uintptr(mkViewPtr)
	view := &MTKView{mkViewPtr}
	delegates[viewAddr].init(view)
}

//export renderDrawInMTKView
func renderDrawInMTKView(mkViewPtr unsafe.Pointer) {
	viewAddr := uintptr(mkViewPtr)
	view := &MTKView{mkViewPtr}
	delegates[viewAddr].draw(view)
}
