#include <MetalKit/MetalKit.h>

void *CAMetalLayer_nextDrawable(void *metalLayer);
void *MTKView_currentRenderPassDescriptor(void *view);
void *MTKView_device(void *view);
void *MTKView_layer(void *view);
void *CAMetalDrawable_texture(void *drawable);
void *MTLDevice_newCommandQueue(void *device);
void *MTLDevice_newLibraryWithData(void *device, void *data);
void *MTLDevice_newLibraryWithSource(void *device, char *source, void *options);
void *MTLLibrary_newFunctionWithName(void *library, char *name);
void *MTLCommandQueue_commandBuffer(void *commandQueue);
void *MTLCommandBuffer_renderCommandEncoderWithDescriptor(void *commandBuffer, void *passDescriptor);
void MTLCommandBuffer_presentDrawable(void *commandBuffer, void *drawable);
void MTLCommandBuffer_commit(void *commandBuffer);
void MTLRenderCommandEncoder_endEncoding(void *commandEncoder);
void MTLRenderPassDescriptor_colorAttachments(void *passDescriptor, MTLLoadAction loadAction, MTLStoreAction storeAction, MTLClearColor clearColor, void *texture);
