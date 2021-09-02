#include <MetalKit/MetalKit.h>
#include <simd/SIMD.h>

vector_float2 new_vector_float2(float x, float y);
vector_float3 new_vector_float3(float x, float y, float z);
vector_float4 new_vector_float4(float x, float y, float z, float w);
matrix_float3x3 new_matrix_float3x3(float *x, float *y, float *z);
matrix_float4x4 new_matrix_float4x4(vector_float4 x, vector_float4 y, vector_float4 z, vector_float4 w);
matrix_float4x4 simd_matrix_multiply(matrix_float4x4 a, matrix_float4x4 b);
vector_float3 simd_vector3_cross(vector_float3 a, vector_float3 b);
vector_float4 simd_vector4_normalize(vector_float4 a);
void *CAMetalLayer_nextDrawable(void *metalLayer);
CGSize CAMetalLayer_drawableSize(void *metalLayer);
void *MTKView_currentRenderPassDescriptor(void *view);
void *MTKView_device(void *view);
void *MTKView_layer(void *view);
void MTKView_set_depthStencilPixelFormat(void *view, MTLPixelFormat pixelFormat);
void *CAMetalDrawable_texture(void *drawable);
void *MTLDevice_newCommandQueue(void *device);
void *MTLDevice_newLibraryWithData(void *device, void *data);
void *MTLDevice_newLibraryWithSource(void *device, char *source, void *options);
void *MTLDevice_newRenderPipelineStateWithDescriptor(void *device, void *pipelineDescriptor);
void *MTLDevice_newBufferWithBytes(void *device, void *data, int length, MTLResourceOptions options);
void *MTLDevice_newBufferWithVectors(void *device, vector_float4 vertices[], int length, MTLResourceOptions options);
void *MTLDevice_newBufferWithInts(void *device, uint16_t vertices[], int length, MTLResourceOptions options);
void *MTLDevice_newDepthStencilStateWithDescriptor(void *device, void *stencilDescriptor);
void *MTLLibrary_newFunctionWithName(void *library, char *name);
void *MTLCommandQueue_commandBuffer(void *commandQueue);
void *MTLCommandBuffer_renderCommandEncoderWithDescriptor(void *commandBuffer, void *passDescriptor);
void MTLCommandBuffer_presentDrawable(void *commandBuffer, void *drawable);
void MTLCommandBuffer_commit(void *commandBuffer);
void MTLRenderCommandEncoder_endEncoding(void *commandEncoder);
void MTLRenderCommandEncoder_setRenderPipelineState(void *commandEncoder, void *pipelineState);
void MTLRenderCommandEncoder_setDepthStencilState(void *commandEncoder, void *depthStencilState);
void MTLRenderCommandEncoder_setVertexBuffer(void *commandEncoder, void *vb, int offset, int atIndex);
void MTLRenderCommandEncoder_drawPrimitives(void *commandEncoder, MTLPrimitiveType type, int start, int count);
void MTLRenderCommandEncoder_drawIndexedPrimitives(void *commandEncoder, MTLPrimitiveType primitiveType, int indexCount, MTLIndexType indexType, void *indexBuffer, int indexBufferOffset);
void MTLRenderCommandEncoder_setFrontFacingWinding(void *commandEncoder, MTLWinding winding);
void MTLRenderCommandEncoder_setCullMode(void *commandEncoder, MTLCullMode cullmode);

void *MTLRenderPipelineDescriptor_new();
void MTLRenderPipelineDescriptor_set_vertexFunction(void *pdesc, void *fn);
void MTLRenderPipelineDescriptor_set_fragmentFunction(void *pdesc, void *fn);
void MTLRenderPipelineDescriptor_set_depthAttachmentPixelFormat(void *pdesc, MTLPixelFormat pixelFormat);

void *MTLRenderPassDescriptor_colorAttachments(void *passDescriptor, int idx);
void colorAttachments_set_loadAction(void *cad, MTLLoadAction loadAction);
void colorAttachments_set_storeAction(void *cad, MTLStoreAction storeAction);
void colorAttachments_set_clearColor(void *cad, MTLClearColor clearColor);
void colorAttachments_set_texture(void *cad, void *texture);

void *MTLRenderPipelineDescriptor_colorAttachments(void *pdesc, int ids);
void colorAttachments_set_pixelFormat(void *cad, MTLPixelFormat pixelFormat);

void *MTLDepthStencilDescriptor_new();
void MTLDepthStencilDescriptor_set_depthCompareFunction(void *dsdesc, MTLCompareFunction dcfun);
void MTLDepthStencilDescriptor_set_depthWriteEnabled(void *dsdesc, bool enabled);

void *MTLBuffer_contents(void *buffer);
