package main

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
int Window() {
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

    mtkView.enableSetNeedsDisplay = YES;
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

    Renderer * renderer = [[Renderer alloc] initWithMetalKitView:mtkView];

    if (!renderer) {
        NSLog(@"Renderer initialization failed");
    }

    [renderer mtkView:mtkView drawableSizeWillChange:mtkView.bounds.size];

    mtkView.delegate = renderer;

    NSLog(@"ok");

    [NSApp activateIgnoringOtherApps:YES];
    [NSApp run];
    return 0;
}
*/
import "C"

func main() {
	C.StartApp()
	C.Menu()
	C.Window()
}
