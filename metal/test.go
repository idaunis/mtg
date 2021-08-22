package metal

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Metal -framework MetalKit
#include "renderer2.h"
*/
import "C"

import (
	"unsafe"
)

type CAMetalLayer struct {
	ptr unsafe.Pointer
}

func (s *CAMetalLayer) NextDrawable() *CAMetalDrawable {
	ptr := C.CAMetalLayer_nextDrawable(s.ptr)
	return &CAMetalDrawable{ptr}
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
	ptr := C.MTLDevice_newBufferWithBytes(s.ptr, data, C.int(size), C.MTLResourceOptions(options))
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

type MTLBuffer struct {
	ptr   unsafe.Pointer
	Count int
}

func (s *MTLDevice) NewRenderPipelineStateWithDescriptor(pipelineDescriptor *MTLRenderPipelineDescriptor) *MTLRenderPipelineState {
	ptr := C.MTLDevice_newRenderPipelineStateWithDescriptor(s.ptr, pipelineDescriptor.ptr)
	return &MTLRenderPipelineState{ptr}
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

type (
	MTLLoadAction      C.MTLLoadAction
	MTLStoreAction     C.MTLStoreAction
	Vector_float4      C.vector_float4
	Matrix_float4x4    [4]C.vector_float4
	Uint16             C.uint16_t
	MTLPixelFormat     C.MTLPixelFormat
	MTLResourceOptions C.MTLResourceOptions
	MTLPrimitiveType   C.MTLPrimitiveType
	MTLIndexType       C.MTLIndexType
)

type MTLClearColor struct {
	Red   float64
	Green float64
	Blue  float64
	Alpha float64
}

const (
	MTLPixelFormatBGRA8Unorm MTLPixelFormat   = C.MTLPixelFormatBGRA8Unorm
	MTLLoadActionClear       MTLLoadAction    = C.MTLLoadActionClear
	MTLStoreActionStore      MTLStoreAction   = C.MTLStoreActionStore
	MTLPrimitiveTypeTriangle MTLPrimitiveType = C.MTLPrimitiveTypeTriangle
	MTLIndexTypeUInt16       MTLIndexType     = C.MTLIndexTypeUInt16

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
