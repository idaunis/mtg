#include <simd/SIMD.h>
#include <MetalKit/MetalKit.h>

#include "renderer.h"
#include "renderer2.h"

void renderDrawInMTKView(MTKView *);
void renderInitWithMetalKitView(MTKView *);

void *CAMetalLayer_nextDrawable(void *metalLayer) {
    return (id<CAMetalDrawable>) [(CAMetalLayer *)metalLayer nextDrawable];
}

void *MTKView_currentRenderPassDescriptor(void *view) {
    return (MTLRenderPassDescriptor *) ((MTKView *) view).currentRenderPassDescriptor;
}

void *MTKView_device(void *view) {
    return (id<MTLDevice>) ((MTKView *) view).device;
}

void *MTKView_layer(void *view) {
    return (CAMetalLayer*) ((MTKView *) view).layer;
}

void *CAMetalDrawable_texture(void *drawable) {
    return (id<MTLTexture>) ((id<CAMetalDrawable>) drawable).texture;
}

void *MTLDevice_newCommandQueue(void *device) {
    return (id<MTLCommandQueue>) [(id<MTLDevice>)device newCommandQueue];
}

void *MTLCommandQueue_commandBuffer(void *commandQueue) {
    return (id<MTLCommandBuffer>) [(id<MTLCommandQueue>)commandQueue commandBuffer];
}

void *MTLCommandBuffer_renderCommandEncoderWithDescriptor(void *commandBuffer, void *passDescriptor) {
    return (id<MTLRenderCommandEncoder>) [(id<MTLCommandBuffer>)commandBuffer renderCommandEncoderWithDescriptor:(MTLRenderPassDescriptor *)passDescriptor];
}

void MTLCommandBuffer_presentDrawable(void *commandBuffer, void *drawable) {
    [(id<MTLCommandBuffer>)commandBuffer presentDrawable:(id<CAMetalDrawable>)drawable];
}

void MTLCommandBuffer_commit(void *commandBuffer) {
    [(id<MTLCommandBuffer>)commandBuffer commit];
}

void MTLRenderCommandEncoder_endEncoding(void *commandEncoder) {
    [(id<MTLRenderCommandEncoder>) commandEncoder endEncoding];
}

void MTLRenderPassDescriptor_colorAttachments(void *passDescriptor, MTLLoadAction loadAction, MTLStoreAction storeAction, MTLClearColor clearColor, void *texture) {
    ((MTLRenderPassDescriptor *) passDescriptor).colorAttachments[0].loadAction = loadAction;
    ((MTLRenderPassDescriptor *) passDescriptor).colorAttachments[0].storeAction = storeAction;
    ((MTLRenderPassDescriptor *) passDescriptor).colorAttachments[0].clearColor = clearColor;
    ((MTLRenderPassDescriptor *) passDescriptor).colorAttachments[0].texture = (id<MTLTexture>) texture;
}

@implementation Renderer
{

}

- (nonnull instancetype)initWithMetalKitView:(nonnull MTKView *)view {
    self = [super init];
    if (self) {
        renderInitWithMetalKitView(view);
    }
    return self;
}

- (void)drawInMTKView:(nonnull MTKView *)view
{
    renderDrawInMTKView(view);
}

- (void)mtkView:(nonnull MTKView *)view drawableSizeWillChange:(CGSize)size
{
    NSLog(@"drawableSizeWillChange");
}

@end
