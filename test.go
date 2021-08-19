package main

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

func (s *CAMetalLayer) nextDrawable() *CAMetalDrawable {
	ptr := C.CAMetalLayer_nextDrawable(s.ptr)
	return &CAMetalDrawable{ptr}
}

type CAMetalDrawable struct {
	ptr unsafe.Pointer
}

func (s *CAMetalDrawable) texture() *MTLTexture {
	ptr := C.CAMetalDrawable_texture(s.ptr)
	return &MTLTexture{ptr}
}

type MTLTexture struct {
	ptr unsafe.Pointer
}

type MTKView struct {
	ptr unsafe.Pointer
}

func (s *MTKView) currentRenderPassDescriptor() *MTLRenderPassDescriptor {
	ptr := C.MTKView_currentRenderPassDescriptor(s.ptr)
	return &MTLRenderPassDescriptor{ptr}
}

func (s *MTKView) device() *MTLDevice {
	ptr := C.MTKView_device(s.ptr)
	return &MTLDevice{ptr}
}

func (s *MTKView) layer() *CAMetalLayer {
	ptr := C.MTKView_layer(s.ptr)
	return &CAMetalLayer{ptr}
}

type MTLDevice struct {
	ptr unsafe.Pointer
}

func (s *MTLDevice) newCommandQueue() *MTLCommandQueue {
	ptr := C.MTLDevice_newCommandQueue(s.ptr)
	return &MTLCommandQueue{ptr}
}

type MTLCommandQueue struct {
	ptr unsafe.Pointer
}

func (s *MTLCommandQueue) commandBuffer() *MTLCommandBuffer {
	ptr := C.MTLCommandQueue_commandBuffer(s.ptr)
	return &MTLCommandBuffer{ptr}
}

type MTLCommandBuffer struct {
	ptr unsafe.Pointer
}

func (s *MTLCommandBuffer) renderCommandEncoderWithDescriptor(passDescriptor *MTLRenderPassDescriptor) *MTLRenderCommandEncoder {
	ptr := C.MTLCommandBuffer_renderCommandEncoderWithDescriptor(s.ptr, passDescriptor.ptr)
	return &MTLRenderCommandEncoder{ptr}
}

func (s *MTLCommandBuffer) presentDrawable(drawable *CAMetalDrawable) {
	C.MTLCommandBuffer_presentDrawable(s.ptr, drawable.ptr)
}

func (s *MTLCommandBuffer) commit() {
	C.MTLCommandBuffer_commit(s.ptr)
}

type MTLRenderCommandEncoder struct {
	ptr unsafe.Pointer
}

func (s *MTLRenderCommandEncoder) endEncoding() {
	C.MTLRenderCommandEncoder_endEncoding(s.ptr)
}

type MTLRenderPassDescriptor struct {
	ptr unsafe.Pointer
}

func (s *MTLRenderPassDescriptor) colorAttachment(ca *ColorAttachment) {
	C.MTLRenderPassDescriptor_colorAttachments(
		s.ptr,
		C.MTLLoadAction(ca.loadAction),
		C.MTLStoreAction(ca.storeAction),
		C.MTLClearColor(ca.clearColor),
		ca.texture.ptr,
	)
}

type (
	MTLLoadAction  C.MTLLoadAction
	MTLStoreAction C.MTLStoreAction
	MTLClearColor  C.MTLClearColor
)

const (
	MTLLoadActionClear  MTLLoadAction  = C.MTLLoadActionClear
	MTLStoreActionStore MTLStoreAction = C.MTLStoreActionStore
)

type ColorAttachment struct {
	texture     *MTLTexture
	loadAction  MTLLoadAction
	storeAction MTLStoreAction
	clearColor  MTLClearColor
}

var (
	device       *MTLDevice
	commandQueue *MTLCommandQueue
	metalLayer   *CAMetalLayer
)

//export renderInitWithMetalKitView
func renderInitWithMetalKitView(mkViewPtr unsafe.Pointer) {
	view := &MTKView{mkViewPtr}

	device = view.device()
	metalLayer = view.layer()
	commandQueue = device.newCommandQueue()

	fmt.Println("InitWithMetalKitView", view, device, metalLayer, commandQueue)
}

//export renderDrawInMTKView
func renderDrawInMTKView(mkViewPtr unsafe.Pointer) {
	view := &MTKView{mkViewPtr}

	drawable := metalLayer.nextDrawable()
	texture := drawable.texture()

	renderPassDescriptor := view.currentRenderPassDescriptor()

	renderPassDescriptor.colorAttachment(&ColorAttachment{
		texture:     texture,
		loadAction:  MTLLoadActionClear,
		storeAction: MTLStoreActionStore,
		clearColor:  MTLClearColor{1, 1, 0, 1},
	})

	commandBuffer := commandQueue.commandBuffer()
	commandEncoder := commandBuffer.renderCommandEncoderWithDescriptor(renderPassDescriptor)

	commandEncoder.endEncoding()
	commandBuffer.presentDrawable(drawable)
	commandBuffer.commit()

	fmt.Println(metalLayer, drawable, renderPassDescriptor, texture, commandEncoder)
}
