package metal

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Metal -framework MetalKit
#import <Cocoa/Cocoa.h>
#import <MetalKit/MetalKit.h>
#include "renderer.h"
int StartApp() {
    [NSAutoreleasePool new];
    [NSApplication sharedApplication];
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    return 0;
}
int Menu() {
    id menubar = [[NSMenu new] autorelease];
    id appMenuItem = [[NSMenuItem new] autorelease];
    [menubar addItem:appMenuItem];
    [NSApp setMainMenu:menubar];
    id appMenu = [[NSMenu new] autorelease];
    id appName = [[NSProcessInfo processInfo] processName];
    id quitTitle = [@"Quit " stringByAppendingString:appName];
    id quitMenuItem = [[[NSMenuItem alloc] initWithTitle:quitTitle
        action:@selector(terminate:) keyEquivalent:@"q"]
          	autorelease];
    [appMenu addItem:quitMenuItem];
    [appMenuItem setSubmenu:appMenu];

    return 0;
}
void *Window() {
    id appName = [[NSProcessInfo processInfo] processName];
    NSWindow * window = [[[NSWindow alloc] initWithContentRect: NSMakeRect(0, 0, 640, 480)
        styleMask: NSWindowStyleMaskTitled | NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskClosable | NSWindowStyleMaskResizable
            backing: NSBackingStoreBuffered defer: NO]
            autorelease];

    [window center];
    [window setTitle:appName];
    [window makeKeyAndOrderFront:nil];

    NSView *view = (NSView *) window.contentView;

    MTKView *mtkView = [[MTKView alloc]
                         initWithFrame: [view frame]
                                device: MTLCreateSystemDefaultDevice()];

    // mtkView.enableSetNeedsDisplay = NO;
    // mtkView.preferredFramesPerSecond = 60;
    mtkView.translatesAutoresizingMaskIntoConstraints = false;
    [view addSubview:mtkView];
    [view addConstraints: [NSLayoutConstraint constraintsWithVisualFormat:@"|[mtkView]|"
                                                                      options: 0
                                                                      metrics: nil
                                                                        views: @{@"mtkView":mtkView}]];

    [view addConstraints: [NSLayoutConstraint constraintsWithVisualFormat:@"V:|[mtkView]|"
                                                                      options: 0
                                                                      metrics: nil
                                                                        views: @{@"mtkView":mtkView}]];
    if (!mtkView.device) {
        NSLog(@"Metal is not supported on this device");
    }

    mtkView.clearColor = MTLClearColorMake(0.0, 0.5, 0.25, 1.0);

    return mtkView;
}
void RenderDelegate(void *mtkView) {
    Renderer * renderer = [[Renderer alloc] initWithMetalKitView:(MTKView *)mtkView];

    if (!renderer) {
        NSLog(@"Renderer initialization failed");
    }

    [renderer mtkView:(MTKView *)mtkView drawableSizeWillChange:((MTKView *)mtkView).bounds.size];

    ((MTKView *)mtkView).delegate = renderer;

    NSLog(@"ok");
}
void RunApp() {
    [NSApp activateIgnoringOtherApps:YES];
    [NSApp run];
}
*/
import "C"

func CreateApp() {
	C.StartApp()
	C.Menu()
}

func CreateWindow() *MTKView {
	ptr := C.Window()
	return &MTKView{ptr}
}

type delegateFuncs struct {
	init, draw func(*MTKView)
}

var delegates map[uintptr]delegateFuncs

func RenderDelegate(view *MTKView, init, draw func(*MTKView)) {
	viewAddr := uintptr(view.ptr)
	if delegates == nil {
		delegates = make(map[uintptr]delegateFuncs)
	}
	delegates[viewAddr] = delegateFuncs{init, draw}
	C.RenderDelegate(view.ptr)
}

func RunApp() {
	C.RunApp()
}
