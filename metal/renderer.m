#include <simd/SIMD.h>
#include <MetalKit/MetalKit.h>

#include "renderer.h"

void renderDrawInMTKView(MTKView *);
void renderInitWithMetalKitView(MTKView *);

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
}

@end
