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

//export test
func test() float64 {
	return float64(0.2)
}

var (
	device       *MTLDevice
	commandQueue *MTLCommandQueue
)

//export renderInitWithMetalKitView
func renderInitWithMetalKitView(mkViewPtr unsafe.Pointer) {
	view := &MTKView{mkViewPtr}
	device = view.device()
	commandQueue = device.newCommandQueue()

	fmt.Println("InitWithMetalKitView", view, device, commandQueue)
}

//export renderDrawInMTKView
func renderDrawInMTKView(metalLayerPtr unsafe.Pointer, mkViewPtr unsafe.Pointer) {
	metalLayer := &CAMetalLayer{metalLayerPtr}
	view := &MTKView{mkViewPtr}

	drawable := metalLayer.nextDrawable()
	texture := drawable.texture()

	renderPassDescriptor := view.currentRenderPassDescriptor()

	renderPassDescriptor.colorAttachment(&ColorAttachment{
		texture:     texture,
		loadAction:  MTLLoadActionClear,
		storeAction: MTLStoreActionStore,
		clearColor:  MTLClearColor{0, 1, 0, 1},
	})

	fmt.Println(metalLayer, drawable, renderPassDescriptor, texture)
}
