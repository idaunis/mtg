package metal

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Metal -framework MetalKit
#include "renderer2.h"
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

func (s *MTLDevice) NewBufferWithBytes(vertices []MBEVertex, size int, options int) {
	fmt.Println("newBufferWithBytes", vertices, size, options)
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

func (s *MTLRenderCommandEncoder) EndEncoding() {
	C.MTLRenderCommandEncoder_endEncoding(s.ptr)
}

type MTLRenderPassDescriptor struct {
	ptr unsafe.Pointer
}

func (s *MTLRenderPassDescriptor) ColorAttachment(ca *ColorAttachment) {
	C.MTLRenderPassDescriptor_colorAttachments(
		s.ptr,
		C.MTLLoadAction(ca.LoadAction),
		C.MTLStoreAction(ca.StoreAction),
		C.MTLClearColor{C.double(ca.ClearColor.Red), C.double(ca.ClearColor.Green), C.double(ca.ClearColor.Blue), C.double(ca.ClearColor.Alpha)},
		ca.Texture.ptr,
	)
}

type (
	MTLLoadAction  C.MTLLoadAction
	MTLStoreAction C.MTLStoreAction
	Vector_float4  C.vector_float4
)

type MTLClearColor struct {
	Red   float64
	Green float64
	Blue  float64
	Alpha float64
}

const (
	MTLLoadActionClear  MTLLoadAction  = C.MTLLoadActionClear
	MTLStoreActionStore MTLStoreAction = C.MTLStoreActionStore

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
