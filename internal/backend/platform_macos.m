//go:build !ios && darwin

#include "_cgo_export.h"

#import <Cocoa/Cocoa.h>

@interface GniAppDelegate : NSObject<NSApplicationDelegate>
@end

@interface GniView : NSView
{
}
@end

@interface GniWindowDelegate : NSObject<NSWindowDelegate>
@end

@implementation GniView
- (void)drawRect:(CGRect)rect {   
    gni_paint(CGRectGetWidth(rect), CGRectGetHeight(rect));
}
@end

@implementation GniAppDelegate
- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
	[[NSRunningApplication currentApplication] activateWithOptions:(NSApplicationActivateAllWindows | NSApplicationActivateIgnoringOtherApps)];
}
@end

@implementation GniWindowDelegate
@end


void gni_exec() {
    @autoreleasepool {      
        [NSApplication sharedApplication];
        [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
        GniAppDelegate *ad = [[GniAppDelegate alloc] init];
		[NSApp setDelegate:ad];

        id mainMenu = [NSMenuItem new];
		id menu = [NSMenu new];
		id hideMenuItem = [[NSMenuItem alloc] initWithTitle:@"Hide"
            action:@selector(hide:)
			keyEquivalent:@"h"];
		[menu addItem:hideMenuItem];
		NSMenuItem *quitMenuItem = [[NSMenuItem alloc] initWithTitle:@"Quit"
			action:@selector(terminate:)
			keyEquivalent:@"q"];
		[menu addItem:quitMenuItem];
		[mainMenu setSubmenu:menu];
		
        id menuBar = [NSMenu new];
		[menuBar addItem:mainMenu];
		[NSApp setMainMenu:menuBar];
        
        NSRect rect = NSMakeRect(0, 0, 800, 600);
        NSUInteger styleMask = NSWindowStyleMaskTitled |
            NSWindowStyleMaskClosable |
            NSWindowStyleMaskMiniaturizable |
			NSWindowStyleMaskResizable;
        NSWindow* window = [[NSWindow alloc] initWithContentRect:rect
                    styleMask:styleMask
                    backing:NSBackingStoreBuffered
                    defer:NO];
	    window.title = @"GNI";
        [window setBackgroundColor:[NSColor whiteColor]];
        [window makeKeyAndOrderFront:NSApp];

        GniView* view = [[GniView alloc] initWithFrame:CGRectZero];
        [window setContentView:view];
        GniWindowDelegate *wd = [[GniWindowDelegate alloc] init];
        [window setDelegate:wd];

        [NSApp run];
    }
}