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

void *MTLDevice_newLibraryWithData(void *device, void *data) {
    return (id<MTLLibrary>) [(id<MTLDevice>)device newLibraryWithData:(dispatch_data_t)data error:nil];
}

void *MTLDevice_newLibraryWithSource(void *device, char *source, void *options) {
    NSString *nsSource = [NSString stringWithUTF8String:source];
    return (id<MTLLibrary>) [(id<MTLDevice>)device newLibraryWithSource:nsSource options:(MTLCompileOptions *)options error:nil];
}

void *MTLDevice_newRenderPipelineStateWithDescriptor(void *device, void *pipelineDescriptor) {
    return (id<MTLRenderPipelineState>) [(id<MTLDevice>)device newRenderPipelineStateWithDescriptor:(MTLRenderPipelineDescriptor *)pipelineDescriptor error:nil];
}

void *MTLDevice_newBufferWithBytes(void *device, void *data, int length, MTLResourceOptions options) {
    return (id<MTLBuffer>) [(id<MTLDevice>)device newBufferWithBytes:data length:length options:options];
}

void *MTLDevice_newBufferWithVectors(void *device, vector_float4 vertices[], int length, MTLResourceOptions options) {
    return (id<MTLBuffer>) [(id<MTLDevice>)device newBufferWithBytes:vertices length:length options:options];
}

void *MTLDevice_newBufferWithInts(void *device, uint16_t vertices[], int length, MTLResourceOptions options) {
    return (id<MTLBuffer>) [(id<MTLDevice>)device newBufferWithBytes:vertices length:length options:options];
}

void *MTLDevice_newDepthStencilStateWithDescriptor(void *device, void *stencilDescriptor) {
    return (id<MTLDepthStencilState>) [(id<MTLDevice>)device newDepthStencilStateWithDescriptor:(MTLDepthStencilDescriptor *)stencilDescriptor];
}

void *MTLLibrary_newFunctionWithName(void *library, char *name) {
    NSString *nsName = [NSString stringWithUTF8String:name];
    return (id<MTLFunction>) [(id<MTLLibrary>)library newFunctionWithName:nsName];
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

void MTLRenderCommandEncoder_setRenderPipelineState(void *commandEncoder, void *pipelineState) {
    [(id<MTLRenderCommandEncoder>) commandEncoder setRenderPipelineState:(id<MTLRenderPipelineState>)pipelineState];
}

void MTLRenderCommandEncoder_setDepthStencilState(void *commandEncoder, void *depthStencilState) {
    [(id<MTLRenderCommandEncoder>) commandEncoder setDepthStencilState:(id<MTLDepthStencilState>)depthStencilState];
}

void MTLRenderCommandEncoder_setFrontFacingWinding(void *commandEncoder, MTLWinding winding) {
    [(id<MTLRenderCommandEncoder>) commandEncoder setFrontFacingWinding:winding];
}

void MTLRenderCommandEncoder_setCullMode(void *commandEncoder, MTLCullMode cullmode) {
    [(id<MTLRenderCommandEncoder>) commandEncoder setCullMode:cullmode];
}

void MTLRenderCommandEncoder_setVertexBuffer(void *commandEncoder, void *vertexBuffer, int offset, int atIndex) {
    [(id<MTLRenderCommandEncoder>) commandEncoder setVertexBuffer:(id<MTLBuffer>)vertexBuffer offset:offset atIndex:atIndex];
}

void MTLRenderCommandEncoder_drawPrimitives(void *commandEncoder, MTLPrimitiveType primitiveType, int start, int count) {
    [(id<MTLRenderCommandEncoder>) commandEncoder drawPrimitives:primitiveType vertexStart:start vertexCount:count];
}

void MTLRenderCommandEncoder_drawIndexedPrimitives(void *commandEncoder, MTLPrimitiveType primitiveType, int indexCount, MTLIndexType indexType, void *indexBuffer, int indexBufferOffset) {    
    [(id<MTLRenderCommandEncoder>) commandEncoder drawIndexedPrimitives:primitiveType indexCount:indexCount indexType:indexType indexBuffer:(id<MTLBuffer>)indexBuffer indexBufferOffset:indexBufferOffset];
}

void *MTLRenderPipelineDescriptor_colorAttachments(void *pdesc, int idx) {
    //returns MTLRenderPipelineColorAttachmentDescriptor *
    return ((MTLRenderPipelineDescriptor *) pdesc).colorAttachments[idx];
}

void colorAttachments_set_pixelFormat(void *cad, MTLPixelFormat pixelFormat) {
    ((MTLRenderPipelineColorAttachmentDescriptor *)cad).pixelFormat = pixelFormat; 
}

void *MTLRenderPassDescriptor_colorAttachments(void *rpdesc, int idx) {
    //returns MTLRenderPassColorAttachmentDescriptor *
    return ((MTLRenderPassDescriptor *) rpdesc).colorAttachments[idx];
}

void colorAttachments_set_loadAction(void *cad, MTLLoadAction loadAction) {
    ((MTLRenderPassColorAttachmentDescriptor *)cad).loadAction = loadAction; 
}

void colorAttachments_set_storeAction(void *cad, MTLStoreAction storeAction) {
    ((MTLRenderPassColorAttachmentDescriptor *)cad).storeAction = storeAction;
}

void colorAttachments_set_clearColor(void *cad, MTLClearColor clearColor) {
    ((MTLRenderPassColorAttachmentDescriptor *)cad).clearColor = clearColor;    
}

void colorAttachments_set_texture(void *cad, void *texture) {
    ((MTLRenderPassColorAttachmentDescriptor *)cad).texture = (id<MTLTexture>) texture;
}

void *MTLRenderPipelineDescriptor_new() {
    return (MTLRenderPipelineDescriptor *) [MTLRenderPipelineDescriptor new];
}

void MTLRenderPipelineDescriptor_set_vertexFunction(void *pdesc, void *fn) {
    ((MTLRenderPipelineDescriptor *) pdesc).vertexFunction = (id<MTLFunction>) fn;
}

void MTLRenderPipelineDescriptor_set_fragmentFunction(void *pdesc, void *fn) {
    ((MTLRenderPipelineDescriptor *) pdesc).fragmentFunction = (id<MTLFunction>) fn;
}

void MTLRenderPipelineDescriptor_set_depthAttachmentPixelFormat(void *pdesc, MTLPixelFormat pixelFormat) {
    ((MTLRenderPipelineDescriptor *) pdesc).depthAttachmentPixelFormat = pixelFormat;
}

void *MTLDepthStencilDescriptor_new() {
    return (MTLDepthStencilDescriptor *) [MTLDepthStencilDescriptor new];
}

void MTLDepthStencilDescriptor_set_depthCompareFunction(void *dsdesc, MTLCompareFunction dcfun) {
    ((MTLDepthStencilDescriptor *) dsdesc).depthCompareFunction = dcfun;
}

void MTLDepthStencilDescriptor_set_depthWriteEnabled(void *dsdesc, bool enabled) {
    ((MTLDepthStencilDescriptor *) dsdesc).depthWriteEnabled = enabled;
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
