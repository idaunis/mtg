#include <MetalKit/MetalKit.h>

double test();

void *CAMetalLayer_nextDrawable(void *metalLayer);
void *MTKView_currentRenderPassDescriptor(void *view);
void *MTKView_device(void *view);
void *CAMetalDrawable_texture(void *drawable);
void *MTLDevice_newCommandQueue(void *device);
void MTLRenderPassDescriptor_colorAttachments(void *passDescriptor, MTLLoadAction loadAction, MTLStoreAction storeAction, MTLClearColor clearColor, void *texture);
