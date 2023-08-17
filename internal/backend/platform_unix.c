//go:build (linux && !android) || freebsd || openbsd

#include "_cgo_export.h"
#include <cairo-xcb.h>
#include <stdio.h>
#include <stdlib.h>
#include <xcb/xcb.h>

#define MIN(X, Y) (((X) < (Y)) ? (X) : (Y))

static xcb_connection_t *conn;

const char *xcb_errors[] = {
    "Success",
    "BadRequest",
    "BadValue",
    "BadWindow",
    "BadPixmap",
    "BadAtom",
    "BadCursor",
    "BadFont",
    "BadMatch",
    "BadDrawable",
    "BadAccess",
    "BadAlloc",
    "BadColor",
    "BadGC",
    "BadIDChoice",
    "BadName",
    "BadLength",
    "BadImplementation",
    "Unknown"
};

const char *xcb_protocol_request_codes[] = {
    "Null",
    "CreateWindow",
    "ChangeWindowAttributes",
    "GetWindowAttributes",
    "DestroyWindow",
    "DestroySubwindows",
    "ChangeSaveSet",
    "ReparentWindow",
    "MapWindow",
    "MapSubwindows",
    "UnmapWindow",
    "UnmapSubwindows",
    "ConfigureWindow",
    "CirculateWindow",
    "GetGeometry",
    "QueryTree",
    "InternAtom",
    "GetAtomName",
    "ChangeProperty",
    "DeleteProperty",
    "GetProperty",
    "ListProperties",
    "SetSelectionOwner",
    "GetSelectionOwner",
    "ConvertSelection",
    "SendEvent",
    "GrabPointer",
    "UngrabPointer",
    "GrabButton",
    "UngrabButton",
    "ChangeActivePointerGrab",
    "GrabKeyboard",
    "UngrabKeyboard",
    "GrabKey",
    "UngrabKey",
    "AllowEvents",
    "GrabServer",
    "UngrabServer",
    "QueryPointer",
    "GetMotionEvents",
    "TranslateCoords",
    "WarpPointer",
    "SetInputFocus",
    "GetInputFocus",
    "QueryKeymap",
    "OpenFont",
    "CloseFont",
    "QueryFont",
    "QueryTextExtents",
    "ListFonts",
    "ListFontsWithInfo",
    "SetFontPath",
    "GetFontPath",
    "CreatePixmap",
    "FreePixmap",
    "CreateGC",
    "ChangeGC",
    "CopyGC",
    "SetDashes",
    "SetClipRectangles",
    "FreeGC",
    "ClearArea",
    "CopyArea",
    "CopyPlane",
    "PolyPoint",
    "PolyLine",
    "PolySegment",
    "PolyRectangle",
    "PolyArc",
    "FillPoly",
    "PolyFillRectangle",
    "PolyFillArc",
    "PutImage",
    "GetImage",
    "PolyText8",
    "PolyText16",
    "ImageText8",
    "ImageText16",
    "CreateColormap",
    "FreeColormap",
    "CopyColormapAndFree",
    "InstallColormap",
    "UninstallColormap",
    "ListInstalledColormaps",
    "AllocColor",
    "AllocNamedColor",
    "AllocColorCells",
    "AllocColorPlanes",
    "FreeColors",
    "StoreColors",
    "StoreNamedColor",
    "QueryColors",
    "LookupColor",
    "CreateCursor",
    "CreateGlyphCursor",
    "FreeCursor",
    "RecolorCursor",
    "QueryBestSize",
    "QueryExtension",
    "ListExtensions",
    "ChangeKeyboardMapping",
    "GetKeyboardMapping",
    "ChangeKeyboardControl",
    "GetKeyboardControl",
    "Bell",
    "ChangePointerControl",
    "GetPointerControl",
    "SetScreenSaver",
    "GetScreenSaver",
    "ChangeHosts",
    "ListHosts",
    "SetAccessControl",
    "SetCloseDownMode",
    "KillClient",
    "RotateProperties",
    "ForceScreenSaver",
    "SetPointerMapping",
    "GetPointerMapping",
    "SetModifierMapping",
    "GetModifierMapping",
    "Unknown"
};

static void handle_error(xcb_generic_error_t *e)
{
    uint8_t clamped_error_code = MIN(e->error_code, (sizeof(xcb_errors) / sizeof(xcb_errors[0])) - 1);
    uint8_t clamped_major_code = MIN(e->major_code, (sizeof(xcb_protocol_request_codes) / sizeof(xcb_protocol_request_codes[0])) - 1);

    fprintf(stderr, "XCB error: %d (%s), sequence: %d, resource id: %d, major code: %d (%s), minor code: %d\n",
        e->error_code, xcb_errors[clamped_error_code],
        e->sequence,
        e->resource_id,
        e->major_code, xcb_protocol_request_codes[clamped_major_code],
        e->minor_code);
}

static xcb_visualtype_t *find_visual(xcb_connection_t *c, xcb_visualid_t visual)
{
    xcb_screen_iterator_t screen_iter = xcb_setup_roots_iterator(xcb_get_setup(c));
    for (; screen_iter.rem; xcb_screen_next(&screen_iter)) {
        xcb_depth_iterator_t depth_iter = xcb_screen_allowed_depths_iterator(screen_iter.data);
        for (; depth_iter.rem; xcb_depth_next(&depth_iter)) {
            xcb_visualtype_iterator_t visual_iter = xcb_depth_visuals_iterator(depth_iter.data);
            for (; visual_iter.rem; xcb_visualtype_next(&visual_iter)) {
                if (visual == visual_iter.data->visual_id) {
                    return visual_iter.data;
                }
            }
        }
    }
    return NULL;
}

static xcb_window_t create_window(xcb_screen_t *s)
{
    uint32_t mask;
    uint32_t values[2];

    xcb_window_t window;

    mask = XCB_CW_BACK_PIXEL | XCB_CW_EVENT_MASK;
    values[0] = s->white_pixel;
    values[1] = XCB_EVENT_MASK_EXPOSURE | XCB_EVENT_MASK_KEY_PRESS;

    window = xcb_generate_id(conn);
    xcb_create_window(conn,
        XCB_COPY_FROM_PARENT, window, s->root,
        0, 0, 640, 480,
        0,
        XCB_WINDOW_CLASS_INPUT_OUTPUT,
        s->root_visual,
        mask, values);

    xcb_map_window(conn, window);
    return window;
}

void gni_exec()
{
    int screen_number;
    cairo_surface_t *surface = NULL;
    cairo_t *cr = NULL;

    conn = xcb_connect(NULL, &screen_number);
    if (xcb_connection_has_error(conn)) {
        fprintf(stderr, "can't connect to an X server\n");
        goto out;
    }

    xcb_screen_iterator_t iter = xcb_setup_roots_iterator(xcb_get_setup(conn));
    for (int i = 0; i < screen_number && iter.rem; i++) {
        xcb_screen_next(&iter);
    }
    if (iter.rem == 0) {
        fprintf(stderr, "no screens available\n");
        goto out;
    }
    xcb_screen_t *s = iter.data;

    xcb_visualtype_t *visual = find_visual(conn, s->root_visual);
    if (visual == NULL) {
        fprintf(stderr, "couldn't find visual\n");
        goto out;
    }
    xcb_window_t w = create_window(s);
    surface = cairo_xcb_surface_create(conn, w, visual, 640, 480);
    cr = cairo_create(surface);
    xcb_flush(conn);

    xcb_generic_event_t *e;
    while ((e = xcb_wait_for_event(conn))) {
        if (!(e->response_type & ~0x80)) {
            handle_error((xcb_generic_error_t *)(e));
            free(e);
            continue;
        }

        switch (e->response_type & ~0x80) {
        case XCB_EXPOSE:
            gni_paint(cr);
            cairo_surface_flush(surface);
            xcb_flush(conn);
            break;
        case XCB_KEY_PRESS: {
            xcb_key_press_event_t *ev;
            ev = (xcb_key_press_event_t *)e;
            if (ev->detail == 9) {
                goto out;
            }
            break;
        }
        }
        free(e);
    }
out:
    cairo_destroy(cr);
    cairo_surface_destroy(surface);
    xcb_disconnect(conn);
}
