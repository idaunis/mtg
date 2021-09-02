#include <simd/SIMD.h>
#include <MetalKit/MetalKit.h>

#include "renderer.h"
#include "renderer2.h"

void renderDrawInMTKView(MTKView *);
void renderInitWithMetalKitView(MTKView *);

vector_float2 new_vector_float2(float x, float y) {
    return vector2(x, y);
}

vector_float3 new_vector_float3(float x, float y, float z) {
    return vector3(x, y, z);
}

vector_float4 new_vector_float4(float x, float y, float z, float w) {
    return vector4(x, y, z, w);
}

matrix_float3x3 new_matrix_float3x3(float *x, float *y, float *z) {
    // NSLog(@"x: %f %f %f", x[0], x[1], x[2]);
    // NSLog(@"y: %f %f %f", y[0], y[1], y[2]);
    // NSLog(@"z: %f %f %f", z[0], z[1], z[2]);
    matrix_float3x3 mat = {vector3(x[0], x[1], x[2]), vector3(y[0], y[1], y[2]), vector3(z[0], z[1], z[2])};
    return mat;
}

matrix_float4x4 new_matrix_float4x4(vector_float4 x, vector_float4 y, vector_float4 z, vector_float4 w) {
    matrix_float4x4 mat = {x, y, z, w};
    return mat;
}

matrix_float4x4 simd_matrix_multiply(matrix_float4x4 a, matrix_float4x4 b) {
    return matrix_multiply(a, b);
}

vector_float3 simd_vector3_cross(vector_float3 x, vector_float3 y) {
    NSLog(@"cros x: %f %f %f", x[0], x[1], x[2]);
    NSLog(@"cros y: %f %f %f", y[0], y[1], y[2]);
    return vector_cross(x, y);
}

vector_float4 simd_vector4_normalize(vector_float4 x) {
    NSLog(@"normalize x: %f %f %f %f", x[0], x[1], x[2], x[3]);
    return vector_normalize(x);
}

void *CAMetalLayer_nextDrawable(void *metalLayer) {
    return (id<CAMetalDrawable>) [(CAMetalLayer *)metalLayer nextDrawable];
}

CGSize CAMetalLayer_drawableSize(void *metalLayer) {
    return ((CAMetalLayer *)metalLayer).drawableSize;
}

void *MTKView_currentRenderPassDescriptor(void *view) {
    return (MTLRenderPassDescriptor *) ((MTKView *) view).currentRenderPassDescriptor;
}

void MTKView_set_depthStencilPixelFormat(void *view, MTLPixelFormat pixelFormat) {
    ((MTKView *)view).depthStencilPixelFormat = pixelFormat;
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

void *MTLBuffer_contents(void *buffer) {
    return ((id<MTLBuffer>) buffer).contents;
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
