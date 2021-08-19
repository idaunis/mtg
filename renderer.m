#include <simd/SIMD.h>
#include <MetalKit/MetalKit.h>

#include "renderer.h"
#include "renderer2.h"

void renderDrawInMTKView(CAMetalLayer *, MTKView *);
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

void *CAMetalDrawable_texture(void *drawable) {
    return (id<MTLTexture>) ((id<CAMetalDrawable>) drawable).texture;
}

void *MTLDevice_newCommandQueue(void *device) {
    return (id<MTLCommandQueue>) [(id<MTLDevice>)device newCommandQueue];
}


void MTLRenderPassDescriptor_colorAttachments(void *passDescriptor, MTLLoadAction loadAction, MTLStoreAction storeAction, MTLClearColor clearColor, void *texture) {
    ((MTLRenderPassDescriptor *) passDescriptor).colorAttachments[0].loadAction = loadAction;
    ((MTLRenderPassDescriptor *) passDescriptor).colorAttachments[0].storeAction = storeAction;
    ((MTLRenderPassDescriptor *) passDescriptor).colorAttachments[0].clearColor = clearColor;
    ((MTLRenderPassDescriptor *) passDescriptor).colorAttachments[0].texture = (id<MTLTexture>) texture;
}

@implementation Renderer
{
    id<MTLDevice> _device;

    id<MTLCommandQueue> _commandQueue;

    CAMetalLayer * metalLayer;
}

- (nonnull instancetype)initWithMetalKitView:(nonnull MTKView *)view {
    self = [super init];
    if (self) {
        renderInitWithMetalKitView(view);
        metalLayer = (CAMetalLayer*) view.layer;
        _device = view.device;
        _commandQueue = [_device newCommandQueue];
    }
    return self;
}

- (void)drawInMTKView:(nonnull MTKView *)view
{
    renderDrawInMTKView(metalLayer, view);
}


/// Called whenever the view needs to render a frame.
- (void)drawInMTKView2:(nonnull MTKView *)view
{
    NSLog(@"draw");
    // The render pass descriptor references the texture into which Metal should draw
    //MTLRenderPassDescriptor *passDescriptor = [MTLRenderPassDescriptor renderPassDescriptor];
    MTLRenderPassDescriptor *passDescriptor = view.currentRenderPassDescriptor;

    if (passDescriptor == nil) {
        return;
    }

    // Get the drawable that will be presented at the end of the frame

    /*
    id<MTLDrawable> drawable = view.currentDrawable;
    id<MTLTexture> texture = [drawable texture];
    */

    NSLog(@"%.20f",test());

    id<CAMetalDrawable> drawable = [metalLayer nextDrawable];
    id<MTLTexture> texture = drawable.texture;
    
    
    passDescriptor.colorAttachments[0].loadAction = MTLLoadActionClear;
    passDescriptor.colorAttachments[0].storeAction = MTLStoreActionStore;
    passDescriptor.colorAttachments[0].clearColor = MTLClearColorMake(1, test(), 1, 1);
    passDescriptor.colorAttachments[0].texture = texture;
    passDescriptor.colorAttachments[0].clearColor = MTLClearColorMake(0, 0, 1, 1);

    id<MTLCommandBuffer> commandBuffer = [_commandQueue commandBuffer];
    
    // Create a render pass and immediately end encoding, causing the drawable to be cleared
    id<MTLRenderCommandEncoder> commandEncoder = [commandBuffer renderCommandEncoderWithDescriptor:passDescriptor];
    
    [commandEncoder endEncoding];

    // Request that the drawable texture be presented by the windowing system once drawing is done
    [commandBuffer presentDrawable:drawable];
    
    [commandBuffer commit];
}


/// Called whenever view changes orientation or is resized
- (void)mtkView:(nonnull MTKView *)view drawableSizeWillChange:(CGSize)size
{
    NSLog(@"drawableSizeWillChange");
}

@end
