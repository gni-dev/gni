package dev.gni;

import android.app.Activity;
import android.content.Context;
import android.os.Bundle;
import android.util.Log;

public final class GniActivity extends Activity {
    private GniView view;

	@Override
	protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

		view = new GniView(this);
		setContentView(view);
		gniCreate(this);
    }

    static {
        System.loadLibrary("gni");
    }

	private static native void gniCreate(Context context);
}
