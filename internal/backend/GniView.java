package dev.gni;

import android.view.View;
import android.graphics.Canvas;
import android.content.Context;

public final class GniView extends View {

    public GniView(Context context) {
		super(context);
	}

    @Override
    protected void onDraw(Canvas canvas) {
        super.onDraw(canvas);
        gniDraw(canvas);
    }

    private static native void gniDraw(Canvas canvas);
}
